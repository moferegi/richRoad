package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxUnzipFiles      = 10000
	maxUnzipSingleSize = 128 * 1024 * 1024
	maxUnzipTotalSize  = 512 * 1024 * 1024
)

// 解压
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	destAbs, err := filepath.Abs(destDir)
	if err != nil {
		return []string{}, err
	}

	var totalSize int64
	fileCount := 0

	for _, f := range zipReader.File {
		fileCount++
		if fileCount > maxUnzipFiles {
			return []string{}, fmt.Errorf("压缩包文件数量超过限制")
		}
		if f.UncompressedSize64 > maxUnzipSingleSize {
			return []string{}, fmt.Errorf("%s 文件过大", f.Name)
		}
		totalSize += int64(f.UncompressedSize64)
		if totalSize > maxUnzipTotalSize {
			return []string{}, fmt.Errorf("压缩包解压后总大小超过限制")
		}

		cleanName := filepath.Clean(strings.ReplaceAll(f.Name, "\\", "/"))
		if cleanName == "" || cleanName == "." {
			continue
		}
		if strings.HasPrefix(cleanName, "..") || strings.Contains(cleanName, "../") {
			return []string{}, fmt.Errorf("%s 文件名不合法", f.Name)
		}

		fpath := filepath.Join(destDir, cleanName)
		fpathAbs, err := filepath.Abs(fpath)
		if err != nil {
			return []string{}, err
		}
		rel, err := filepath.Rel(destAbs, fpathAbs)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return []string{}, fmt.Errorf("%s 文件路径越界", f.Name)
		}

		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return []string{}, err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				_ = inFile.Close()
				return []string{}, err
			}

			_, err = io.Copy(outFile, inFile)
			closeOutErr := outFile.Close()
			closeInErr := inFile.Close()
			if err != nil {
				return []string{}, err
			}
			if closeOutErr != nil {
				return []string{}, closeOutErr
			}
			if closeInErr != nil {
				return []string{}, closeInErr
			}
		}
	}
	return paths, nil
}
