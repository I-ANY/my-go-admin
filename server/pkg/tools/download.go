package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CompressFilesToZipWithTmpDir(fileAbsPath, tmpDir string) (string, error) {
	_, err := os.Stat(tmpDir)
	if err != nil {
		tmpDir = os.TempDir()
	}
	filename := ""
	ss := strings.Split(filepath.Base(fileAbsPath), ".")
	if len(ss) > 1 {
		filename = fmt.Sprintf("%s.zip", strings.Join(ss[:len(ss)-1], "."))
	} else {
		filename = fmt.Sprintf("%s.zip", fileAbsPath)
	}
	zipAbsPath := filepath.Join(tmpDir, filename)
	err = CompressFilesToZip(zipAbsPath, fileAbsPath)
	if err != nil {
		return zipAbsPath, err
	}
	return zipAbsPath, nil
}

func GinDownloadFile(ctx *gin.Context, fileAbsPath, filename string) error {
	if len(filename) == 0 {
		filename = filepath.Base(fileAbsPath)
	}
	stat, err := os.Stat(fileAbsPath)
	if err != nil {
		return errors.WithStack(err)
	}
	ctx.Header("Content-Type", "application/octet-stream")
	// 使用 UTF-8 编码处理中文文件名
	encodedFilename := url.QueryEscape(filename)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedFilename))
	ctx.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))
	file, err := os.Open(fileAbsPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()
	bs := make([]byte, 1024*1024*5)
	// 使用 io.Copy 流式传输文件
	_, err = io.CopyBuffer(ctx.Writer, file, bs)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GinDownFileAndCompress(ctx *gin.Context, fileAbsPath, tmpDir, downloadZipName string) error {
	zipAbsPath, err := CompressFilesToZipWithTmpDir(fileAbsPath, tmpDir)
	if err != nil {
		err = errors.WithMessage(err, "compress files to zip failed")
		return err
	}
	defer func() {
		if len(zipAbsPath) > 0 {
			os.Remove(zipAbsPath)
		}
	}()
	err = GinDownloadFile(ctx, zipAbsPath, downloadZipName)
	if err != nil {
		err = errors.WithMessage(err, "download file failed")
		return err
	}
	return nil
}
