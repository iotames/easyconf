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
func (c JsonConf) Save(v any, filename string) error {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	fpath := filepath.Join(c.dirPath, filename)
	// 临时文件须与目标同目录，保证 rename 在同一文件系统上原子生效
	tmp, err := os.CreateTemp(filepath.Dir(fpath), filepath.Base(fpath)+".tmp-*")
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name()) // rename 成功后临时文件已不存在，此处报错无害
	}()
	// 权限位与旧实现（OpenFile 0777）保持一致，避免行为漂移
	if err := tmp.Chmod(0o777); err != nil {
		return err
	}
	if _, err := tmp.Write(b); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), fpath)
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
