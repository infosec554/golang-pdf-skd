# 📄 Golang PDF SDK

Go dasturlash tilida PDF fayllari bilan ishlash uchun kuchli va oson SDK.

## 🚀 O'rnatish

```bash
go get github.com/infosec554/golang-pdf-sdk
```

## 📋 Talablar

- **Go 1.21+**
- **Gotenberg** (Word/Excel/PowerPoint konvertatsiya uchun) - [gotenberg.dev](https://gotenberg.dev)
- **pdftoppm** (PDF dan JPG uchun) - poppler-utils paketi

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# macOS
brew install poppler

# Gotenberg docker orqali
docker run -d -p 3000:3000 gotenberg/gotenberg:8
```

## 💡 Tez Boshlash

```go
package main

import (
    "fmt"
    "os"
    
    pdfsdk "github.com/infosec554/golang-pdf-sdk"
)

func main() {
    // SDK ni ishga tushirish
    pdf := pdfsdk.New("http://localhost:3000")
    
    // PDF faylni o'qish
    input, _ := os.ReadFile("fayl.pdf")
    
    // Kompressiya qilish
    output, err := pdf.Compress().CompressBytes(input)
    if err != nil {
        panic(err)
    }
    
    // Natijani saqlash
    os.WriteFile("kichik.pdf", output, 0644)
    
    fmt.Printf("Hajmi: %d -> %d bayt\n", len(input), len(output))
}
```

---

## 📚 Barcha Funksiyalar

### 1️⃣ PDF Kompressiya

PDF hajmini kamaytirish:

```go
pdf := pdfsdk.New("http://localhost:3000")

// Baytlar bilan ishlash
input, _ := os.ReadFile("katta.pdf")
output, err := pdf.Compress().CompressBytes(input)

// Fayl bilan ishlash
err := pdf.Compress().CompressFile("input.pdf", "output.pdf")
```

---

### 2️⃣ PDF Birlashtirish (Merge)

Bir nechta PDF larni birlashtirib bitta PDF qilish:

```go
pdf := pdfsdk.New("http://localhost:3000")

// Baytlar bilan
pdf1, _ := os.ReadFile("1.pdf")
pdf2, _ := os.ReadFile("2.pdf")
pdf3, _ := os.ReadFile("3.pdf")

output, err := pdf.Merge().MergeBytes([][]byte{pdf1, pdf2, pdf3})
os.WriteFile("birlashgan.pdf", output, 0644)

// Fayllar bilan
err := pdf.Merge().MergeFiles(
    []string{"1.pdf", "2.pdf", "3.pdf"}, 
    "birlashgan.pdf",
)
```

---

### 3️⃣ PDF Bo'lish (Split)

PDF ni sahifalar bo'yicha bo'lish:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("kitob.pdf")

// Sahifa diapazonlari bo'yicha (ZIP qaytaradi)
zipBytes, err := pdf.Split().SplitBytes(input, "1-5,6-10,11-20")
os.WriteFile("qismlar.zip", zipBytes, 0644)

// Har bir sahifani alohida olish
pages, err := pdf.Split().SplitToPages(input)
for i, page := range pages {
    os.WriteFile(fmt.Sprintf("sahifa_%d.pdf", i+1), page, 0644)
}
```

---

### 4️⃣ PDF Aylantirish (Rotate)

Sahifalarni 90°, 180°, 270° ga aylantirish:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("fayl.pdf")

// Barcha sahifalarni 90° ga aylantirish
output, err := pdf.Rotate().RotateBytes(input, 90, "all")

// Faqat 1-3 sahifalarni 180° ga
output, err := pdf.Rotate().RotateBytes(input, 180, "1-3")

// Fayl bilan
err := pdf.Rotate().RotateFile("input.pdf", "output.pdf", 270, "all")
```

---

### 5️⃣ Watermark (Suv belgisi)

PDF ga matn watermark qo'shish:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("hujjat.pdf")

// Oddiy watermark
output, err := pdf.Watermark().AddWatermarkBytes(input, "MAXFIY", nil)

// Sozlamalar bilan
options := &service.WatermarkOptions{
    FontSize: 72,           // Shrift o'lchami
    Position: "diagonal",   // "diagonal", "center"
    Opacity:  0.3,          // Shaffoflik (0.0 - 1.0)
    Color:    "red",        // Rang
}
output, err := pdf.Watermark().AddWatermarkBytes(input, "QORALAMA", options)
```

---

### 6️⃣ PDF Himoyalash (Parol qo'yish)

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("maxfiy.pdf")

// Parol qo'yish
protected, err := pdf.Protect().ProtectBytes(input, "parol123")
os.WriteFile("himoyalangan.pdf", protected, 0644)

// Fayl bilan
err := pdf.Protect().ProtectFile("input.pdf", "protected.pdf", "parol123")
```

---

### 7️⃣ PDF Qulfini Ochish

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("himoyalangan.pdf")

// Parolni olib tashlash
unlocked, err := pdf.Unlock().UnlockBytes(input, "parol123")
os.WriteFile("ochilgan.pdf", unlocked, 0644)
```

---

### 8️⃣ PDF dan Rasm (JPG)

PDF sahifalarini JPG rasmlarga aylantirish:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("prezentatsiya.pdf")

// ZIP fayl olish (barcha rasmlar ichida)
zipBytes, err := pdf.PDFToJPG().ConvertBytes(input)
os.WriteFile("rasmlar.zip", zipBytes, 0644)

// Har bir rasmni alohida olish
images, err := pdf.PDFToJPG().ConvertToImages(input)
for i, img := range images {
    os.WriteFile(fmt.Sprintf("sahifa_%d.jpg", i+1), img, 0644)
}
```

---

### 9️⃣ Rasmlardan PDF

JPG/PNG rasmlarni bitta PDF ga birlashtirish:

```go
pdf := pdfsdk.New("http://localhost:3000")

// Rasm fayllarini o'qish
img1, _ := os.ReadFile("rasm1.jpg")
img2, _ := os.ReadFile("rasm2.png")
img3, _ := os.ReadFile("rasm3.jpg")

// PDF yaratish
pdfBytes, err := pdf.JPGToPDF().ConvertMultipleBytes(
    [][]byte{img1, img2, img3},
    []string{"rasm1.jpg", "rasm2.png", "rasm3.jpg"},
)
os.WriteFile("rasmlar.pdf", pdfBytes, 0644)

// Fayllar bilan
err := pdf.JPGToPDF().ConvertFiles(
    []string{"1.jpg", "2.jpg", "3.png"},
    "albom.pdf",
)
```

---

### 🔟 Word/Excel/PowerPoint → PDF

**⚠️ Gotenberg kerak!**

```go
import "context"

pdf := pdfsdk.New("http://localhost:3000")
ctx := context.Background()

// Word → PDF
docxBytes, _ := os.ReadFile("hujjat.docx")
pdfBytes, err := pdf.WordToPDF().ConvertBytes(ctx, docxBytes, "hujjat.docx")

// Excel → PDF  
xlsxBytes, _ := os.ReadFile("jadval.xlsx")
pdfBytes, err := pdf.ExcelToPDF().ConvertBytes(ctx, xlsxBytes, "jadval.xlsx")

// PowerPoint → PDF
pptxBytes, _ := os.ReadFile("slayd.pptx")
pdfBytes, err := pdf.PowerPointToPDF().ConvertBytes(ctx, pptxBytes, "slayd.pptx")
```

---

## 🔧 Sozlamalar

`.env` fayl orqali:

```env
GOTENBERG_URL=http://localhost:3000
SERVICE_NAME=my-pdf-app
LOGGER_LEVEL=info
```

---

## 📦 Loyiha Strukturasi

```
golang-pdf-sdk/
├── pdfsdk.go           # Asosiy SDK entry point
├── config/             # Konfiguratsiya
├── pkg/
│   ├── gotenberg/      # Gotenberg client
│   └── logger/         # Logger
└── service/
    ├── compress.go     # PDF kompressiya
    ├── merge.go        # PDF birlashtirish
    ├── split.go        # PDF bo'lish
    ├── rotate.go       # PDF aylantirish
    ├── watermark.go    # Watermark qo'shish
    ├── protect.go      # Parol qo'yish
    ├── unlock.go       # Parolni ochish
    ├── pdf_to_jpg.go   # PDF → JPG
    ├── jpgtopdf.go     # JPG → PDF
    ├── word_to_pdf.go  # Word → PDF
    ├── excel_to_pdf.go # Excel → PDF
    └── powerpoint.go   # PPT → PDF
```

---

## 📄 Litsenziya

MIT License - batafsil [LICENSE](LICENSE) ga qarang.

## 🤝 Hissa Qo'shish

Pull Request lar qabul qilinadi!
