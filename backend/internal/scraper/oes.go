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
		for y := 2022; y <= latestAvailable; y++ {
			years = append(years, y)
		}
	} else {
		years = []int{latestAvailable}
	}

	// Set up file logger
	logDir := filepath.Join("logs")
	sl, err := NewScrapeLogger(logDir)
	if err != nil {
		log.Printf("Warning: could not create scrape log file: %v (logging to stdout only)", err)
	}
	if sl != nil {
		defer sl.Close()
		sl.Log("=== OES Scrape Started ===")
		sl.Log("Backfill: %v | Years: %v | LocalZipDir: %q", backfill, years, localZipDir)
	}

	logMsg := func(format string, args ...interface{}) {
		if sl != nil {
			sl.Log(format, args...)
		} else {
			log.Printf(format, args...)
		}
	}

	scrapeLog := models.ScrapeLog{
		Source:    "bls_oes",
		Status:    "running",
		StartedAt: time.Now(),
	}
	database.DB.Create(&scrapeLog)

	totalRecords := 0
	var scrapeErrors []string

	// Each BLS year has a state file (st) and a metro file (ma)
	fileSuffixes := []struct {
		suffix    string
		label     string
		areaType  int
	}{
		{"st", "state", 2},
		{"ma", "metro", 4},
	}

	for _, year := range years {
		yy := fmt.Sprintf("%02d", year%100)
		logMsg("--- Year %d ---", year)

		for _, fs := range fileSuffixes {
			// Check if data for this year+file already exists
			var count int64
			database.DB.Model(&models.SalarySnapshot{}).
				Joins("JOIN regions ON regions.id = salary_snapshots.region_id").
				Where("salary_snapshots.year = ? AND regions.area_type = ?", year, fs.areaType).Count(&count)
			if count > 0 {
				logMsg("SKIP %d %s: already has %d records", year, fs.label, count)
				continue
			}

			url := fmt.Sprintf("https://www.bls.gov/oes/special-requests/oesm%s%s.zip", yy, fs.suffix)
			logMsg("DOWNLOAD %d %s: %s", year, fs.label, url)

			xlsxPath, cleanup, err := getExcelFile(year, yy, fs.suffix, localZipDir)
			if err != nil {
				errMsg := fmt.Sprintf("FAIL %d %s download/extract: %v", year, fs.label, err)
				logMsg(errMsg)
				scrapeErrors = append(scrapeErrors, errMsg)
				continue
			}

			logMsg("PARSE %d %s: %s", year, fs.label, xlsxPath)
			regions, snapshots, err := ProcessExcelFile(xlsxPath, year)
			cleanup()
			if err != nil {
				errMsg := fmt.Sprintf("FAIL %d %s parse: %v", year, fs.label, err)
				logMsg(errMsg)
				scrapeErrors = append(scrapeErrors, errMsg)
				continue
			}

			// Log what was parsed
			areaTypeCounts := make(map[int]int)
			for _, r := range regions {
				areaTypeCounts[r.AreaType]++
			}
			logMsg("PARSED %d %s: %d regions (area_types: %v), %d snapshots", year, fs.label, len(regions), areaTypeCounts, len(snapshots))

			logMsg("IMPORT %d %s: writing to database...", year, fs.label)
			imported, err := importToDatabase(regions, snapshots, year)
			if err != nil {
				errMsg := fmt.Sprintf("FAIL %d %s import: %v", year, fs.label, err)
				logMsg(errMsg)
				scrapeErrors = append(scrapeErrors, errMsg)
				continue
			}

			totalRecords += imported
			logMsg("OK %d %s: imported %d records", year, fs.label, imported)
		}
	}

	now := time.Now()
	if len(scrapeErrors) > 0 {
		scrapeLog.Status = "partial"
		scrapeLog.ErrorMessage = fmt.Sprintf("%d errors; see log file for details", len(scrapeErrors))
		logMsg("=== Completed with %d errors ===", len(scrapeErrors))
		for i, e := range scrapeErrors {
			logMsg("  Error %d: %s", i+1, e)
		}
	} else {
		scrapeLog.Status = "success"
		logMsg("=== Completed successfully ===")
	}
	scrapeLog.RecordsFound = totalRecords
	scrapeLog.CompletedAt = &now
	database.DB.Save(&scrapeLog)

	logMsg("Total records imported: %d", totalRecords)
	if sl != nil {
		logMsg("Log file: %s", sl.Path)
	}

	return nil
}

// getExcelFile returns the path to the xlsx file and a cleanup function.
func getExcelFile(year int, yy string, suffix string, localZipDir string) (string, func(), error) {
	var zipPath string
	var tempDir string

	if localZipDir != "" {
		// Look for local zip file
		zipPath = filepath.Join(localZipDir, fmt.Sprintf("oesm%s%s.zip", yy, suffix))
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			return "", func() {}, fmt.Errorf("local zip not found: %s", zipPath)
		}
	} else {
		// Download from BLS
		url := fmt.Sprintf("https://www.bls.gov/oes/special-requests/oesm%s%s.zip", yy, suffix)
		log.Printf("Downloading %s", url)

		var err error
		tempDir, err = os.MkdirTemp("", "oes-download-*")
		if err != nil {
			return "", func() {}, fmt.Errorf("create temp dir: %w", err)
		}

		zipPath = filepath.Join(tempDir, fmt.Sprintf("oesm%s%s.zip", yy, suffix))
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

	xlsxPath, err := extractXlsx(zipPath, extractDir, suffix)
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
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://www.bls.gov/oes/tables.htm")

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

// extractXlsx extracts the target xlsx from a zip. The prefix selects which
// file to extract when a zip contains multiple xlsx files:
//   - "state" zips: look for state_*.xlsx or oesm*st*.xlsx
//   - "ma" (metro) zips: look for MSA_*.xlsx
//   - "bos" zips: look for BOS_*.xlsx
//
// If prefix is empty, the first xlsx found is returned.
func extractXlsx(zipPath, destDir, prefix string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var fallback *zip.File
	for _, f := range r.File {
		base := filepath.Base(f.Name)
		if filepath.Ext(base) != ".xlsx" || strings.HasPrefix(base, "~$") {
			continue
		}

		if fallback == nil {
			fallback = f
		}

		upperBase := strings.ToUpper(base)
		match := false
		switch prefix {
		case "st":
			match = strings.HasPrefix(upperBase, "STATE_") || strings.Contains(upperBase, "ST")
		case "ma":
			match = strings.HasPrefix(upperBase, "MSA_")
		case "bos":
			match = strings.HasPrefix(upperBase, "BOS_")
		default:
			match = true
		}

		if match {
			return extractZipFile(f, destDir)
		}
	}

	// Fall back to first xlsx if no prefix match
	if fallback != nil {
		return extractZipFile(fallback, destDir)
	}

	return "", fmt.Errorf("no xlsx file found in zip")
}

func extractZipFile(f *zip.File, destDir string) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	destPath := filepath.Join(destDir, filepath.Base(f.Name))
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func importToDatabase(regions []OESRegion, snapshots []OESSnapshot, year int) (int, error) {
	// Build region lookup: AreaCode -> DB Region
	regionMap := make(map[string]models.Region)

	for _, r := range regions {
		var dbRegion models.Region
		result := database.DB.Where("area_code = ?", r.AreaCode).First(&dbRegion)
		if result.Error != nil {
			// Not found — create it
			dbRegion = models.Region{
				AreaCode:  r.AreaCode,
				AreaTitle: r.AreaTitle,
				AreaType:  r.AreaType,
				State:     r.State,
				StateCode: r.StateCode,
			}
			if err := database.DB.Create(&dbRegion).Error; err != nil {
				return 0, fmt.Errorf("create region %s: %w", r.AreaCode, err)
			}
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