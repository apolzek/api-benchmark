package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"

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

	// Analyze Data
	successCount := 0
	statusCodeFrequency := make(map[int]int)
	timeSeries := make(plotter.XYs, len(data))
	for i, d := range data {
		if d.Code >= 200 && d.Code < 300 {
			successCount++
		}
		statusCodeFrequency[d.Code]++
		timeSeries[i].X = float64(i) // Use index for X-axis
		timeSeries[i].Y = float64(d.Latency) / 1e6 // Convert latency to milliseconds
	}

	successRate := float64(successCount) / float64(len(data)) * 100
	fmt.Printf("Success Rate: %.2f%%\n", successRate)

	// Generate Latency Chart
	generateLineChart(timeSeries, "Latency Over Time", "latency.png")

	// Generate Status Code Chart
	generateBarChart(statusCodeFrequency, "Status Codes Over Time", "status_codes.png")
}

func generateLineChart(points plotter.XYs, title, filename string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Request Number"
	p.Y.Label.Text = "Latency (ms)"

	line, err := plotter.NewLine(points)
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

	var statusCodeLabels []string
	var statusCodeCounts plotter.Values

	for code, count := range frequency {
		statusCodeLabels = append(statusCodeLabels, fmt.Sprintf("%d", code))
		statusCodeCounts = append(statusCodeCounts, float64(count))
	}

	bars, err := plotter.NewBarChart(statusCodeCounts, vg.Points(20))
	if err != nil {
		log.Fatalf("Failed to create bar chart: %v", err)
	}

	bars.Color = color.RGBA{R: 255, A: 255}

	p.Add(bars)
	p.NominalX(statusCodeLabels...) // Label X-axis with status code strings

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}
}

