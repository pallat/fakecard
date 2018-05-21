package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type changeable interface {
	Set(x, y int, c color.Color)
	image.Image
}

func main() {
	imgfile, err := os.Open("conner_id_card.png")
	if err != nil {
		panic(err.Error())
	}
	defer imgfile.Close()

	img, err := png.Decode(imgfile)
	if err != nil {
		panic(err.Error())
	}

	newImg := mix(img)

	var cimg changeable
	var ok bool
	if cimg, ok = newImg.(changeable); ok {
		colBlk := color.RGBA{20, 20, 20, 255}
		colBlue := color.RGBA{20, 20, 220, 255}
		colRed := color.RGBA{220, 20, 20, 255}

		addLabel(cimg, 80, 30, "บัตรประจำตัวประชาชน", colBlk, 18)
		addLabel(cimg, 260, 30, "Thai National ID Card", colBlue, 18)

		addLabel(cimg, 80, 50, "เลขประจำตัวประชาชน", colBlk, 12)
		addLabel(cimg, 80, 65, "Identification Number", colBlue, 12)
		addLabel(cimg, 220, 55, "1 2345 67890 12 1", colBlk, 18)
		addLabel(cimg, 60, 90, "ชื่อตัวและชื่อสกุล", colBlk, 12)
		addLabel(cimg, 170, 90, "เจ้าพระยา ประยุทธ จานโอชา", colBlk, 18)

		addLabel(cimg, 140, 110, "Name    Mr. Prayuth", colBlue, 12)
		addLabel(cimg, 140, 130, "Last Name   ChanOCha", colBlue, 12)
		addLabel(cimg, 160, 150, "เกิดวันที่  19 ต.ค. 2499", colBlk, 12)
		addLabel(cimg, 160, 170, "Date of Birth    19 Oct. 1956", colBlue, 12)

		addLabel(cimg, 50, 220, "ที่อยู่ 1/11 ราบสิบเอ็ด เขตบางบัว", colBlk, 12)
		addLabel(cimg, 50, 240, "แขวงบางเขน จ.กรุงเทพมหานคร", colBlk, 12)

		addLabel(cimg, 50, 260, "18 ส.ค. 2560", colBlk, 11)
		addLabel(cimg, 50, 275, "วันออกบัตร", colBlk, 11)
		addLabel(cimg, 50, 290, "18 Aug. 2017", colBlue, 11)
		addLabel(cimg, 50, 305, "Date of Issue", colBlue, 11)

		addLabel(cimg, 280, 260, "18 ส.ค. 2560", colBlk, 11)
		addLabel(cimg, 280, 275, "วันบัตรหมดอายุ", colBlk, 11)
		addLabel(cimg, 280, 290, "18 Aug. 2017", colBlue, 11)
		addLabel(cimg, 280, 305, "Date of Expiry", colBlue, 11)

		addLabel(cimg, 370, 295, "10070206191152", colRed, 11)
	} else {
		fmt.Println("no luck")
	}

	f, err := os.Create("id_card.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if err := png.Encode(f, cimg); err != nil {
		panic(err)
	}

}

func addLabel(img draw.Image, x, y int, label string, col color.RGBA, size float64) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: fontface(size),
		Dot:  point,
	}
	d.DrawString(label)
}

func fontface(size float64) font.Face {
	ttf, _ := ioutil.ReadFile("dtac2013_bl.ttf")
	f, _ := truetype.Parse(ttf)
	return truetype.NewFace(f, &truetype.Options{
		Size:              size,
		DPI:               72,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	})
}

func mix(background image.Image) image.Image {
	image2, err := os.Open("images.jpeg")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	second, err := jpeg.Decode(image2)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image2.Close()

	//resize
	n_second := resize.Resize(110, 135, second, resize.NearestNeighbor)
	//resize

	offset := image.Pt(360, 150)
	b := background.Bounds()
	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, background, image.ZP, draw.Src)
	draw.Draw(image3, n_second.Bounds().Add(offset), n_second, image.ZP, draw.Over)

	return image3
}
