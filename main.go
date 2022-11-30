package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	pngSuffix = ".png"
	PathSep   = string(os.PathSeparator)
)

func main() {
	rootDir := getCurrentAbPathByExecutable()
	// 配置读取
	content, err := os.ReadFile(rootDir + PathSep + "coverPic.conf")
	if err != nil {
		println("打开配置文件失败: " + err.Error())
		return
	}
	confArr := strings.Split(string(content), "\n")
	if len(confArr) < 5 {
		println("配置文件coverPic.conf内容行数不对, 一共5行, srcDir, qrCodeDir, resultDir, X, Y")
		return
	}
	srcDir := confArr[0]
	if strings.HasPrefix(srcDir, "."+PathSep) {
		srcDir = strings.Replace(srcDir, "."+PathSep, rootDir+PathSep, 1)
	}
	qrCodeDir := confArr[1]
	if strings.HasPrefix(qrCodeDir, "."+PathSep) {
		qrCodeDir = strings.Replace(qrCodeDir, "."+PathSep, rootDir+PathSep, 1)
	}
	resultDir := confArr[2]
	if strings.HasPrefix(resultDir, "."+PathSep) {
		resultDir = strings.Replace(resultDir, "."+PathSep, rootDir+PathSep, 1)
	}
	X, err := strconv.ParseInt(confArr[3], 10, 64)
	if err != nil {
		println("配置X坐标值错误, 必须是个整数: " + confArr[3])
		return
	}
	Y, err := strconv.ParseInt(confArr[4], 10, 64)
	if err != nil {
		println("配置Y坐标值错误, 必须是个整数: " + confArr[3])
		return
	}

	// 文件夹检测
	srcList, err := ListDir(srcDir, pngSuffix)
	if err != nil {
		println("读取源目录失败: " + err.Error())
		return
	}
	if len(srcList) <= 0 {
		println("确定这个模板图目录下有.png文件吗: " + srcDir)
		return
	}

	qrCodeList, err := ListDir(qrCodeDir, pngSuffix)
	if err != nil {
		println("读取源目录失败: " + err.Error())
		return
	}
	if len(qrCodeList) <= 0 {
		println("确定这个二维码目录下有.png文件吗: " + srcDir)
		return
	}
	err = os.Mkdir(resultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		println("创建结果目录失败: " + err.Error())
		return
	}
	count := 0
	sTime:=time.Now()
	for _, srcFile := range srcList {
		templateDir := strings.Replace(srcFile, pngSuffix, "", -1)
		templateDir = resultDir + PathSep + templateDir
		err = os.Mkdir(templateDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			println("创建模板目录失败: " + err.Error())
			return
		}
		for _, qrCodeFile := range qrCodeList {
			coverPic(srcDir+PathSep+srcFile,
				qrCodeDir+PathSep+qrCodeFile,
				templateDir+PathSep+qrCodeFile,
				int(X), int(Y))
			count++
			println(fmt.Sprintf("完成: %s/%s 总耗时: %d秒 进度: %d",
				srcFile, qrCodeFile, time.Now().Sub(sTime).Milliseconds()/1000, count*100/len(qrCodeList)*len(srcList)))
		}
	}
}

// 获取当前执行程序所在的绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// ListDir 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		println("read dir: " + dirPth + " fail: " + err.Error())
		return nil, err
	}
	suffix = strings.ToLower(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() || !strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
			continue
		}
		files = append(files, fi.Name())
	}
	return files, nil
}

func coverPic(srcFile, coverFile, resultFile string, X, Y int) {
	baseX := X
	baseY := Y
	srcPNG, err := readPNG(srcFile)
	if err != nil {
		println("read srcPNG fail: " + err.Error())
		return
	}
	qrCode, err := readPNG(coverFile)
	if err != nil {
		println("read qd code fail: " + err.Error())
		return
	}
	srcMax := srcPNG.Bounds().Max
	qrMax := qrCode.Bounds().Max
	if qrMax.X+baseX > srcMax.X || qrMax.Y+baseY > srcMax.Y {
		println("x or y over max")
		return
	}
	//println("src size: " + srcMax.String())
	resultImg := image.NewRGBA(image.Rect(0, 0, srcMax.X, srcMax.Y))
	// 复制srcPNG
	for x := 0; x < srcMax.X; x++ {
		for y := 0; y < srcMax.Y; y++ {
			p := srcPNG.At(x, y)
			resultImg.Set(x, y, p)
		}
	}
	// 覆盖二维码
	for x := 0; x < qrMax.X; x++ {
		for y := 0; y < qrMax.Y; y++ {
			resultImg.Set(x+baseX, y+baseY, qrCode.At(x, y))
		}
	}
	outFile, err := os.Create(resultFile)
	defer outFile.Close()
	if err != nil {
		println("create result png fail: " + err.Error())
		return
	}
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, resultImg)
	if err != nil {
		println("encode result file fail: " + err.Error())
		return
	}
	err = b.Flush()
	if err != nil {
		println("flush result file fail: " + err.Error())
		return
	}
}

func readPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}
