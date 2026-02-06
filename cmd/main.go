// Bu fayl SDK dan foydalanish namunasi
// go run cmd/main.go bilan ishga tushiring

package main

import (
	"fmt"
	"os"

	pdfsdk "github.com/infosec554/golang-pdf-sdk"
)

func main() {
	fmt.Println("🚀 Golang PDF SDK - Namunalar")
	fmt.Println("================================")

	// SDK ni ishga tushirish (Gotenberg URL)
	pdf := pdfsdk.New("http://localhost:3000")

	// Namuna 1: Kompressiya
	fmt.Println("\n📦 1. PDF Kompressiya:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Compress().CompressBytes(input)
		if err != nil {
			fmt.Println("   ❌ Xato:", err)
		} else {
			os.WriteFile("compressed.pdf", output, 0644)
			saving := 100 - (float64(len(output))/float64(len(input)))*100
			fmt.Printf("   ✅ %d → %d bayt (%.1f%% tejaldi)\n", len(input), len(output), saving)
		}
	} else {
		fmt.Println("   ⚠️  test.pdf topilmadi, o'tkazib yuborildi")
	}

	// Namuna 2: Aylantirish
	fmt.Println("\n🔄 2. PDF Aylantirish (90°):")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Rotate().RotateBytes(input, 90, "all")
		if err != nil {
			fmt.Println("   ❌ Xato:", err)
		} else {
			os.WriteFile("rotated.pdf", output, 0644)
			fmt.Println("   ✅ rotated.pdf yaratildi")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf topilmadi")
	}

	// Namuna 3: Watermark
	fmt.Println("\n💧 3. Watermark qo'shish:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Watermark().AddWatermarkBytes(input, "MAXFIY HUJJAT", nil)
		if err != nil {
			fmt.Println("   ❌ Xato:", err)
		} else {
			os.WriteFile("watermarked.pdf", output, 0644)
			fmt.Println("   ✅ watermarked.pdf yaratildi")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf topilmadi")
	}

	// Namuna 4: Himoyalash
	fmt.Println("\n🔒 4. PDF Himoyalash:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		output, err := pdf.Protect().ProtectBytes(input, "parol123")
		if err != nil {
			fmt.Println("   ❌ Xato:", err)
		} else {
			os.WriteFile("protected.pdf", output, 0644)
			fmt.Println("   ✅ protected.pdf yaratildi (parol: parol123)")
		}
	} else {
		fmt.Println("   ⚠️  test.pdf topilmadi")
	}

	// Namuna 5: PDF dan JPG
	fmt.Println("\n🖼️  5. PDF dan JPG:")
	if input, err := os.ReadFile("test.pdf"); err == nil {
		images, err := pdf.PDFToJPG().ConvertToImages(input)
		if err != nil {
			fmt.Println("   ❌ Xato:", err)
		} else {
			os.MkdirAll("pages", 0755)
			for i, img := range images {
				os.WriteFile(fmt.Sprintf("pages/sahifa_%d.jpg", i+1), img, 0644)
			}
			fmt.Printf("   ✅ %d ta rasm yaratildi (pages/ papkada)\n", len(images))
		}
	} else {
		fmt.Println("   ⚠️  test.pdf topilmadi")
	}

	fmt.Println("\n================================")
	fmt.Println("✅ Namunalar tugadi!")
	fmt.Println("\n💡 Sinab ko'rish uchun 'test.pdf' faylini qo'ying")
}
