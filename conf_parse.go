package easyconf

import (
	"fmt"
	"os"

	"strings"

	"github.com/iotames/miniutils"
)

func (cf *Conf) Parse() error {
	var err error
	var content []byte

	for _, file := range cf.files {
		if !miniutils.IsPathExists(file) {
			err = createEnvFile(file, cf.DefaultString())
			if err != nil {
				panic(err)
			}
			fmt.Printf("Create file %s SUCCESS\n", file)
		}
	}

	filenum := len(cf.files)
	lasti := filenum - 1
	for i := 0; i < filenum; i++ {
		readfile := cf.files[lasti-i]
		fmt.Println(readfile)
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

	// for _, arg := range cf.items {
	// 	fmt.Printf("-----cf.item--k(%s)---v(%s)--default(%s)--\n", arg.Name, arg.GetValue(), arg.GetDefaultValue())
	// }

	return nil
}

func GetConfStrByLine(line string) (itemk, itemv string) {
	v := strings.TrimSpace(line)
	if strings.Index(v, "#") == 0 {
		return
	}
	if strings.Contains(v, "=") {
		itemsplit := strings.Split(v, "=")
		itemk = strings.TrimSpace(itemsplit[0])
		if strings.Contains(itemk, "#") {
			itemk = ""
			return
		}
		itemv = strings.TrimSpace(itemsplit[1])
		if strings.Index(itemv, `"`) == 0 && itemv[len(itemv)-1] == '"' {
			itemv = itemv[1 : len(itemv)-1]
		}
		if strings.Index(itemv, `'`) == 0 && itemv[len(itemv)-1] == '\'' {
			itemv = itemv[1 : len(itemv)-1]
		}
	}
	return
}
