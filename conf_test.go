package easyconf

import (
	"fmt"
	"testing"
)

func TestConfLine(t *testing.T) {
	lines := []string{
		`NAME1=VALUE1`,
		`NAME2 = VALUE2`,
		`NAME3  = 'VALUE3'`,
		`NAME4 =  "VALUE4"`,
		`NAME5 ='VALUE5'`,
		`NAME6 ="VALUE6"`,
		`NAME7 ="VALUE7" # 备注1`,
		`NAME8 ='VALUE8' # 备注2`,
	}
	for i, line := range lines {
		k, v := GetConfStrByLine(line)
		if k != fmt.Sprintf(`NAME%d`, i+1) || v != fmt.Sprintf(`VALUE%d`, i+1) {
			t.Fatal("func GetConfStrByLine test err")
		}
	}
}

func TestConf(t *testing.T) {
	cf := NewConf(".env", "default1.env", "defalut2.env")
	var version string
	var isauto bool
	var webport int
	var domains []string
	domainUsage := []string{
		"1. 允许直连的域名放在 routing.rules 数组中",
		"2. 当路由匹配到一个规则时就会跳出匹配而不会对之后的规则进行匹配；",
	}
	version = "v1.0.1"
	cf.StringVar(&version, "VERSION", "v1.0.0", "版本号")
	cf.BoolVar(&isauto, "IS_AUTO", false, "开启自动")
	cf.IntVar(&webport, "WEB_PORT", 8080, "web服务器端口")
	webport = 8888
	cf.StringListVar(&domains, "DOMAINS", []string{"baidu.com", "taobao.com"}, "允许的域名列表", domainUsage...)
	err := cf.Parse()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("---111--VERSION(%s)--IS_AUTO(%t)--WEB_PORT(%d)--DOMAINS(%v)---\n", version, isauto, webport, domains)
	webport = 8899
	version = "v1.99.9"
	isauto = true
	domains = []string{"amazon.com", "google.com"}
	// cf.UpdateFile("")
	// cf.UpdateFile("update.env")
	// t.Logf("--222--VERSION(%s)--IS_AUTO(%t)--WEB_PORT(%d)--DOMAINS(%v)---\n", version, isauto, webport, domains)
}
