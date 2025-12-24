package metrics

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

// ExportToJSON exports benchmark results to a JSON file
func ExportToJSON(results []*common.BenchmarkResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// ExportToCSV exports benchmark results to a CSV file
func ExportToCSV(results []*common.BenchmarkResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Queue Type",
		"Message Count",
		"Duration (s)",
		"Throughput (msg/s)",
		"MB/s",
		"Avg Latency (ms)",
		"P50 Latency (ms)",
		"P95 Latency (ms)",
		"P99 Latency (ms)",
		"Min Latency (ms)",
		"Max Latency (ms)",
		"Success Count",
		"Error Count",
		"Bytes Processed",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, result := range results {
		row := []string{
			result.QueueType,
			strconv.Itoa(result.MessageCount),
			fmt.Sprintf("%.2f", result.Duration.Seconds()),
			fmt.Sprintf("%.2f", result.Throughput),
			fmt.Sprintf("%.2f", result.MBPerSecond),
			fmt.Sprintf("%.2f", float64(result.AvgLatency.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.P50Latency.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.P95Latency.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.P99Latency.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.MinLatency.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.MaxLatency.Microseconds())/1000.0),
			strconv.Itoa(result.SuccessCount),
			strconv.Itoa(result.ErrorCount),
			strconv.FormatInt(result.BytesProcessed, 10),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// PrintResults prints benchmark results to console
func PrintResults(result *common.BenchmarkResult) {
	fmt.Println("\n" + "=" * 80)
	fmt.Printf("Benchmark Results: %s\n", result.QueueType)
	fmt.Println("=" * 80)
	fmt.Printf("Messages:           %d\n", result.MessageCount)
	fmt.Printf("Duration:           %v\n", result.Duration)
	fmt.Printf("Throughput:         %.2f msg/s\n", result.Throughput)
	fmt.Printf("Bandwidth:          %.2f MB/s\n", result.MBPerSecond)
	fmt.Printf("Success Count:      %d\n", result.SuccessCount)
	fmt.Printf("Error Count:        %d\n", result.ErrorCount)
	fmt.Printf("Bytes Processed:    %d (%.2f MB)\n", result.BytesProcessed, float64(result.BytesProcessed)/(1024*1024))
	fmt.Println("\nLatency Statistics:")
	fmt.Printf("  Min:              %.2f ms\n", float64(result.MinLatency.Microseconds())/1000.0)
	fmt.Printf("  Avg:              %.2f ms\n", float64(result.AvgLatency.Microseconds())/1000.0)
	fmt.Printf("  P50:              %.2f ms\n", float64(result.P50Latency.Microseconds())/1000.0)
	fmt.Printf("  P95:              %.2f ms\n", float64(result.P95Latency.Microseconds())/1000.0)
	fmt.Printf("  P99:              %.2f ms\n", float64(result.P99Latency.Microseconds())/1000.0)
	fmt.Printf("  Max:              %.2f ms\n", float64(result.MaxLatency.Microseconds())/1000.0)
	fmt.Println("=" * 80 + "\n")
}

// CompareResults prints a comparison of multiple benchmark results
func CompareResults(results []*common.BenchmarkResult) {
	fmt.Println("\n" + "=" * 100)
	fmt.Println("Benchmark Comparison")
	fmt.Println("=" * 100)
	fmt.Printf("%-25s %-15s %-15s %-15s %-15s\n", "Queue Type", "Throughput", "MB/s", "Avg Latency", "P99 Latency")
	fmt.Println("-" * 100)

	for _, result := range results {
		fmt.Printf("%-25s %-15.2f %-15.2f %-15.2f %-15.2f\n",
			result.QueueType,
			result.Throughput,
			result.MBPerSecond,
			float64(result.AvgLatency.Microseconds())/1000.0,
			float64(result.P99Latency.Microseconds())/1000.0,
		)
	}

	fmt.Println("=" * 100 + "\n")
}

// GenerateReport generates a comprehensive benchmark report
func GenerateReport(results []*common.BenchmarkResult, outputDir string) error {
	timestamp := time.Now().Format("20060102-150405")

	// Export to JSON
	jsonFile := fmt.Sprintf("%s/benchmark-results-%s.json", outputDir, timestamp)
	if err := ExportToJSON(results, jsonFile); err != nil {
		return err
	}
	fmt.Printf("JSON report saved to: %s\n", jsonFile)

	// Export to CSV
	csvFile := fmt.Sprintf("%s/benchmark-results-%s.csv", outputDir, timestamp)
	if err := ExportToCSV(results, csvFile); err != nil {
		return err
	}
	fmt.Printf("CSV report saved to: %s\n", csvFile)

	return nil
}
