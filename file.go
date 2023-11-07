package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
func ListDir(dirPth string, suffixList []string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		println("read dir: " + dirPth + " fail: " + err.Error())
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		for _, suffix := range suffixList {
			if strings.HasSuffix(strings.ToLower(fi.Name()), strings.ToLower(suffix)) {
				files = append(files, fi.Name())
			}
		}
	}
	return files, nil
}

// 检查两个源目录，创建目标目录
func getListAndMkTargetDir() {
	var err error
	// 文件夹检测
	srcList, err = ListDir(srcDir, fileSuffix)
	if err != nil {
		println("读取源目录失败: " + err.Error())
		os.Exit(0)
	}
	if len(srcList) <= 0 {
		println("确定这个模板图目录下有.png文件吗: " + srcDir)
		os.Exit(0)
	}

	coverList, err = ListDir(coverDir, fileSuffix)
	if err != nil {
		println("读取源目录失败: " + err.Error())
		os.Exit(0)
	}
	if len(coverList) <= 0 {
		println("确定这个二维码目录下有.png文件吗: " + srcDir)
		os.Exit(0)
	}
	err = os.Mkdir(resultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		println("创建结果目录失败: " + err.Error())
		os.Exit(0)
	}
}
