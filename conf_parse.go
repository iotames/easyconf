package easyconf

import (
	"fmt"
	"os"

	"regexp"
	"strings"

	"github.com/iotames/miniutils"
)

// Parse 从配置文件和系统环境变量解析配置
func (cf *Conf) Parse() error {
	var err error
	var content []byte

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
		content, err = os.ReadFile(readfile)
		if err != nil {
			panic(err)
		}
		contstr := string(content)
		lines := strings.Split(contstr, "\n")
		for _, line := range lines {
			itemk, itemv := GetConfStrByLine(line)
			if itemk == "" || itemv == "" {
				continue
			}
			// fmt.Printf("-----ReadFile(%s)-----k(%s)--v(%s)--------\n", readfile, itemk, itemv)
			cf.setItemVar(itemk, itemv)
		}
	}

	// 从系统环境变量获取配置，配置文件上的同名配置会被覆盖
	return cf.SetItemVarByEnv()
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
