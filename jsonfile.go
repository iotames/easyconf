package easyconf

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// JsonConf 以 JSON 文件格式保存和读取配置
type JsonConf struct {
	dirPath string
}

// NewJsonConf 创建 JsonConf 实例，dirPath 为配置文件存放目录
func NewJsonConf(dirPath string) *JsonConf {
	return &JsonConf{dirPath: dirPath}
}

// Save 将配置 v 保存到 dirPath/filename 文件中
// filename 不包含目录路径。应该包含 .json 后缀，如 config.json。
func (c JsonConf) Save(v any, filename string) error {
	var err error
	var b []byte
	b, err = json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	var f io.WriteCloser
	fpath := filepath.Join(c.dirPath, filename)
	f, err = os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

// Read 从 dirPath/filename 文件中读取配置到 v
// filename 不包含目录路径。应该包含 .json 后缀，如 config.json。
func (c JsonConf) Read(v any, filename string) error {
	var b []byte
	var err error
	b, err = os.ReadFile(filepath.Join(c.dirPath, filename))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
