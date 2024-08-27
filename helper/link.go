package helper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type LinkCreator struct {
	src string // 软链源文件
}

func NewLinkCreator(src string) *LinkCreator {
	return &LinkCreator{
		src: src,
	}
}

func (c *LinkCreator) LinkCmderBin(key any) error {
	dst := fmt.Sprintf("%s", key)
	os.Remove(dst)
	return os.Symlink(c.src, dst)
}

// LinkCmderBin 创建软链
func LinkCmderBin(key any) error {
	dirName, _ := os.Getwd()
	creator := NewLinkCreator(filepath.Join(dirName, path.Base(os.Args[0])))
	return creator.LinkCmderBin(key)
}
