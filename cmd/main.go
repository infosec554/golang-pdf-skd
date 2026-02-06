
package main

import (
	"fmt"
	"os"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
	fmt.Println("🚀 Golang PDF SDK - Examples")
	fmt.Println("==============================")

	pdf := pdfsdk.New("http://localhost:3000")

	fmt.Println("\n📦 1. PDF Compression:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Compress().CompressBytes(input)
		if err != nil {
			fmt.Println("   ❌ Error:", err)
		} else {
			os.WriteFile("compressed.pdf", output, 0644)
			saving := 100 - (float64(len(output))/float64(len(input)))*100
			fmt.Printf("   ✅ %d → %d bytes (%.1f%% saved)\n", len(input), len(output), saving)
		}
	} else {
		fmt.Println("   ⚠️  test.pdf not found, skipped")
	}

	fmt.Println("\n🔄 2. PDF Rotation (90°):")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Rotate().RotateBytes(input, 90, "all")
		if err != nil {
			fmt.Println("   ❌ Error:", err)
		} else {
			os.WriteFile("rotated.pdf", output, 0644)
			fmt.Println("   ✅ Created rotated.pdf")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf not found")
	}

	fmt.Println("\n💧 3. Add Watermark:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)
		if err != nil {
			fmt.Println("   ❌ Error:", err)
		} else {
			os.WriteFile("watermarked.pdf", output, 0644)
			fmt.Println("   ✅ Created watermarked.pdf")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf not found")
	}

	fmt.Println("\n🔒 4. Protect PDF:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Protect().ProtectBytes(input, "password123")
		if err != nil {
			fmt.Println("   ❌ Error:", err)
		} else {
			os.WriteFile("protected.pdf", output, 0644)
			fmt.Println("   ✅ Created protected.pdf (password: password123)")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf not found")
	}

	fmt.Println("\n🖼️  5. PDF to JPG Images:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		images, err := pdf.PDFToJPG().ConvertToImages(input)
		if err != nil {
			fmt.Println("   ❌ Error:", err)
		} else {
			os.MkdirAll("pages", 0755)
			for i, img := range images {
				os.WriteFile(fmt.Sprintf("pages/page_%d.jpg", i+1), img, 0644)
			}
			fmt.Printf("   ✅ Created %d images in pages/ folder\n", len(images))
		}
	} else {
		fmt.Println("   ⚠️  test.pdf not found")
	}

	fmt.Println("\n==============================")
	fmt.Println("✅ Examples completed!")
	fmt.Println("\n💡 Place a 'test.pdf' file to run all examples")
}
