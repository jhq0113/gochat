package utils

import (
	"os"
	"syscall"
)

// FileExists 判断文件是否存在
func FileExists(fileName string) (exists bool, err error) {
	_, err = os.Stat(fileName)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
	}

	return
}

// MkDir 创建目录
func MkDir(dir string, perm os.FileMode) (err error) {
	exists, err := FileExists(dir)
	if err != nil {
		return
	}

	if exists {
		return
	}

	return os.MkdirAll(dir, perm)
}

// FileTime 文件时间
func FileTime(fileName string) (createTime, lastAccessTime, lastWriteTime int64, err error) {
	info, err := os.Stat(fileName)
	if err != nil {
		return
	}

	attr := info.Sys().(*syscall.Stat_t)
	createTime = attr.Ctimespec.Sec
	lastAccessTime = attr.Atimespec.Sec
	lastWriteTime = attr.Mtimespec.Sec
	return
}
