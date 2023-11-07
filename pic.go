package main

import (
	"bufio"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

func coverPic(srcFile, coverFile, resultFile string, X, Y int) {
	baseX := X
	baseY := Y
	// 读取两个图片源文件
	srcPicImage, srcType, err := readPic(srcFile)
	if err != nil {
		println("read srcPicImage fail: " + err.Error())
		os.Exit(0)
	}
	coverPicImage, _, err := readPic(coverFile)
	if err != nil {
		println("read qd code fail: " + err.Error())
		os.Exit(0)
	}
	// 校验大小
	srcMax := srcPicImage.Bounds().Max
	coverMax := coverPicImage.Bounds().Max
	if coverMax.X+baseX > srcMax.X || coverMax.Y+baseY > srcMax.Y {
		println("x or y over max")
		os.Exit(0)
	}

	// 创建目标image
	var resultImg draw.Image
	switch srcType { // 使用底图后缀
	case "jpeg", "jpg":
		resultFile = resultFile + jpgSuffix
		resultImg = image.NewCMYK(image.Rect(0, 0, srcMax.X, srcMax.Y)) // jpg使用cmyk色彩模式
	default:
		resultFile = resultFile + pngSuffix
		resultImg = image.NewRGBA64(image.Rect(0, 0, srcMax.X, srcMax.Y))
	}

	draw.Draw(resultImg, srcPicImage.Bounds(), srcPicImage, image.Point{}, draw.Src)
	start := image.Pt(X, Y)
	draw.Draw(resultImg, image.Rectangle{start, start.Add(coverPicImage.Bounds().Size())},
		coverPicImage, image.Point{}, draw.Over)

	outFile, err := os.Create(resultFile)
	defer outFile.Close()
	if err != nil {
		println("create result png fail: " + err.Error())
		os.Exit(0)
	}
	b := bufio.NewWriter(outFile)
	switch srcType { // 使用底图格式
	case "jpeg", "jpg":
		err = jpeg.Encode(b, resultImg, &jpeg.Options{100}) // 转cmyk色彩模式
	default:
		err = png.Encode(b, resultImg)
	}
	if err != nil {
		println("encode result file fail: " + err.Error())
		os.Exit(0)
	}
	err = b.Flush()
	if err != nil {
		println("flush result file fail: " + err.Error())
		os.Exit(0)
	}
}

//func convertCMYK(src draw.Image) draw.Image {
//	bound := src.Bounds()
//	result := image.NewCMYK(image.Rect(0, 0, bound.Max.X, bound.Max.Y))
//	for i := 0; i < bound.Max.X; i++ {
//		for j := 0; j < bound.Max.Y; j++ {
//			r, g, b, a := src.At(i, j).RGBA()
//			result.Set(i, j, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
//		}
//	}
//	return result
//}

func readPic(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(f)
}
