package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// TarGZ, 归档操作函数
//
// 参数:
//
//	dest: 归档文件路径
//	filenames: 配置文件路径切片
//
// 返回值:
//
//	错误
func TarGZ(dest string, filenames []string) (err error) {
	if filenames == nil || len(filenames) < 1 {
		return errors.New("filename list is null")
	}

	// 归档目录
	destDir := filepath.Dir(dest)

	// 归档文件初始化
	d, _ := os.Create(dest) // 创建归档文件
	defer d.Close()
	gw := gzip.NewWriter(d) // 转换归档文件为gzip压缩文件
	defer gw.Close()
	tgzw := tar.NewWriter(gw) // 转换归档文件为tar.gz/tgz归档文件
	defer tgzw.Close()

	// 开始归档
	for _, f := range filenames {

		// 初始化被归档文件信息
		relDir, relErr := filepath.Rel(destDir, filepath.Dir(f))
		if relErr != nil {
			return relErr
		}
		fd, opErr := os.Open(f)
		if opErr != nil {
			return opErr
		}

		// 开始归档
		compErr := compress(fd, relDir, tgzw)
		fd.Close()
		if compErr != nil {
			return compErr
		}

	}
	return
}

// compress, 归档压缩子函数
//
// 参数:
//
//	fd: 被归档文件的系统文件对象指针
//	prefix: 被归档文件的目录路径
//	tgzw: tar文件对象指针
//
// 返回值:
//
//	错误
func compress(fd *os.File, prefix string, tgzw *tar.Writer) error {

	// 加载被归档文件信息
	info, infoErr := fd.Stat()
	if infoErr != nil {
		return infoErr
	}

	// 文件归档初始化
	header, hErr := tar.FileInfoHeader(info, "")
	if hErr != nil {
		return hErr
	}
	// 调整被归档文件目录信息
	header.Name = filepath.Join(prefix, header.Name)
	// 完成该文件归档
	twErr := tgzw.WriteHeader(header)
	if twErr != nil {
		return twErr
	}
	_, ioErr := io.Copy(tgzw, fd)
	if ioErr != nil {
		return ioErr
	}

	return nil
}
