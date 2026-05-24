package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 前端传来文件片与当前片为什么文件的第几片
// 后端拿到以后比较次分片是否上传 或者是否为不完全片
// 前端发送每片多大
// 前端告知是否为最后一片且是否完成

const (
	breakpointDir = "./breakpointDir/"
	finishDir     = "./fileDir/"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: BreakPointContinue
//@description: 断点续传
//@param: content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string
//@return: error, string

func BreakPointContinue(content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string) (string, error) {
	if !isSafePathSegment(fileName) || !isSafePathSegment(fileMd5) {
		return "", errors.New("文件名或路径不合法")
	}
	path := filepath.Join(breakpointDir, fileMd5)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return path, err
	}
	pathC, err := makeFileContent(content, fileName, path, contentNumber)
	return pathC, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CheckMd5
//@description: 检查Md5
//@param: content []byte, chunkMd5 string
//@return: CanUpload bool

func CheckMd5(content []byte, chunkMd5 string) (CanUpload bool) {
	fileMd5 := MD5V(content)
	if fileMd5 == chunkMd5 {
		return true // 可以继续上传
	} else {
		return false // 切片不完整，废弃
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: makeFileContent
//@description: 创建切片内容
//@param: content []byte, fileName string, FileDir string, contentNumber int
//@return: string, error

func makeFileContent(content []byte, fileName string, FileDir string, contentNumber int) (string, error) {
	if !isSafePathSegment(fileName) {
		return "", errors.New("文件名或路径不合法")
	}
	path := filepath.Join(FileDir, fileName+"_"+strconv.Itoa(contentNumber))
	f, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return path, err
	}

	return path, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: makeFileContent
//@description: 创建切片文件
//@param: fileName string, FileMd5 string
//@return: error, string

func MakeFile(fileName string, FileMd5 string) (string, error) {
	if !isSafePathSegment(fileName) || !isSafePathSegment(FileMd5) {
		return "", errors.New("文件名或路径不合法")
	}
	chunkDir := filepath.Join(breakpointDir, FileMd5)
	rd, err := os.ReadDir(chunkDir)
	if err != nil {
		return filepath.Join(finishDir, fileName), err
	}
	_ = os.MkdirAll(finishDir, os.ModePerm)
	finishPath := filepath.Join(finishDir, fileName)
	fd, err := os.OpenFile(finishPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return finishPath, err
	}
	defer fd.Close()
	for k := range rd {
		content, _ := os.ReadFile(filepath.Join(chunkDir, fileName+"_"+strconv.Itoa(k)))
		_, err = fd.Write(content)
		if err != nil {
			_ = os.Remove(finishPath)
			return finishPath, err
		}
	}
	return finishPath, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: RemoveChunk
//@description: 移除切片
//@param: FileMd5 string
//@return: error

func RemoveChunk(FileMd5 string) error {
	if !isSafePathSegment(FileMd5) {
		return errors.New("路径不合法")
	}
	err := os.RemoveAll(filepath.Join(breakpointDir, FileMd5))
	return err
}

func isSafePathSegment(name string) bool {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return false
	}
	if strings.Contains(trimmed, "..") {
		return false
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return false
	}
	return trimmed == filepath.Base(trimmed)
}
