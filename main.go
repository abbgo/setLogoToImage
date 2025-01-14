package main

import (
	"image"
	_ "image/gif"
	"image/jpeg"  // PNG desteği
	_ "image/png" // PNG desteği ekleniyor
	"log"
	"os"

	"golang.org/x/image/draw"
)

func main() {
	// Orijinal resmi aç
	srcImage, err := os.Open("images.png")
	if err != nil {
		log.Fatalf("Orijinal resim açılamadı: %v", err)
	}
	defer srcImage.Close()

	// Orijinal resmi decode et
	mainImg, _, err := image.Decode(srcImage)
	if err != nil {
		log.Fatalf("Orijinal resim decode edilemedi: %v", err)
	}

	// Logo resmini aç
	logoImage, err := os.Open("logo.png")
	if err != nil {
		log.Fatalf("Logo resmi açılamadı: %v", err)
	}
	defer logoImage.Close()

	// Logo resmini decode et
	logoImg, _, err := image.Decode(logoImage)
	if err != nil {
		log.Fatalf("Logo resmi decode edilemedi: %v", err)
	}

	// Ana resmin boyutlarını al
	mainBounds := mainImg.Bounds()
	mainWidth := mainBounds.Dx()
	mainHeight := mainBounds.Dy()

	// Logo için yeni boyutlar belirle
	newLogoWidth := mainWidth / 12 // Ana resmin genişliğinin 1/5'i kadar
	newLogoHeight := newLogoWidth  // Kare olarak ölçekle
	resizedLogo := image.NewRGBA(image.Rect(0, 0, newLogoWidth, newLogoHeight))

	// Logoyu yeniden boyutlandır
	draw.CatmullRom.Scale(resizedLogo, resizedLogo.Bounds(), logoImg, logoImg.Bounds(), draw.Over, nil)

	// Yeni bir RGBA görüntü oluştur
	outputImg := image.NewRGBA(mainBounds)

	// Orijinal resmi kopyala
	draw.Draw(outputImg, mainBounds, mainImg, image.Point{0, 0}, draw.Src)

	// Yeniden boyutlandırılmış logoyu alt sağ köşeye yerleştir
	offsetX := mainWidth - newLogoWidth - 10   // Sağ kenardan 10 piksel içeride
	offsetY := mainHeight - newLogoHeight - 10 // Alt kenardan 10 piksel yukarıda
	logoPosition := image.Pt(offsetX, offsetY)

	// Yeniden boyutlandırılmış logoyu ana resme çiz
	draw.Draw(outputImg, resizedLogo.Bounds().Add(logoPosition), resizedLogo, image.Point{0, 0}, draw.Over)

	// Yeni dosyayı oluştur ve kaydet
	outputFile, err := os.Create("output_with_logo.jpg")
	if err != nil {
		log.Fatalf("Çıktı dosyası oluşturulamadı: %v", err)
	}
	defer outputFile.Close()

	// Çıktıyı JPEG olarak kaydet
	err = jpeg.Encode(outputFile, outputImg, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatalf("Çıktı dosyası kaydedilemedi: %v", err)
	}

	log.Println("Logo başarıyla eklendi ve 'output_with_logo.jpg' dosyasına kaydedildi.")
}
