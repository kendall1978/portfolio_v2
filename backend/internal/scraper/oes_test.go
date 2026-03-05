package scraper

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func testdataPath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", name)
}

func TestProcessExcelFile(t *testing.T) {
	path := testdataPath("test_oes.xlsx")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Test fixture not found: %s", path)
	}

	regions, snapshots, err := ProcessExcelFile(path, 2024)
	if err != nil {
		t.Fatalf("ProcessExcelFile failed: %v", err)
	}

	// Should have 2 unique regions: Alabama (state) and Kansas City (metro)
	if len(regions) != 2 {
		t.Errorf("Expected 2 regions, got %d", len(regions))
	}

	// Should have 2 snapshots: Web Developers in AL + Web Developers in KC
	// Rows skipped: 00-0000 (group total), 11-0000 (major group), suppressed anesthesiologist
	if len(snapshots) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(snapshots))
	}

	// Verify Alabama state region
	var alRegion *OESRegion
	for i := range regions {
		if regions[i].AreaCode == "01" {
			alRegion = &regions[i]
			break
		}
	}
	if alRegion == nil {
		t.Fatal("Alabama region not found")
	}
	if alRegion.AreaTitle != "Alabama" {
		t.Errorf("Expected AreaTitle 'Alabama', got %q", alRegion.AreaTitle)
	}
	if alRegion.AreaType != 2 {
		t.Errorf("Expected AreaType 2, got %d", alRegion.AreaType)
	}
	if alRegion.StateCode != "AL" {
		t.Errorf("Expected StateCode 'AL', got %q", alRegion.StateCode)
	}
	if alRegion.State != "Alabama" {
		t.Errorf("Expected State 'Alabama', got %q", alRegion.State)
	}

	// Verify Kansas City metro region
	var kcRegion *OESRegion
	for i := range regions {
		if regions[i].AreaCode == "C2818" {
			kcRegion = &regions[i]
			break
		}
	}
	if kcRegion == nil {
		t.Fatal("Kansas City region not found")
	}
	if kcRegion.AreaType != 4 {
		t.Errorf("Expected AreaType 4, got %d", kcRegion.AreaType)
	}
	if kcRegion.State != "Missouri" {
		t.Errorf("Expected State 'Missouri', got %q", kcRegion.State)
	}

	// Verify Alabama Web Developer snapshot
	var alSnapshot *OESSnapshot
	for i := range snapshots {
		if snapshots[i].AreaCode == "01" && snapshots[i].OccCode == "15-1254" {
			alSnapshot = &snapshots[i]
			break
		}
	}
	if alSnapshot == nil {
		t.Fatal("Alabama Web Developer snapshot not found")
	}
	if alSnapshot.MedianSalary != 56160 {
		t.Errorf("Expected MedianSalary 56160, got %v", alSnapshot.MedianSalary)
	}
	if alSnapshot.Percentile25 != 46800 {
		t.Errorf("Expected Percentile25 46800, got %v", alSnapshot.Percentile25)
	}
	if alSnapshot.TotalEmployment != 1200 {
		t.Errorf("Expected TotalEmployment 1200, got %d", alSnapshot.TotalEmployment)
	}
	if alSnapshot.Year != 2024 {
		t.Errorf("Expected Year 2024, got %d", alSnapshot.Year)
	}
}