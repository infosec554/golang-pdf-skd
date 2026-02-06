package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
	fmt.Println("🚀 Golang PDF SDK v2.1 - Examples")
	fmt.Println("==================================")

	// Initialize SDK with custom options
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		GotenbergURL:        "http://localhost:3000",
		MaxWorkers:          10,  // Max 10 parallel operations
		MaxIdleConns:        100, // Connection pool size
		MaxIdleConnsPerHost: 10,
		RequestTimeout:      5 * time.Minute,
	})
	defer sdk.Close()

	fmt.Printf("\n📊 SDK Version: %s\n", pdfsdk.Version)
	fmt.Printf("📊 Max Workers: %d\n", sdk.Stats().MaxWorkers)

	// Read test PDF
	input, err := os.ReadFile("test.pdf")
	if err != nil {
		fmt.Println("\n⚠️  test.pdf not found - skipping PDF operations")
		fmt.Println("💡 Place a 'test.pdf' file in this directory to test")
		showAPIExamples()
		return
	}

	fmt.Printf("\n📁 Input file size: %d bytes\n", len(input))

	// Example 1: Get PDF Info
	fmt.Println("\n📋 1. PDF Info:")
	info, err := sdk.Info().GetInfoBytes(input)
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		fmt.Printf("   ✅ Pages: %d, Version: %s, Encrypted: %v\n",
			info.PageCount, info.Version, info.Encrypted)
	}

	// Example 2: Compression
	fmt.Println("\n📦 2. PDF Compression:")
	output, err := sdk.Compress().CompressBytes(input)
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.WriteFile("compressed.pdf", output, 0644)
		saving := 100 - (float64(len(output))/float64(len(input)))*100
		fmt.Printf("   ✅ %d → %d bytes (%.1f%% saved)\n", len(input), len(output), saving)
	}

	// Example 3: Rotation
	fmt.Println("\n🔄 3. PDF Rotation (90°):")
	output, err = sdk.Rotate().RotateBytes(input, 90, "all")
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.WriteFile("rotated.pdf", output, 0644)
		fmt.Println("   ✅ Created rotated.pdf")
	}

	// Example 4: Watermark
	fmt.Println("\n💧 4. Add Watermark:")
	output, err = sdk.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.WriteFile("watermarked.pdf", output, 0644)
		fmt.Println("   ✅ Created watermarked.pdf")
	}

	// Example 5: Protection
	fmt.Println("\n🔒 5. Protect PDF:")
	output, err = sdk.Protect().ProtectBytes(input, "password123")
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.WriteFile("protected.pdf", output, 0644)
		fmt.Println("   ✅ Created protected.pdf (password: password123)")
	}

	// Example 6: Pipeline (chained operations)
	fmt.Println("\n⛓️  6. Pipeline (Compress → Watermark → Protect):")
	output, err = sdk.Pipeline().
		Compress().
		Watermark("DRAFT", nil).
		Protect("secret").
		Execute(input)
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.WriteFile("pipeline_result.pdf", output, 0644)
		fmt.Println("   ✅ Created pipeline_result.pdf")
	}

	// Example 7: Batch Processing
	fmt.Println("\n🔄 7. Batch Processing (parallel compression):")
	inputs := [][]byte{input, input, input} // 3 copies for demo
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	results := sdk.Batch(3).CompressBatch(ctx, inputs)
	successCount := 0
	for _, r := range results {
		if r.Error == nil {
			successCount++
		}
	}
	fmt.Printf("   ✅ Processed %d/%d PDFs in parallel\n", successCount, len(inputs))

	// Example 8: PDF to JPG
	fmt.Println("\n🖼️  8. PDF to JPG Images:")
	images, err := sdk.PDFToJPG().ConvertToImages(input)
	if err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		os.MkdirAll("pages", 0755)
		for i, img := range images {
			os.WriteFile(fmt.Sprintf("pages/page_%d.jpg", i+1), img, 0644)
		}
		fmt.Printf("   ✅ Created %d images in pages/ folder\n", len(images))
	}

	// Show stats
	stats := sdk.Stats()
	fmt.Println("\n==================================")
	fmt.Printf("📊 SDK Stats: Processed=%d, Active=%d/%d workers\n",
		stats.ProcessedTasks, stats.ActiveWorkers, stats.MaxWorkers)
	fmt.Println("✅ All examples completed!")
}

func showAPIExamples() {
	fmt.Print("\n📚 API Usage Examples:")
	fmt.Print(`
// Initialize with options
sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
    GotenbergURL: "http://localhost:3000",
    MaxWorkers:   10,
})
defer sdk.Close()

// Get PDF info
info, _ := sdk.Info().GetInfoBytes(pdfBytes)
fmt.Println("Pages:", info.PageCount)

// Compress
compressed, _ := sdk.Compress().CompressBytes(input)

// Pipeline (chain operations)
result, _ := sdk.Pipeline().
    Compress().
    Watermark("CONFIDENTIAL", nil).
    Protect("password").
    Execute(input)

// Batch processing (parallel)
ctx := context.Background()
results := sdk.Batch(5).CompressBatch(ctx, [][]byte{pdf1, pdf2, pdf3})

// Worker pool control
sdk.Workers().Acquire()  // Get worker slot
defer sdk.Workers().Release()
// ... do work ...
`)
}
