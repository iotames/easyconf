package easyconf

import (
	"fmt"
	// "os"
	"strconv"
	"strings"
)

// type IConfItem interface {
// 	// GetName() string
// 	// GetTitle() string
// 	GetValue() string
// 	GetValueAny() any
// 	// GetDefaultValue() string
// }

type ConfItem struct {
	Name, Title         string
	Value, DefaultValue any
	Usage               []string
}

//	func (item ConfItem) GetValueAny() any {
//		return item.Value
//	}
func (item *ConfItem) setValueVar(vv string) error {
	var err error
	v := item.Value
	vv = strings.TrimSpace(vv)
	switch val := v.(type) {
	case *int:
		*val, err = strconv.Atoi(vv)
	case *float64:
		*val, err = strconv.ParseFloat(vv, 64)
	case *bool:
		if strings.EqualFold(vv, "true") {
			*val = true
		}
	case *string:
		*val = vv
	case *[]string:
		vsplit := strings.Split(vv, ",")
		var vlist []string
		for _, v1 := range vsplit {
			vlist = append(vlist, strings.TrimSpace(v1))
		}
		*val = vlist
	default:
		err = fmt.Errorf("配置项%s的值为不支持的变量类型(%t)", item.Name, v)
	}
	return err
}

func (item ConfItem) GetValue() string {
	return anyToString(item.Value, item.Name)
}

func (item ConfItem) GetDefaultValue() string {
	return anyToString(item.DefaultValue, item.Name)
}

func (item ConfItem) toString(isDefaultValue bool) string {
	var result []string
	result = append(result, fmt.Sprintf("# %s. The default value is: %s", item.Title, item.GetDefaultValue()))
	if len(item.Usage) > 0 {
		for _, v := range item.Usage {
			result = append(result, fmt.Sprintf("# %s", v))
		}
	}
	if isDefaultValue {
		result = append(result, fmt.Sprintf("%s = %s", item.Name, item.GetDefaultValue()))
	} else {
		result = append(result, fmt.Sprintf("%s = %s", item.Name, item.GetValue()))
	}
	return strings.Join(result, "\n")
}

func (item ConfItem) String() string {
	return item.toString(false)
}

func (item ConfItem) DefaultString() string {
	return item.toString(true)
}

func anyToString(v any, k string) string {
	result := ""
	switch val := v.(type) {
	case nil:
		panic(fmt.Errorf("配置项%s不能为nil", k))
	case int, int64:
		result = fmt.Sprintf("%d", val)
	case *int:
		result = fmt.Sprintf("%d", *val)
	case *float64:
		result = fmt.Sprintf("%.2f", *val)
	case float64:
		result = fmt.Sprintf("%.2f", val)
	case bool:
		result = fmt.Sprintf("%t", val)
	case *bool:
		result = fmt.Sprintf("%t", *val)
	case string:
		result = val
	case *string:
		result = *val
	case *[]string:
		result = strings.Join(*val, ",")
	case []string:
		result = strings.Join(val, ",")
	default:
		panic(fmt.Errorf("配置项%s的值为不支持的变量类型:%t", k, v))
	}
	return result
}

func newConfItem(pval any, name string, defval any, title string, usage ...string) *ConfItem {
	return &ConfItem{
		Value:        pval,
		Name:         name,
		DefaultValue: defval,
		Title:        title,
		Usage:        usage,
	}
}

// func getEnvDefaultStr(key, defval string) string {
// 	v, ok := os.LookupEnv(key)
// 	if !ok {
// 		return defval
// 	}
// 	return v
// }

// func getEnvDefaultBool(key string, defval bool) bool {
// 	v, ok := os.LookupEnv(key)
// 	if !ok {
// 		return defval
// 	}
// 	return strings.EqualFold(v, "true")
// }

// func getEnvDefaultInt(key string, defval int) int {
// 	v, ok := os.LookupEnv(key)
// 	if !ok {
// 		return defval
// 	}
// 	vv, _ := strconv.Atoi(v)
// 	return vv
// }

// // getEnvDefaultStrList 切片的每个元素去掉收尾空格，空字符串对应长度为0的空切片。
// func getEnvDefaultStrList(key string, defval string, sep string) []string {
// 	v, ok := os.LookupEnv(key)
// 	if !ok {
// 		v = defval
// 	}
// 	v = strings.TrimSpace(v)
// 	if v == "" {
// 		return []string{}
// 	}
// 	vv := strings.Split(v, sep)
// 	var result []string
// 	for _, iv := range vv {
// 		vvv := strings.TrimSpace(iv)
// 		if vvv != "" {
// 			result = append(result, vvv)
// 		}
// 	}
// 	return result
// }
