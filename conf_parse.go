package easyconf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iotames/miniutils"
)

// Parse 从配置文件和系统环境变量解析配置
func (cf *Conf) Parse() error {
	var err error
	cf.kvmp = make(map[string]*string, len(cf.items))
	for _, file := range cf.files {
		if !miniutils.IsPathExists(file) {
			err = createEnvFile(file, cf.DefaultString())
			if err != nil {
				panic(err)
			}
			// fmt.Printf("Create file %s SUCCESS\n", file)
		}
	}

	filenum := len(cf.files)
	lasti := filenum - 1
	for i := 0; i < filenum; i++ {
		readfile := cf.files[lasti-i]
		// 按顺序，后面的文件配置，覆盖前面的文件配置
		cf.SetValuesByEnvFile(readfile)
	}

	// 从系统环境变量获取配置，配置文件上的同名配置会被覆盖
	err = cf.SetValuesByEnv()
	if err != nil {
		return err
	}

	// 命令行参数优先级最高。
	errs := cf.SetValuesByCmdArgs()
	var errstr []string
	for _, err := range errs {
		if err != nil {
			errstr = append(errstr, err.Error())
		}
	}
	if len(errstr) > 0 {
		return fmt.Errorf("命令行参数解析错误:%s", strings.Join(errstr, ";"))
	}
	return nil
}

var seps []string = []string{`"`, `'`}

func GetConfStrByLine(line string) (itemk, itemv string) {
	remarkk := "#"
	v := strings.TrimSpace(line)
	if strings.Index(v, remarkk) == 0 {
		// 忽略以注释符 # 开头的一整行
		return
	}
	eqIndex := strings.Index(v, "=")
	if eqIndex > 0 {
		// 等号 = 左边为配置名
		itemk = strings.TrimSpace(v[:eqIndex])
		// 等号 = 右边为配置值
		itemv = strings.TrimSpace(v[eqIndex+1:])

		// 忽略含注释号 # 的配置名
		if strings.Contains(itemk, remarkk) {
			itemk = ""
			return
		}

		// 配置值使用双引号 " 或单引号 ' 包裹，则提取出来。
		for _, sep := range seps {
			if strings.Index(itemv, sep) == 0 && itemv[len(itemv)-1] == sep[0] {
				itemv = itemv[1 : len(itemv)-1]

				// 解决不适配 name="value" # remark="a=b"
				if strings.Contains(itemv, remarkk) {
					re := regexp.MustCompile(fmt.Sprintf(`(\%s( |\t|)*%s)`, sep, remarkk))
					substr := re.FindString(itemv)
					if substr != "" {
						itemv = strings.Split(itemv, substr)[0]
					}
				}

				return
			}
		}

		remarkIndex := strings.Index(itemv, remarkk)

		// 忽略以注释符 # 开头的配置值
		if remarkIndex == 0 {
			itemv = ""
			return
		}

		if remarkIndex > 0 {
			// 配置值使用双引号 " 或单引号 ' 包裹，则提取出来。
			for _, sep := range seps {
				lastSepIndex := strings.LastIndex(itemv, sep)

				// 适配 name="hello#word", name="value" # remark
				if strings.Index(itemv, sep) == 0 && lastSepIndex > 0 {
					itemv = itemv[1:lastSepIndex]
					// 解决不适配 name="value" # remark="a=b"
					if strings.Contains(itemv, remarkk) {
						re := regexp.MustCompile(fmt.Sprintf(`(\%s( |\t)*%s)`, sep, remarkk))
						substr := re.FindString(itemv)
						if substr != "" {
							itemv = strings.Split(itemv, substr)[0]
						}
					}
					return
				}
			}
			itemv = strings.TrimSpace(itemv[:remarkIndex])
		}
	}

	return
}
