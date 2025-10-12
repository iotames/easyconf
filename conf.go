package easyconf

// import "path/filepath"

import (
	"fmt"
	"os"

	"strings"
)

type Conf struct {
	files []string
	kvmp  map[string]*string
	items []*ConfItem
}

// NewConf 定义配置文件。留空默认为: .env, default.env
func NewConf(files ...string) *Conf {
	defaultFiles := []string{".env", "default.env"}
	if len(files) == 0 {
		files = defaultFiles
	}
	for _, ff := range files {
		if ff == "" {
			panic("配置文件路径不能为空")
		}
	}
	return &Conf{files: files}
}

// DefaultString 默认配置的文件内容
func (cf Conf) DefaultString() string {
	var result []string
	for _, item := range cf.items {
		result = append(result, item.DefaultString())
	}
	return strings.Join(result, "\n\n")
}

func (cf Conf) String() string {
	var result []string
	for _, item := range cf.items {
		result = append(result, item.String())
	}
	return strings.Join(result, "\n\n")
}

func createEnvFile(fpath, content string) error {
	if fpath == "" {
		return fmt.Errorf("createEnvFile empty file")
	}
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("create env file(%s)err(%v)", fpath, err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("write env file(%s)err(%v)", fpath, err)
	}
	return f.Close()
}
