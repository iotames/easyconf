package easyconf

import (
	"fmt"
	// "os"
	"strconv"
	"strings"
)

// ConfItem 配置项信息。
// 包含配置名，标题，配置值，默认值，使用说明。
type ConfItem struct {
	Name         string // 配置项键名。键名为空字符串，则该项的值可能为注释。
	Title        string
	ValueStr     string // 配置项的字符串原始值。对于只有字符串类型的环境变量很有用。
	Value        any    // 配置项的值，类型为指针，引用传递。
	DefaultValue any    // 配置项的默认值，值传递。
	Usage        []string
}

func parseIntList(val *[]int, vv string, defaultVal []int) error {
	var err error
	vsplit := strings.Split(vv, ",")
	var vlist []int
	var vint int
	for _, v1 := range vsplit {
		vint, err = strconv.Atoi(strings.TrimSpace(v1))
		if err != nil {
			break
		}
		vlist = append(vlist, vint)
	}
	*val = vlist
	if err != nil {
		*val = defaultVal
	}
	return err
}

func parseStringList(val *[]string, vv string) {
	vsplit := strings.Split(vv, ",")
	var vlist []string
	for _, v1 := range vsplit {
		vlist = append(vlist, strings.TrimSpace(v1))
	}
	*val = vlist
}

func (item *ConfItem) setValueVar(vv string) error {
	var err error
	v := item.Value
	vv = strings.TrimSpace(vv)
	switch val := v.(type) {
	case *int:
		*val, err = strconv.Atoi(vv)
		if err != nil {
			*val = item.DefaultValue.(int)
		}
	case *[]int:
		err = parseIntList(val, vv, item.DefaultValue.([]int))
	case *float64:
		*val, err = strconv.ParseFloat(vv, 64)
		if err != nil {
			*val = item.DefaultValue.(float64)
		}
	case *bool:
		if strings.EqualFold(vv, "true") {
			*val = true
		} else {
			*val = false
		}
	case *string:
		*val = vv
	case *[]string:
		parseStringList(val, vv)
	default:
		err = fmt.Errorf("配置项%s的值为不支持的变量类型(%t)", item.Name, v)
	}
	return err
}

func (item ConfItem) GetValue() string {
	switch val := item.Value.(type) {
	case nil:
		panic(fmt.Errorf("配置项%s的指针不能为nil", item.Name))
	case *int:
		return anyToString(*val, item.Name)
	case *float64:
		return anyToString(*val, item.Name)
	case *bool:
		return anyToString(*val, item.Name)
	case *string:
		return anyToString(*val, item.Name)
	case *[]string:
		return anyToString(*val, item.Name)
	case *[]int:
		return anyToString(*val, item.Name)
	default:
		panic(fmt.Errorf("配置项%s的值为不支持的变量类型:%t", item.Name, item.Value))
	}
}

func (item ConfItem) GetDefaultValue() string {
	return anyToString(item.DefaultValue, item.Name)
}

func (item ConfItem) toString(isDefaultValue bool) string {
	var result []string
	var defval, realval string
	if item.Title != "" {
		// ADD COMMENT: Title
		result = append(result, fmt.Sprintf("# %s", item.Title))
	}
	if item.DefaultValue != nil && item.Value != nil {
		switch item.DefaultValue.(type) {
		case string:
			defval = fmt.Sprintf(`"%s"`, item.GetDefaultValue())
			realval = fmt.Sprintf(`"%s"`, item.GetValue())
		default:
			defval = item.GetDefaultValue()
			realval = item.GetValue()
		}
		// ADD COMMENT: default value
		result = append(result, fmt.Sprintf(`# The default value is: %s`, defval))
	}
	if len(item.Usage) > 0 {
		for _, v := range item.Usage {
			// ADD COMMENT: usage
			result = append(result, fmt.Sprintf("# %s", v))
		}
	}
	if item.Value != nil {
		var val string
		if isDefaultValue {
			val = defval
		} else {
			val = realval
		}
		// ADD name and value pair
		result = append(result, fmt.Sprintf(`%s = %s`, item.Name, val))
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
		panic(fmt.Errorf("配置项%s的值不能为nil", k))
	case int:
		result = fmt.Sprintf("%d", val)
	case float64:
		result = fmt.Sprintf("%.6f", val)
	case bool:
		result = fmt.Sprintf("%t", val)
	case string:
		result = val
	case []string:
		result = strings.Join(val, ",")
	case []int:
		var vvv []string
		for _, v1 := range val {
			vvv = append(vvv, fmt.Sprintf("%d", v1))
		}
		result = strings.Join(vvv, ",")
	default:
		panic(fmt.Errorf("配置项%s的值为不支持的变量类型:%T", k, v))
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
