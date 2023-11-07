package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	pngSuffix = ".png"
	jpgSuffix = ".jpg"
	PathSep   = string(os.PathSeparator)
)

var (
	fileSuffix = []string{pngSuffix, jpgSuffix}

	// 初始化时写入配置
	srcDir    string
	coverDir  string
	resultDir string
	X         int64
	Y         int64

	// 目标文件检测
	srcList   []string
	coverList []string
)

func main() {
	loadConf()              // 读取相对路径配置文件(./coverPic.conf)，加载到全局变量
	getListAndMkTargetDir() // 检查两个原目录文件，创建目标目录

	count := 0
	sTime := time.Now()
	for _, srcFile := range srcList {
		resultItemDir := srcFile
		for _, suffix := range fileSuffix { // 文件名去后缀
			resultItemDir = strings.Replace(srcFile, suffix, "", -1)
		}
		resultItemDir = resultDir + PathSep + resultItemDir
		err := os.Mkdir(resultItemDir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			println("创建模板目录失败: " + err.Error())
			return
		}
		for _, coverFile := range coverList {
			resultFileName := coverFile
			for _, suffix := range fileSuffix { // 文件名去后缀
				resultFileName = strings.Replace(resultFileName, suffix, "", -1)
			}
			coverPic(srcDir+PathSep+srcFile,
				coverDir+PathSep+coverFile,
				resultItemDir+PathSep+resultFileName, // ./${result}/${template}/${cover}
				int(X), int(Y))
			count++
			println(fmt.Sprintf("完成: %s/%s 总耗时: %d秒 进度: %d",
				srcFile, coverFile, time.Now().Sub(sTime).Milliseconds()/1000, count*100/len(coverList)*len(srcList)))
		}
	}
}
