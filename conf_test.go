package easyconf

import (
	"testing"
)

func TestGetConf(t *testing.T) {
	cf := NewConf(".env", "default1.env", "defalut2.env")
	var version string
	var isauto bool
	var webport int
	var domains []string
	domainUsage := []string{
		"1. 允许直连的域名放在 routing.rules 数组中",
		"2. 当路由匹配到一个规则时就会跳出匹配而不会对之后的规则进行匹配；",
	}
	cf.StringVar(&version, "VERSION", "v1.0.0", "版本号")
	cf.BoolVar(&isauto, "IS_AUTO", false, "开启自动")
	cf.IntVar(&webport, "WEB_PORT", 8080, "web服务器端口")
	cf.StringListVar(&domains, "DOMAINS", []string{"baidu.com", "taobao.com"}, "允许的域名列表", domainUsage...)

	webport = 88898745
	version = "v1.23.9"
	err := cf.Parse()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("---VERSION(%s)--IS_AUTO(%t)--WEB_PORT(%d)--DOMAINS(%v)---\n", version, isauto, webport, domains)
}
