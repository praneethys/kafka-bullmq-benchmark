package metrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

func TestExportToJSON(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-results.json")

	results := []*common.BenchmarkResult{
		{
			QueueType:      "Test Queue",
			MessageCount:   1000,
			Duration:       10 * time.Second,
			Throughput:     100.0,
			SuccessCount:   1000,
			ErrorCount:     0,
			BytesProcessed: 1024000,
		},
	}

	err := ExportToJSON(results, filename)
	if err != nil {
		t.Fatalf("ExportToJSON failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("JSON file was not created")
	}

	// Verify file is not empty
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("JSON file is empty")
	}
}

func TestExportToCSV(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test-results.csv")

	results := []*common.BenchmarkResult{
		{
			QueueType:      "Test Queue 1",
			MessageCount:   1000,
			Duration:       10 * time.Second,
			Throughput:     100.0,
			MBPerSecond:    1.5,
			AvgLatency:     5 * time.Millisecond,
			P50Latency:     4 * time.Millisecond,
			P95Latency:     8 * time.Millisecond,
			P99Latency:     10 * time.Millisecond,
			MinLatency:     1 * time.Millisecond,
			MaxLatency:     15 * time.Millisecond,
			SuccessCount:   1000,
			ErrorCount:     0,
			BytesProcessed: 1024000,
		},
		{
			QueueType:      "Test Queue 2",
			MessageCount:   2000,
			Duration:       20 * time.Second,
			Throughput:     100.0,
			MBPerSecond:    2.0,
			AvgLatency:     6 * time.Millisecond,
			P50Latency:     5 * time.Millisecond,
			P95Latency:     9 * time.Millisecond,
			P99Latency:     12 * time.Millisecond,
			MinLatency:     2 * time.Millisecond,
			MaxLatency:     18 * time.Millisecond,
			SuccessCount:   2000,
			ErrorCount:     5,
			BytesProcessed: 2048000,
		},
	}

	err := ExportToCSV(results, filename)
	if err != nil {
		t.Fatalf("ExportToCSV failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("CSV file was not created")
	}

	// Verify file is not empty
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("CSV file is empty")
	}
}

func TestGenerateReport(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	results := []*common.BenchmarkResult{
		{
			QueueType:      "Test Queue",
			MessageCount:   1000,
			Duration:       10 * time.Second,
			Throughput:     100.0,
			SuccessCount:   1000,
			ErrorCount:     0,
			BytesProcessed: 1024000,
		},
	}

	err := GenerateReport(results, tempDir)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	// Check if both JSON and CSV files were created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("Expected at least 2 files (JSON and CSV), got %d", len(files))
	}

	hasJSON := false
	hasCSV := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			hasJSON = true
		}
		if filepath.Ext(file.Name()) == ".csv" {
			hasCSV = true
		}
	}

	if !hasJSON {
		t.Error("JSON file not found in report")
	}

	if !hasCSV {
		t.Error("CSV file not found in report")
	}
}

func TestPrintResults(t *testing.T) {
	result := &common.BenchmarkResult{
		QueueType:      "Test Queue",
		MessageCount:   1000,
		Duration:       10 * time.Second,
		Throughput:     100.0,
		MBPerSecond:    1.5,
		AvgLatency:     5 * time.Millisecond,
		P50Latency:     4 * time.Millisecond,
		P95Latency:     8 * time.Millisecond,
		P99Latency:     10 * time.Millisecond,
		MinLatency:     1 * time.Millisecond,
		MaxLatency:     15 * time.Millisecond,
		SuccessCount:   1000,
		ErrorCount:     0,
		BytesProcessed: 1024000,
	}

	// This should not panic
	PrintResults(result)
}

func TestCompareResults(t *testing.T) {
	results := []*common.BenchmarkResult{
		{
			QueueType:   "Queue 1",
			Throughput:  100.0,
			MBPerSecond: 1.5,
			AvgLatency:  5 * time.Millisecond,
			P99Latency:  10 * time.Millisecond,
		},
		{
			QueueType:   "Queue 2",
			Throughput:  150.0,
			MBPerSecond: 2.0,
			AvgLatency:  3 * time.Millisecond,
			P99Latency:  8 * time.Millisecond,
		},
	}

	// This should not panic
	CompareResults(results)
}
