package tools

import (
	"archive/zip"
	"compress/flate"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

func CompressFilesToZip(zipAbsPath string, fileAbsPaths ...string) error {
	basePath := filepath.Dir(zipAbsPath)
	_, err := os.Stat(basePath)
	if err != nil {
		return errors.WithStack(err)
	}

	// 创建 zip 文件
	zipFile, err := os.Create(zipAbsPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer zipFile.Close()

	// 创建一个新的 zip 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 设置压缩级别为 Deflate
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	for _, fileAbsPath := range fileAbsPaths {
		err = compressFileToZip(fileAbsPath, zipWriter)
		if err != nil {
			return errors.WithMessagef(err, "compress file(%v) failed", fileAbsPath)
		}
	}
	return nil
}

func compressFileToZip(fileAbsPath string, zipWriter *zip.Writer) error {
	// 打开源文件
	file, err := os.Open(fileAbsPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	// 创建新的 zip 文件头
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return errors.WithStack(err)
	}
	header.Name = filepath.Base(fileAbsPath)
	header.Method = zip.Deflate // 设置压缩方法

	// 将文件头写入 zip 文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.WithStack(err)
	}

	// 复制文件内容到 zip 文件
	_, err = io.Copy(writer, file)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
