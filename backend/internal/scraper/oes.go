package scraper

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kendall1978/project_salary/backend/internal/database"
	"github.com/kendall1978/project_salary/backend/internal/models"
	"github.com/xuri/excelize/v2"
)

// OESRegion holds parsed region data from one Excel row.
type OESRegion struct {
	AreaCode  string
	AreaTitle string
	AreaType  int
	StateCode string
	State     string
}

// OESSnapshot holds parsed salary data from one Excel row.
type OESSnapshot struct {
	AreaCode        string
	OccCode         string
	OccTitle        string
	TotalEmployment int
	MedianSalary    float64
	MeanSalary      float64
	Percentile10    float64
	Percentile25    float64
	Percentile75    float64
	Percentile90    float64
	Year            int
}

// ProcessExcelFile reads an OES Excel file and returns parsed regions and snapshots.
// It skips group totals (OCC_CODE ending in -0000) and rows with suppressed median salary.
func ProcessExcelFile(path string, year int) ([]OESRegion, []OESSnapshot, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.Rows(sheets[0])
	if err != nil {
		return nil, nil, fmt.Errorf("get rows: %w", err)
	}
	defer rows.Close()

	// Read header row to build column index map
	if !rows.Next() {
		return nil, nil, fmt.Errorf("empty sheet")
	}
	headerCells, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("read header: %w", err)
	}

	colIdx := make(map[string]int)
	for i, name := range headerCells {
		colIdx[strings.TrimSpace(name)] = i
	}

	// Verify required columns exist
	required := []string{"AREA", "AREA_TITLE", "AREA_TYPE", "PRIM_STATE", "OCC_CODE", "OCC_TITLE", "A_MEDIAN"}
	for _, col := range required {
		if _, ok := colIdx[col]; !ok {
			return nil, nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	regionSet := make(map[string]OESRegion)
	var snapshots []OESSnapshot

	for rows.Next() {
		cells, err := rows.Columns()
		if err != nil {
			continue
		}

		getCol := func(name string) string {
			idx, ok := colIdx[name]
			if !ok || idx >= len(cells) {
				return ""
			}
			return strings.TrimSpace(cells[idx])
		}

		occCode := getCol("OCC_CODE")
		if ShouldSkipOccCode(occCode) {
			continue
		}

		// Skip rows with suppressed median salary
		medianStr := getCol("A_MEDIAN")
		median, ok := ParseNumber(medianStr)
		if !ok {
			continue
		}

		areaCode := getCol("AREA")
		areaTitle := getCol("AREA_TITLE")
		areaTypeStr := getCol("AREA_TYPE")
		stateCode := getCol("PRIM_STATE")

		areaType, _ := strconv.Atoi(areaTypeStr)

		// Upsert region
		if _, exists := regionSet[areaCode]; !exists {
			regionSet[areaCode] = OESRegion{
				AreaCode:  areaCode,
				AreaTitle: areaTitle,
				AreaType:  areaType,
				StateCode: stateCode,
				State:     StateName(stateCode),
			}
		}

		totEmp, _ := ParseInt(getCol("TOT_EMP"))
		mean, _ := ParseNumber(getCol("A_MEAN"))
		p10, _ := ParseNumber(getCol("A_PCT10"))
		p25, _ := ParseNumber(getCol("A_PCT25"))
		p75, _ := ParseNumber(getCol("A_PCT75"))
		p90, _ := ParseNumber(getCol("A_PCT90"))

		snapshots = append(snapshots, OESSnapshot{
			AreaCode:        areaCode,
			OccCode:         occCode,
			OccTitle:        getCol("OCC_TITLE"),
			TotalEmployment: totEmp,
			MedianSalary:    median,
			MeanSalary:      mean,
			Percentile10:    p10,
			Percentile25:    p25,
			Percentile75:    p75,
			Percentile90:    p90,
			Year:            year,
		})
	}

	regions := make([]OESRegion, 0, len(regionSet))
	for _, r := range regionSet {
		regions = append(regions, r)
	}

	return regions, snapshots, nil
}

// FetchOESData downloads and imports OES data for the given years.
// If localZipDir is non-empty, it looks for zip files there instead of downloading.
// Set backfill=true to import years 2018 through current; false for latest year only.
func FetchOESData(backfill bool, localZipDir string) error {
	currentYear := time.Now().Year()
	latestAvailable := currentYear - 1 // OES data lags by ~1 year

	var years []int
	if backfill {
		for y := 2018; y <= latestAvailable; y++ {
			years = append(years, y)
		}
	} else {
		years = []int{latestAvailable}
	}

	scrapeLog := models.ScrapeLog{
		Source:    "bls_oes",
		Status:    "running",
		StartedAt: time.Now(),
	}
	database.DB.Create(&scrapeLog)

	totalRecords := 0

	for _, year := range years {
		// Check if data for this year already exists
		var count int64
		database.DB.Model(&models.SalarySnapshot{}).Where("year = ?", year).Count(&count)
		if count > 0 {
			log.Printf("OES data for %d already exists (%d records), skipping", year, count)
			continue
		}

		yy := fmt.Sprintf("%02d", year%100)

		xlsxPath, cleanup, err := getExcelFile(year, yy, localZipDir)
		if err != nil {
			log.Printf("Failed to get Excel file for %d: %v", year, err)
			continue
		}

		regions, snapshots, err := ProcessExcelFile(xlsxPath, year)
		cleanup()
		if err != nil {
			log.Printf("Failed to process Excel file for %d: %v", year, err)
			continue
		}

		imported, err := importToDatabase(regions, snapshots, year)
		if err != nil {
			log.Printf("Failed to import data for %d: %v", year, err)
			continue
		}

		totalRecords += imported
		log.Printf("Imported %d records for year %d", imported, year)
	}

	now := time.Now()
	scrapeLog.Status = "success"
	scrapeLog.RecordsFound = totalRecords
	scrapeLog.CompletedAt = &now
	database.DB.Save(&scrapeLog)

	return nil
}

// getExcelFile returns the path to the xlsx file and a cleanup function.
func getExcelFile(year int, yy string, localZipDir string) (string, func(), error) {
	var zipPath string
	var tempDir string

	if localZipDir != "" {
		// Look for local zip file
		zipPath = filepath.Join(localZipDir, fmt.Sprintf("oesm%sst.zip", yy))
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			return "", func() {}, fmt.Errorf("local zip not found: %s", zipPath)
		}
	} else {
		// Download from BLS
		url := fmt.Sprintf("https://www.bls.gov/oes/special-requests/oesm%sst.zip", yy)
		log.Printf("Downloading %s", url)

		var err error
		tempDir, err = os.MkdirTemp("", "oes-download-*")
		if err != nil {
			return "", func() {}, fmt.Errorf("create temp dir: %w", err)
		}

		zipPath = filepath.Join(tempDir, fmt.Sprintf("oesm%sst.zip", yy))
		if err := downloadFile(url, zipPath); err != nil {
			os.RemoveAll(tempDir)
			return "", func() {}, fmt.Errorf("download: %w", err)
		}
	}

	// Extract xlsx from zip
	extractDir, err := os.MkdirTemp("", "oes-extract-*")
	if err != nil {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
		return "", func() {}, fmt.Errorf("create extract dir: %w", err)
	}

	xlsxPath, err := extractXlsx(zipPath, extractDir)
	if err != nil {
		os.RemoveAll(extractDir)
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
		return "", func() {}, fmt.Errorf("extract: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(extractDir)
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	}

	return xlsxPath, cleanup, nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractXlsx(zipPath, destDir string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Ext(f.Name) != ".xlsx" {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		destPath := filepath.Join(destDir, filepath.Base(f.Name))
		out, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return "", err
		}

		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return "", err
		}

		return destPath, nil
	}

	return "", fmt.Errorf("no xlsx file found in zip")
}

func importToDatabase(regions []OESRegion, snapshots []OESSnapshot, year int) (int, error) {
	// Build region lookup: AreaCode -> DB Region
	regionMap := make(map[string]models.Region)

	for _, r := range regions {
		var dbRegion models.Region
		result := database.DB.Where("area_code = ?", r.AreaCode).FirstOrCreate(&dbRegion, models.Region{
			AreaCode:  r.AreaCode,
			AreaTitle: r.AreaTitle,
			AreaType:  r.AreaType,
			State:     r.State,
			StateCode: r.StateCode,
		})
		if result.Error != nil {
			return 0, fmt.Errorf("upsert region %s: %w", r.AreaCode, result.Error)
		}
		regionMap[r.AreaCode] = dbRegion
	}

	// Batch insert snapshots
	imported := 0
	batchSize := 500
	batch := make([]models.SalarySnapshot, 0, batchSize)

	for _, s := range snapshots {
		region, ok := regionMap[s.AreaCode]
		if !ok {
			continue
		}

		batch = append(batch, models.SalarySnapshot{
			RegionID:        region.ID,
			Source:          "bls",
			OccCode:         s.OccCode,
			OccTitle:        s.OccTitle,
			TotalEmployment: s.TotalEmployment,
			MedianSalary:    s.MedianSalary,
			MeanSalary:      s.MeanSalary,
			Percentile10:    s.Percentile10,
			Percentile25:    s.Percentile25,
			Percentile75:    s.Percentile75,
			Percentile90:    s.Percentile90,
			Year:            s.Year,
			SnapshotDate:    time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
		})

		if len(batch) >= batchSize {
			result := database.DB.Create(&batch)
			if result.Error != nil {
				log.Printf("Batch insert error (continuing): %v", result.Error)
			} else {
				imported += int(result.RowsAffected)
			}
			batch = batch[:0]
		}
	}

	// Insert remaining
	if len(batch) > 0 {
		result := database.DB.Create(&batch)
		if result.Error != nil {
			log.Printf("Final batch insert error: %v", result.Error)
		} else {
			imported += int(result.RowsAffected)
		}
	}

	return imported, nil
}