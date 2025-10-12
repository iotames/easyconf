package easyconf

import (
	"fmt"
	"os"
	"testing"
)

func TestConfLine(t *testing.T) {
	lines := []string{
		`NAME0 =VALUE0`,
		`NAME1=VALUE1`,
		`NAME2= VALUE2`,
		`NAME3  = 'VALUE3'`,
		`NAME4 =  "VALUE4"`,
		`NAME5= 'VALUE5'`,
		`NAME6 ="VALUE6"`,
		`NAME7 ="VALUE7" # 备注1`,
		`NAME8 = 'VALUE8' # 备注2`,
		`NAME9 = "NAME9=VALUE9" # 备注3`,
		`NAME10='NAME10=VALUE10' # 备注4`,
		`NAME11 =NAME11=VALUE11 # 备注5`,
		`NAME12 ="NAME12#VALUE12" # 备注6`,
		`NAME13="NAME13#VALUE13" # remark="a=b"`,
	}
	for i, line := range lines {
		okval := ""
		k, v := GetConfStrByLine(line)
		if i < 9 {
			okval = fmt.Sprintf(`VALUE%d`, i)
		}
		if i >= 9 && i <= 11 {
			okval = fmt.Sprintf("NAME%d=VALUE%d", i, i)
		}
		if i >= 12 {
			okval = fmt.Sprintf("NAME%d#VALUE%d", i, i)
		}
		if k != fmt.Sprintf(`NAME%d`, i) || v != okval {
			t.Fatal(fmt.Errorf("value(%s) err for %s", v, okval))
		}
	}
}

func TestConf(t *testing.T) {
	cf := NewConf()
	var version string
	var isbool1, isbool2 bool
	var webport int
	var domains []string
	var intlist []int
	var costDays float64
	domainUsage := []string{
		"1. 允许直连的域名放在 routing.rules 数组中",
		"2. 当路由匹配到一个规则时就会跳出匹配而不会对之后的规则进行匹配；",
	}
	version = "v1.0.1"
	const DEFAULT_VERSION = "v1.0.0"
	const DEFAULT_WEBPORT = 8080
	var DEFAULT_DOMAINS = []string{"baidu.com", "taobao.com"}
	var DEFAULT_INTLIST = []int{2, 4, 6, 8}
	var DEFAULT_COST_DAYS = 3.25
	cf.StringVar(&version, "VERSION", DEFAULT_VERSION, "版本号")
	cf.BoolVar(&isbool1, "IS_BOOL1", false, "默认关闭")
	cf.BoolVar(&isbool2, "IS_BOOL2", true, "默认开启")
	cf.IntVar(&webport, "WEB_PORT", DEFAULT_WEBPORT, "web服务器端口")
	cf.AddComment("这个是注释说明", "继续添加注释说明A")
	cf.StringListVar(&domains, "DOMAINS", DEFAULT_DOMAINS, "允许的域名列表", domainUsage...)
	cf.IntListVar(&intlist, "INTLIST", DEFAULT_INTLIST, "整数列表")
	cf.Float64Var(&costDays, "COST_DAYS", DEFAULT_COST_DAYS, "耗时天数", "完成工期花费的时间/天")
	webport = 8888
	err := cf.Parse(false)
	if err != nil {
		t.Fatal(err)
	}

	// 验证默认值
	if version != DEFAULT_VERSION || isbool1 || !isbool2 || webport != DEFAULT_WEBPORT {
		t.Fatal(fmt.Errorf(`默认值设置错误isbool1(%t)--isbool2(%t)`, isbool1, isbool2))
	}
	if costDays != DEFAULT_COST_DAYS {
		t.Fatal(fmt.Errorf("COST_DAYS err--val(%v)---default(%v)", costDays, DEFAULT_COST_DAYS))
	}

	for i, d := range DEFAULT_DOMAINS {
		if d != domains[i] {
			t.Fatal("[]string 默认值设置错误")
		}
	}

	for i, d := range DEFAULT_INTLIST {
		if d != intlist[i] {
			t.Fatal("[]int 默认值设置错误")
		}
	}

	// t.Logf("---111--VERSION(%s)--IS_BOOL1(%t)--WEB_PORT(%d)--DOMAINS(%v)---\n", version, isbool1, webport, domains)

	// 更新测试
	webport = 8899
	version = "v1.99.9"
	isbool1 = true
	isbool2 = false
	domains = []string{"amazon.com", "google.com"}
	intlist = []int{1, 3, 7}
	costDays = 6.79
	err = cf.UpdateFile("")
	if err != nil {
		t.Fatal(err)
	}
	err = cf.UpdateFile("update.env")
	if err != nil {
		t.Fatal(err)
	}

	// 验证更新
	updatedWebport := webport
	updatedVersion := version
	updatedDomains := domains
	updatedIntlist := intlist
	updateCostDays := costDays

	cf = NewConf()
	cf.StringVar(&version, "VERSION", DEFAULT_VERSION, "版本号")
	cf.BoolVar(&isbool1, "IS_BOOL1", false, "默认关闭")
	cf.BoolVar(&isbool2, "IS_BOOL2", true, "默认开启")
	cf.IntVar(&webport, "WEB_PORT", DEFAULT_WEBPORT, "web服务器端口")
	cf.StringListVar(&domains, "DOMAINS", DEFAULT_DOMAINS, "允许的域名列表", domainUsage...)
	cf.IntListVar(&intlist, "INTLIST", DEFAULT_INTLIST, "整数列表")
	cf.Float64Var(&costDays, "COST_DAYS", DEFAULT_COST_DAYS, "耗时天数", "完成工期花费的时间/天")
	webport = 8888
	err = cf.Parse(false)
	if err != nil {
		t.Fatal(err)
	}
	if version != updatedVersion || !isbool1 || isbool2 || webport != updatedWebport {
		t.Fatal(fmt.Errorf(`配置更新验证失败isbool1(%t)--isbool2(%t)`, isbool1, isbool2))
	}
	if costDays != updateCostDays {
		t.Fatal(fmt.Errorf("COST_DAYS err--costDays-val(%v)---expected(%v)", costDays, updateCostDays))
	}
	for i, d := range updatedDomains {
		if d != domains[i] {
			t.Fatal("[]string 更新值设置错误")
		}
	}
	for i, d := range updatedIntlist {
		if d != intlist[i] {
			t.Fatal("[]int 更新值设置错误")
		}
	}

	// 验证默认值

	// t.Logf("--222--VERSION(%s)--IS_BOOL1(%t)--WEB_PORT(%d)--DOMAINS(%v)---\n", version, isbool1, webport, domains)
	f, err := os.OpenFile(".env", os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("")
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	cf = NewConf()
	cf.StringVar(&version, "VERSION", DEFAULT_VERSION, "版本号")
	cf.BoolVar(&isbool1, "IS_BOOL1", false, "默认关闭")
	cf.BoolVar(&isbool2, "IS_BOOL2", true, "默认开启")
	cf.IntVar(&webport, "WEB_PORT", DEFAULT_WEBPORT, "web服务器端口")
	cf.StringListVar(&domains, "DOMAINS", DEFAULT_DOMAINS, "允许的域名列表", domainUsage...)
	cf.IntListVar(&intlist, "INTLIST", DEFAULT_INTLIST, "整数列表")
	err = cf.Parse(false)
	if err != nil {
		t.Fatal(err)
	}
	// 验证默认值
	if version != DEFAULT_VERSION || isbool1 || !isbool2 || webport != DEFAULT_WEBPORT {
		t.Fatal(fmt.Errorf(`默认值设置错误isbool1(%t)--isbool2(%t)`, isbool1, isbool2))
	}
	for i, d := range DEFAULT_DOMAINS {
		if d != domains[i] {
			t.Fatal("[]string 默认值设置错误")
		}
	}
	for i, d := range DEFAULT_INTLIST {
		if d != intlist[i] {
			t.Fatal("[]int 默认值设置错误")
		}
	}
	os.Remove(".env")
	os.Remove("default.env")
}
