package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Data struct {
	Seq       int    `json:"seq"`
	Code      int    `json:"code"`
	Latency   int    `json:"latency"`
	Timestamp string `json:"timestamp"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <datafile>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var data []Data
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var d Data
		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}
		data = append(data, d)
	}

	successCount := 0
	statusCodeFrequency := make(map[int]int)
	latencies := make([]int, len(data))

	for i, d := range data {
			if d.Code >= 200 && d.Code < 300 {
				successCount++
			}
			statusCodeFrequency[d.Code]++
			latencies[i] = d.Latency
		}

	averageLatency := calculateAverageLatency(latencies)
	minLatency := calculateMinLatency(latencies)
	maxLatency := calculateMaxLatency(latencies)
	p99Latency := calculateP99Latency(latencies)

	successRate := float64(successCount) / float64(len(data)) * 100
	fmt.Printf("Success Rate: %.2f%%\n", successRate)
	fmt.Printf("Average Latency: %.2f ms\n", float64(averageLatency)/1e6)
	fmt.Printf("Minimum Latency: %.2f ms\n", float64(minLatency)/1e6)
	fmt.Printf("Maximum Latency: %.2f ms\n", float64(maxLatency)/1e6)
	fmt.Printf("p99 Latency: %.2f ms\n", float64(p99Latency)/1e6)

	generateLineChart(latencies, "Latency Over Time", "latency.png")
	generateBarChart(statusCodeFrequency, "Status Codes Over Time", "status_codes.png")
}

func calculateAverageLatency(latencies []int) int {
	totalLatency := 0
	for _, latency := range latencies {
		totalLatency += latency
	}
	return totalLatency / len(latencies)
}

func calculateMinLatency(latencies []int) int {
	if len(latencies) == 0 {
		return 0
	}
	min := latencies[0]
	for _, latency := range latencies {
		if latency < min {
			min = latency
		}
	}
	return min
}

func calculateMaxLatency(latencies []int) int {
	if len(latencies) == 0 {
		return 0
	}
	max := latencies[0]
	for _, latency := range latencies {
		if latency > max {
			max = latency
		}
	}
	return max
}

func calculateP99Latency(latencies []int) int {
	if len(latencies) == 0 {
		return 0
	}
	sort.Ints(latencies)
	p99Index := int(float64(len(latencies)) * 0.99)
	return latencies[p99Index]
}

func generateLineChart(latencies []int, title, filename string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "Request Number"
	p.Y.Label.Text = "Latency (ms)"

	pts := make(plotter.XYs, len(latencies))
	for i, latency := range latencies {
		pts[i].X = float64(i)
		pts[i].Y = float64(latency) / 1e6 // Convert latency to milliseconds
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("Failed to create line plotter: %v", err)
	}
	line.Color = color.RGBA{B: 255, A: 255}
	p.Add(line)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}
}

func generateBarChart(frequency map[int]int, title, filename string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "Status Code"
	p.Y.Label.Text = "Count"

	var codes []int
	for code := range frequency {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	var counts plotter.Values
	for _, code := range codes {
		counts = append(counts, float64(frequency[code]))
	}

	bars, err := plotter.NewBarChart(counts, vg.Points(20))
	if err != nil {
		log.Fatalf("Failed to create bar chart: %v", err)
	}
	bars.Color = color.RGBA{R: 255, A: 255}
	p.Add(bars)

	// Label X-axis with status code strings
	var codeLabels []string
	for _, code := range codes {
		codeLabels = append(codeLabels, fmt.Sprintf("%d", code))
	}
	p.NominalX(codeLabels...)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}
}

