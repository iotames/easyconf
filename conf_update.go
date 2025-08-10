package easyconf

import (
	"fmt"
	"os"
)

func (cf *Conf) AddComment(title string, comment ...string) {
	cf.addItem(nil, "", nil, title, comment...)
}

func (cf *Conf) addItem(pval any, name string, defval any, title string, usage ...string) {
	cf.items = append(cf.items, newConfItem(pval, name, defval, title, usage...))
}

func (cf *Conf) setItemVar(k, v string) error {
	var err error
	for _, arg := range cf.items {
		if arg.Name == "" {
			// 注释语句 Name 为空字符
			continue
		}
		if arg.Name == k {
			err1 := arg.setValueVar(v)
			if err1 != nil {
				err = err1
			}
		}
	}
	return err
}

// SetItemVarByEnv 从操作系统的环境变量获取配置
func (cf *Conf) SetItemVarByEnv() error {
	var err error
	for _, item := range cf.items {
		if item.Name == "" {
			// 注释语句 Name 为空字符
			continue
		}
		v := os.Getenv(item.Name)
		if v != "" {
			err1 := item.setValueVar(v)
			if err1 != nil {
				err = err1
			}
		}
	}
	return err
}

func (cf *Conf) UpdateFile(fpath string) error {
	var f *os.File
	var err error
	if fpath == "" {
		fpath = cf.files[0]
	}
	f, err = os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("open file(%s)err(%v)", fpath, err)
	}
	_, err = f.WriteString(cf.String())
	if err != nil {
		return fmt.Errorf("write file(%s)err(%v)", fpath, err)
	}
	return f.Close()
}

// DisableEnvFile 彻底禁用环境变量配置文件。只从系统环境变量获取。
func (cf *Conf) DisableEnvFile() {
	cf.disableEnvFile = true
}
