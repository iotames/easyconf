package easyconf

import (
	"encoding/json"
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
// 原子写：先写同目录临时文件，再 rename 替换——写入中途崩溃不留半截文件，
// 并根治直接覆写时新内容变短、文件尾部残留旧字节导致 JSON 损坏的问题。
// 权限与 UpdateFile 一致：已有文件沿用原权限，新建文件 0644。
func (c JsonConf) Save(v any, filename string) error {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	return writeFileAtomic(filepath.Join(c.dirPath, filename), b, 0o644)
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
