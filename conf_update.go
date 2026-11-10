package easyconf

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// AddComment 添加纯注释行到配置输出中。
// title 为注释标题（生成为 # title），comment 为多行注释内容。
// 这些注释行无关联的配置键名，在解析时会被跳过。
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
	for _, item := range cf.items {
		if item.Name == k {
			err1 := item.setValueVar(v)
			if err1 != nil {
				err = err1
			} else {
				item.ValueStr = v
			}
		}
	}
	return err
}

// SetItemValue 设置指定配置项的值，并更新指针指向的变量。键名 k 不可为空，值 v 允许为空字符串。
func (cf *Conf) SetItemValue(k, v string) error {
	return cf.setItemVar(k, v)
}

// SetValuesByCmdArgs 从命令行参数获取配置。优先级高
// 允许设置值为空字符串。
func (cf *Conf) SetValuesByCmdArgs() []error {
	var errs []error
	for _, item := range cf.items {
		// 注释语句 Name 为空字符
		if item.Name != "" {
			vstr := item.ValueStr
			v := item.Value
			switch val := v.(type) {
			case *int:
				flag.IntVar(val, item.Name, *val, strings.Join(item.Usage, ";"))
			case *float64:
				flag.Float64Var(val, item.Name, *val, strings.Join(item.Usage, ";"))
			case *bool:
				flag.BoolVar(val, item.Name, *val, strings.Join(item.Usage, ";"))
			case *string:
				flag.StringVar(val, item.Name, *val, strings.Join(item.Usage, ";"))
			case *[]string:
				flag.StringVar(&item.ValueStr, item.Name, vstr, strings.Join(item.Usage, ";"))
			case *[]int:
				flag.StringVar(&item.ValueStr, item.Name, vstr, strings.Join(item.Usage, ";"))
			default:
				errs = append(errs, fmt.Errorf("设置项%s配置值%s不支持变量类型(%T)", item.Name, vstr, v))
			}
		}
	}
	flag.Parse()
	for _, item := range cf.items {
		// 注释语句 Name 为空字符
		if item.Name != "" {
			v := item.Value
			switch val := v.(type) {
			case *[]string:
				parseStringList(val, item.ValueStr)
			case *[]int:
				parseIntList(val, item.ValueStr, *val)
			}
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
			err1 := item.setValueVar(v)
			if err1 != nil {
				err = err1
			} else {
				item.ValueStr = v
			}
		}
	}
	return err
}

// SetValuesByEnvFile 从env配置文件更新配置项。优先级低。
// 配置值为空字符串会被忽略
func (cf *Conf) SetValuesByEnvFile(envfile string) error {
	content, err := os.ReadFile(envfile)
	if err != nil {
		return err
	}
	contstr := string(content)
	lines := strings.Split(contstr, "\n")
	// 解析env文件的每一行
	errStrList := []string{}
	for _, line := range lines {
		itemk, itemv := GetConfStrByLine(line)
		if itemk == "" || itemv == "" {
			continue
		}
		// fmt.Printf("-----ReadFile(%s)-----k(%s)=v(%s)--------\n", envfile, itemk, itemv)
		err = cf.setItemVar(itemk, itemv)
		if err != nil {
			errStrList = append(errStrList, err.Error())
		}
	}
	if len(errStrList) > 0 {
		return fmt.Errorf("env文件(%s)解析错误(SetValuesByEnvFile error):%s", envfile, strings.Join(errStrList, ";"))
	}
	return nil
}

// UpdateFile 将当前配置值持久化到指定文件。
//
// fpath 留空时默认更新 cf.files[0]（第一个配置文件）。
// 采用增量更新策略：只替换已有配置项的值行，保留原文件中的注释、空行和引号风格；
// 重建的行统一为「key = value」格式，原行的缩进与等号两侧间距会被归一化（解析时两侧本就会去空白）。
// 原文件中不存在的配置项会被追加到文件末尾。
// 若文件不存在，会自动创建并用当前值全量写入（权限 0644；已有文件沿用原权限）。
//
// 写入是原子的（临时文件 + rename），不会留下半截文件。
// 落盘前会用本库解析器把每一行读回自检：取值（如含 # 或引号）若无法原样读回，
// 宁可返回错误也不写坏文件——避免「保存成功、重启后取值变了」。
func (cf *Conf) UpdateFile(fpath string) error {
	if fpath == "" {
		fpath = cf.files[0]
	}

	content, err := os.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在时，用当前值全量创建
			content = []byte(cf.String())
		} else {
			return fmt.Errorf("read file(%s) err(%v)", fpath, err)
		}
	}

	// 统一处理 CRLF 为 LF
	text := strings.ReplaceAll(string(content), "\r", "")
	lines := strings.Split(text, "\n")
	updatedKeys := make(map[string]bool)

	for i, line := range lines {
		key, _ := GetConfStrByLine(line)
		if key == "" {
			continue
		}
		for _, item := range cf.items {
			if item.Name == key {
				newLine, err := rebuildValueLine(line, key, item.GetValue())
				if err != nil {
					return err
				}
				lines[i] = newLine
				updatedKeys[key] = true
				break
			}
		}
	}

	// 追加原文件中不存在的新 key
	hasNew := false
	for _, item := range cf.items {
		if item.Name != "" && !updatedKeys[item.Name] {
			if !hasNew {
				// 第一个新 key 前：除非文件末尾已是空行，否则加空行分隔
				if len(lines) > 0 && lines[len(lines)-1] != "" {
					lines = append(lines, "")
				}
			} else {
				// 后续新 key 之间始终加空行分隔
				lines = append(lines, "")
			}
			for _, s := range strings.Split(item.String(), "\n") {
				lines = append(lines, s)
			}
			hasNew = true
		}
	}

	// 落盘前自检：每个配置项写进文件的文本必须能被本库解析器原样读回。
	// 写方与读方必须自洽——否则会出现「保存成功、重启后取值变了」这种最难排查的问题。
	// 自检不通过就报错，不把读不回来的值写进文件。
	written := make(map[string]string)
	for _, line := range lines {
		k, v := GetConfStrByLine(line)
		if k == "" {
			continue
		}
		written[k] = v
	}
	for _, item := range cf.items {
		if item.Name == "" {
			continue
		}
		got, ok := written[item.Name]
		if !ok {
			return fmt.Errorf("配置项 %s 未能写入文件 %s，请检查该键所在行是否被注释或格式破坏", item.Name, fpath)
		}
		if got != item.GetValue() {
			// 不回显取值：配置项可能含口令或 DSN，错误信息会进入接口响应与日志。
			return fmt.Errorf("配置项 %s 写入文件 %s 后读回的值与内存中的值不一致，请去掉取值中的引号或 # 后重试",
				item.Name, fpath)
		}
	}

	return writeFileAtomic(fpath, []byte(strings.Join(lines, "\n")), 0o644)
}

// UpdateByMap 通过 map 批量更新配置项的值，并将结果持久化到文件。
//
// mp 的键为配置项名称，值为新值。空 map 或 nil map 不会产生任何变更。
// 如果某个键名在注册的配置项中不存在，setItemVar 会静默跳过该键。
// 遇到错误时立即返回，已更新成功的值不会回滚。
func (cf *Conf) UpdateByMap(mp map[string]string, fpath string) error {
	if mp == nil {
		return nil
	}
	var err error
	for k, v := range mp {
		if err = cf.setItemVar(k, v); err != nil {
			return err
		}
	}
	return cf.UpdateFile(fpath)
}

// rebuildValueLine 替换一行中的值为新值，同时保留：
// - 字符串的引号风格（有引号保留引号，无引号不加）
// - 等号后的行内注释（# comment）
//
// 注意：重建后的行统一为「key = value」格式，原行的缩进与等号两侧间距会被归一化
// （本库解析器两侧均做 TrimSpace，不影响读回结果）。
//
// 写进文件的文本必须能被本库的解析器原样读回，否则会出现「运行期生效、重启后被改写」：
// 含 # 的裸值会被当成行内注释截断，含首尾空格的值会被裁掉。
// 所以取值先按原风格生成，再把拼好的整行交给解析器试算；确实无法表达的取值直接报错，不写坏文件。
func rebuildValueLine(line, key, newValue string) (string, error) {
	suffix, quote := splitValueLine(line)
	literal, err := valueLiteral(key, newValue, quote, suffix)
	if err != nil {
		return "", fmt.Errorf("配置项 %s %w", key, err)
	}
	return assembleValueLine(key, literal, suffix), nil
}

// assembleValueLine 按落盘格式拼出一行：「key = literal」加可选的行内注释。
func assembleValueLine(key, literal, suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s = %s %s", key, literal, suffix)
	}
	return fmt.Sprintf("%s = %s", key, literal)
}

// splitValueLine 拆出原行的行内注释与引号风格，供重建时保留。
// 只认以 # 开头的尾随内容为注释；其余无法识别的尾随内容一并丢弃，避免拼出解析不了的行。
func splitValueLine(line string) (suffix, quote string) {
	eqIndex := strings.Index(line, "=")
	if eqIndex < 0 {
		return "", ""
	}
	trimmed := strings.TrimSpace(line[eqIndex+1:])
	if trimmed == "" {
		return "", ""
	}
	if trimmed[0] == '"' || trimmed[0] == '\'' {
		q := string(trimmed[0])
		last := strings.LastIndex(trimmed, q)
		if last <= 0 {
			return "", q
		}
		suffix = strings.TrimSpace(trimmed[last+1:])
		if !strings.HasPrefix(suffix, "#") {
			suffix = ""
		}
		return suffix, q
	}
	if idx := strings.Index(trimmed, "#"); idx >= 0 {
		return trimmed[idx:], ""
	}
	return "", ""
}

// valueLiteral 返回能原样读回配置文件的取值文本。
// 优先沿用原行的引号风格，其次裸文本，最后换用另一种引号包裹。
//
// 校验对象是拼好的整行（含保留下来的行内注释），不是孤立的取值：
// 注释尾巴同样参与解析，只校验取值会放过「取值本身能读回、拼上注释后读回变了」的情况。
func valueLiteral(key, v, preferQuote, suffix string) (string, error) {
	if strings.ContainsAny(v, "\r\n") {
		return "", fmt.Errorf("取值含换行，无法写进单行配置")
	}
	candidates := make([]string, 0, 3)
	if preferQuote != "" && !strings.Contains(v, preferQuote) {
		candidates = append(candidates, preferQuote+v+preferQuote)
	}
	candidates = append(candidates, v)
	for _, q := range []string{`"`, `'`} {
		if q == preferQuote || strings.Contains(v, q) {
			continue
		}
		candidates = append(candidates, q+v+q)
	}
	for _, c := range candidates {
		if lineReadsBack(key, c, suffix, v) {
			return c, nil
		}
	}
	// 不回显取值：错误信息会进入接口响应与日志，取值可能含口令或 DSN。
	return "", fmt.Errorf("取值无法写进配置文件后原样读回：请去掉其中的引号或 # 后再保存")
}

// lineReadsBack 用本库的解析器验证：按落盘格式拼出的整行读回来是否仍是 want。
// 不靠人工推导解析规则，解析规则变了这里也不会给出错误结论。
func lineReadsBack(key, literal, suffix, want string) bool {
	k, got := GetConfStrByLine(assembleValueLine(key, literal, suffix))
	return k == key && got == want
}
