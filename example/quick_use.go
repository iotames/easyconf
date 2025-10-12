package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iotames/easyconf"
)

// 可修改 envfile 文件中配置项的值，验证值是否更新
const DEFAULT_ENV_FILE = ".env"

const DEFAULT_DB_HOST = "127.0.0.1"
const DEFAULT_DB_PORT = 3306

var envfile string

var cf *easyconf.Conf

var DbHost string
var DbPort int
var DbEnable bool
var AllowIPs []string
var AgeRange []int

func main() {
	fmt.Printf("-------可修改envfile(%s)文件中配置项的默认值，验证修改值是否更新-----\n", getEnvFile())
	fmt.Printf("-----DbHost(%s)--DbPort(%d)--DbEnable(%t)--\n", DbHost, DbPort, DbEnable)
	fmt.Printf("-----AllowIPs(%+v)--AgeRange(%v)--\n", AllowIPs, AgeRange)
	if DbHost == DEFAULT_DB_HOST {
		DbHost = "192.168.1.19"
		DbPort = 3308
		err := cf.UpdateFile(getEnvFile())
		if err != nil {
			panic(err)
		}
		fmt.Printf("SUCCESS: Update Env File:%s\n", getEnvFile())
	}
}

func init() {
	flag.StringVar(&envfile, "envfile", "", "配置文件路径")
	// flag.Parse()
	cf = easyconf.NewConf(getEnvFile())
	cf.StringVar(&DbHost, "DB_HOST", DEFAULT_DB_HOST, "数据库主机地址")
	cf.IntVar(&DbPort, "DB_PORT", DEFAULT_DB_PORT, "数据库地址端口号")
	cf.BoolVar(&DbEnable, "DB_ENABLE", false, "是否启用数据库")
	cf.StringListVar(&AllowIPs, "ALLOW_IP_LIST", []string{"192.168.1.6", "192.168.2.8"}, "允许访问的IP列表，每个IP用逗号(,)隔开。")
	cf.IntListVar(&AgeRange, "AGE_RANGE", []int{3, 6}, "年龄范围", "填写2个正整数,中间用逗号,隔开")
	err := cf.Parse(true)
	if err != nil {
		panic(err)
	}
}

func getEnvFile() string {
	efile := envfile
	if efile == "" {
		efile = os.Getenv("ENV_FILE")
	}
	if efile == "" {
		efile = DEFAULT_ENV_FILE
	}
	return efile
}
