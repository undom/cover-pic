package main

import (
	"os"
	"strconv"
	"strings"
)

func loadConf() {
	rootDir := getCurrentAbPathByExecutable()
	// 配置读取
	content, err := os.ReadFile(rootDir + PathSep + "coverPic.conf")
	if err != nil {
		println("打开配置文件失败: " + err.Error())
		os.Exit(0)
	}
	confArr := strings.Split(string(content), "\n")
	if len(confArr) < 5 {
		println("配置文件coverPic.conf内容行数不对, 一共5行, srcDir, coverDir, resultDir, X, Y")
		os.Exit(0)
	}
	srcDir = confArr[0]
	if strings.HasPrefix(srcDir, "."+PathSep) {
		srcDir = strings.Replace(srcDir, "."+PathSep, rootDir+PathSep, 1)
	}
	coverDir = confArr[1]
	if strings.HasPrefix(coverDir, "."+PathSep) {
		coverDir = strings.Replace(coverDir, "."+PathSep, rootDir+PathSep, 1)
	}
	resultDir = confArr[2]
	if strings.HasPrefix(resultDir, "."+PathSep) {
		resultDir = strings.Replace(resultDir, "."+PathSep, rootDir+PathSep, 1)
	}
	X, err = strconv.ParseInt(confArr[3], 10, 64)
	if err != nil {
		println("配置X坐标值错误, 必须是个整数: " + confArr[3])
		os.Exit(0)
	}
	Y, err = strconv.ParseInt(confArr[4], 10, 64)
	if err != nil {
		println("配置Y坐标值错误, 必须是个整数: " + confArr[3])
		os.Exit(0)
	}
}
