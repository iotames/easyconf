## 简介

简单的Go应用程序配置工具，用于生成和读取 `env` 配置文件的环境变量。

可以很方便地定义 `配置名`，`默认值`，`注释标题`，`注释说明`（多行）。

通过简单的方法读取和修改配置值。

优先级：`命令行同名参数值` > `环境变量值` > `配置文件配置值`


## 新建配置


导入 `easyconf` 工具包:

```
import (
	"github.com/iotames/easyconf"
)
```

可初始化1个或多个配置文件，优先级从左到右。默认: `.env`, `default.env`

```
// 生成 .env, default.env 两份配置文件。
cf := easyconf.NewConf()

// 生成一份 myconf.env 自定义配置文件。
cf = easyconf.NewConf("myconf.env")

// 生成多份配置文件
cf = easyconf.NewConf(".env", "common.env", "default.env")
```

使用 `Parse(fale)` 方法读取文件中的配置值。若文件不存在，则创建。

```
var DbHost string
var DbPort int
cf := easyconf.NewConf()
cf.StringVar(&DbHost, "DB_HOST", "127.0.0.1", "数据库主机地址")
cf.IntVar(&DbPort, "DB_PORT", 3306, "数据库地址端口号")
cf.Parse(false) // 默认创建 .env, default.env 两份文件。
```

## 更新配置

使用 `UpdateFile()` 方法更新配置。需要指定配置文件路径，留空则默认更新第一个。

```
cf := easyconf.NewConf(".env")
var DbHost string
cf.StringVar(&DbHost, "DB_HOST", "127.0.0.1", "数据库主机地址")
cf.Parse(false)

DbHost = "192.168.1.6"

// 可指定一个更新的配置文件。留空则默认更新第一个。
err := cf.UpdateFile("")
if err != nil {
	panic(err)
}
```

## 完整示例

```
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
	fmt.Printf("-----DbHost(%s)--DbPort(%d)--Dbnable(%t)--\n", DbHost, DbPort, DbEnable)
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
	flag.Parse()
	cf = easyconf.NewConf(getEnvFile())
	cf.StringVar(&DbHost, "DB_HOST", DEFAULT_DB_HOST, "数据库主机地址")
	cf.IntVar(&DbPort, "DB_PORT", DEFAULT_DB_PORT, "数据库地址端口号")
	cf.BoolVar(&DbEnable, "DB_ENABLE", false, "是否启用数据库")
	cf.StringListVar(&AllowIPs, "ALLOW_IP_LIST", []string{"192.168.1.6", "192.168.2.8"}, "允许访问的IP列表，每个IP用逗号(,)隔开。")
	cf.IntListVar(&AgeRange, "AGE_RANGE", []int{3, 6}, "年龄范围", "填写2个正整数,中间用逗号,隔开")
	cf.Parse(false)
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


```