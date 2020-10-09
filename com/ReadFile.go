package com

import (
	"fmt"
	"log"
	"os"
	"redisCluster/asset"
	"strings"
)

type Conf struct {
	configFile map[string][]byte
}

var tempDir = getTempDir()

func (c *Conf) readFile() {

	c.configFile = make(map[string][]byte)
	names := asset.AssetNames()

	for _, filePath := range names {
		fileByte, err := asset.Asset(filePath)
		if err != nil {
			log.Println(err, filePath)
		} else {
			c.configFile[filePath] = fileByte
		}
	}
}

// 将dir解压到当前目录：根据生成的.go文件，将其解压为当前文件
func (c *Conf) Restore() {
	c.readFile()
	var resourceFile *os.File
	defer resourceFile.Close()
	for key, value := range c.configFile {

		index := strings.LastIndex(key, "/")

		createSuffixDir(key[:index])

		if !isExist(tempDir + key) {
			resourceFile, err2 := os.Create(tempDir + key)

			if err2 != nil {
				log.Fatal("创建资源文件  ", err2)
				os.Exit(1)
			}
			resourceFile.Write(value)
			if strings.HasSuffix(key, ".sh") || strings.Contains(key, "/bin") {
				resourceFile.Chmod(os.ModePerm)
			}
			resourceFile.Close()
		}
	}
}

/*获取临时生成的文件夹*/
func getTempDir() string {
	var tempDir strings.Builder
	tempDir.WriteString(os.TempDir())
	tempDir.WriteString("/cluster/")

	return tempDir.String()
}

/*生成临时文件夹*/
func createSuffixDir(suffix string) {
	var suffixDir strings.Builder
	suffixDir.WriteString(tempDir)
	suffixDir.WriteString(suffix)
	if !isExist(suffixDir.String()) {
		err := os.MkdirAll(suffixDir.String(), os.ModePerm)
		if err != nil {
			log.Fatal("存储文件,创建临时文件夹出现错误")
		}
	}
}

//判断文件或文件夹是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
