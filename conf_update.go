package easyconf

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func (cf *Conf) AddComment(title string, comment ...string) {
	cf.addItem(nil, "", nil, title, comment...)
}

func (cf *Conf) addItem(pval any, name string, defval any, title string, usage ...string) {
	cf.items = append(cf.items, newConfItem(pval, name, defval, title, usage...))
}

// setItemVar 设置配置项的值。
// 允许设置值为空字符串。
func (cf *Conf) setItemVar(k, v string) error {
	var err error
	if k == "" {
		return fmt.Errorf("配置项的键不能为空")
	}
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

// SetValuesByCmdArgs 从命令行参数获取配置。优先级高
// 允许设置值为空字符串。
// TODO 对于bool类型的参数解析可能会有BUG
func (cf *Conf) SetValuesByCmdArgs() []error {
	for _, item := range cf.items {
		// 注释语句 Name 为空字符
		if item.Name != "" {
			v, ok := cf.kvmp[item.Name]
			if !ok || v == nil {
				cf.kvmp[item.Name] = new(string)
			}
			flag.StringVar(cf.kvmp[item.Name], item.Name, *cf.kvmp[item.Name], strings.Join(item.Usage, ";"))
		}
	}
	flag.Parse()
	var errs []error
	for k, v := range cf.kvmp {
		if err := cf.setItemVar(k, *v); err != nil {
			errs = append(errs, fmt.Errorf("设置项%s配置值%s设置失败:%v", k, *v, err))
		}
	}
	return errs
}

// SetValuesByEnv 从操作系统的环境变量获取配置。优先级中
// 配置值为空字符串会被忽略
func (cf *Conf) SetValuesByEnv() error {
	var err error
	for _, item := range cf.items {
		if item.Name == "" {
			// 注释语句 Name 为空字符
			continue
		}
		v := os.Getenv(item.Name)
		if v != "" {
			cf.kvmp[item.Name] = &v
			err1 := item.setValueVar(v)
			if err1 != nil {
				err = err1
			}
		}
	}
	return err
}

// SetValuesByEnvFile 从env配置文件更新配置项。优先级低。
// 配置值为空字符串会被忽略
func (cf *Conf) SetValuesByEnvFile(envfile string) {
	content, err := os.ReadFile(envfile)
	if err != nil {
		panic(err)
	}
	contstr := string(content)
	lines := strings.Split(contstr, "\n")
	// 解析env文件的每一行
	for _, line := range lines {
		itemk, itemv := GetConfStrByLine(line)
		if itemk == "" || itemv == "" {
			continue
		}
		// fmt.Printf("-----ReadFile(%s)-----k(%s)--v(%s)--------\n", readfile, itemk, itemv)
		cf.kvmp[itemk] = &itemv
		cf.setItemVar(itemk, itemv)
	}
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
