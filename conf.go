package easyconf

// import "path/filepath"

import (
	"fmt"
	"os"

	"strings"
)

type Conf struct {
	files []string
	items []*ConfItem
}

// NewConf 定义配置文件。留空默认为: .env, default.env
func NewConf(files ...string) *Conf {
	defaultFiles := []string{".env", "default1.env", "default2.env"}
	if len(files) == 0 {
		files = defaultFiles
	}
	return &Conf{files: files}
}

func (cf *Conf) StringVar(pval *string, name string, defval, title string, usage ...string) {
	item := newConfItem(pval, name, defval, title, usage...)
	cf.items = append(cf.items, item)
}

func (cf *Conf) BoolVar(pval *bool, name string, defval bool, title string, usage ...string) {
	item := newConfItem(pval, name, defval, title, usage...)
	cf.items = append(cf.items, item)
}

func (cf *Conf) IntVar(pval *int, name string, defval int, title string, usage ...string) {
	item := newConfItem(pval, name, defval, title, usage...)
	cf.items = append(cf.items, item)
}

func (cf *Conf) StringListVar(pval *[]string, name string, defval []string, title string, usage ...string) {
	item := newConfItem(pval, name, defval, title, usage...)
	cf.items = append(cf.items, item)
}

func (cf *Conf) setItemVar(k, v string) error {
	var err error
	for _, arg := range cf.items {
		if arg.Name == k {
			err1 := arg.setValueVar(v)
			if err1 != nil {
				err = err1
			}
		}
	}
	return err
}

func (cf Conf) DefaultString() string {
	var result []string
	for _, item := range cf.items {
		result = append(result, item.DefaultString())
	}
	return strings.Join(result, "\n\n")
}

func createEnvFile(fpath, content string) error {
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
