# AGENTS.md - easyconf

## 项目概述

`easyconf` 是一个 Go 配置管理库，用于从 `.env` 格式的配置文件、系统环境变量和命令行参数中读取配置。模块路径：`github.com/iotames/easyconf`，依赖 `github.com/iotames/miniutils`（仅用于 `IsPathExists` 工具函数）。

## 常用命令

```bash
# 运行测试
go test ./...

# 运行指定测试函数
go test -run TestConf -v ./...
go test -run TestConfLine -v ./...
```

没有 Makefile、CI 配置或 lint 工具。Go 1.22.1+。

## 架构与数据流

### 核心结构

```
Conf (conf.go)
  ├── files []string        # 配置文件路径列表，优先级从左到右递减
  └── items []*ConfItem     # 注册的配置项

ConfItem (conf_item.go)
  ├── Name, Title, Usage    # 元信息（用于生成注释）
  ├── Value any             # 指针，运行时值，引用传递
  └── DefaultValue any      # 默认值，值传递
```

### 配置优先级（高 → 低）

1. **命令行参数** — `SetValuesByCmdArgs()` 调用 `flag.Parse()`，与 Go 原生 flag 库集成
2. **系统环境变量** — `SetValuesByEnv()` 读取 `os.Getenv()`
3. **配置文件列表** — `SetValuesByEnvFile()` 逆序读取，前面的文件覆盖后面的

### Parse 流程 (`conf_parse.go`)

```
Parse(flagParse bool)
  ├── 1. 遍历 cf.files，不存在的文件用 DefaultString() 创建
  ├── 2. 逆序遍历配置文件，逐行解析 KEY=VALUE 并赋值
  ├── 3. 从 os.Getenv() 覆盖同名配置
  └── 4. 若 flagParse=true，通过 flag.Parse() 用命令行参数覆盖
```

### 文件职责

| 文件 | 职责 |
|------|------|
| `conf.go` | `Conf` 结构体、`NewConf()`、`DefaultString()`/`String()`、`createEnvFile()` |
| `conf_define.go` | 类型化配置定义方法：`StringVar`、`IntVar`、`BoolVar`、`StringListVar`、`IntListVar`、`Float64Var` |
| `conf_item.go` | `ConfItem` 结构体、设置/获取值、`anyToString`、`.env` 文件内容格式化 |
| `conf_parse.go` | `Parse()` 主流程、`GetConfStrByLine()` 行解析器 |
| `conf_update.go` | `addItem()`、`setItemVar()`、`SetValuesByEnvFile()`、`SetValuesByEnv()`、`SetValuesByCmdArgs()`、`UpdateFile()`、`AddComment()` |
| `jsonfile.go` | 独立的 `JsonConf` 类型，提供 JSON 文件的读写（与 `Conf` 无关） |
| `atomicfile.go` | `writeFileAtomic()` 原子写（临时文件 + rename + 目录 fsync），`UpdateFile` 与 `JsonConf.Save` 共用 |
| `conf_test.go` | 行解析测试 + 完整读写更新流程测试 |

## 关键约定与注意事项

### `.env` 文件格式

- 逐行解析，格式为 `KEY = VALUE`
- `#` 开头的行是注释
- 值可以用单引号或双引号包裹，内部含 `#` 时会被正确处理
- 列表类型（`StringListVar`/`IntListVar`）的值用逗号分隔，如 `AGE_RANGE = 3,6`

### 写回约定（`UpdateFile` / `JsonConf.Save`）

- **原子写**：同目录临时文件 → `fsync` → rename 替换；POSIX 上再 `fsync` 父目录。两者共用 `writeFileAtomic()`。
- **权限**：已有文件沿用原权限，新建文件 `0644`。
- **写前自检**：每一行都要能被 `GetConfStrByLine()` 原样读回；取值含 `#`/引号时自动补引号，无法表达的取值报错且不改动原文件。校验对象是拼好的整行（含保留的行内注释）。
- **不回显取值**：错误信息只报键名与原因，取值可能含口令/DSN。
- 重建行的缩进与等号两侧间距会被归一化；无法识别的尾随内容（非 `#` 注释）会被丢弃。

### 配置优先级实现细节

- 逆序遍历配置文件列表：最后一个文件先读，第一个文件最后读，实现"前面覆盖后面"
- `SetValuesByEnv()` 中**空字符串会被跳过**，不会覆盖已有值
- `SetValuesByEnvFile()` 中**空键或空值会被跳过**

### `flag.Parse()` 冲突

`Parse(flagParse=true)` 内部会调用 `flag.Parse()`。如果调用方自己也调用了 `flag.Parse()`（如 `conf_test.go:46` 的示例中注释掉的 `flag.Parse()`），会发生冲突。调用方应删除自己的 `flag.Parse()`，由 `Parse(true)` 统一调用。

### 类型断言是 `any`

`ConfItem.Value` 和 `ConfItem.DefaultValue` 都是 `any` 类型，所有类型判断通过 type switch 完成。支持的类型：`*int`、`*float64`、`*bool`、`*string`、`*[]string`、`*[]int`。不支持的类型会 panic 或返回错误。

### 测试文件会在当前目录写入 `.env` 和 `default.env`

测试运行时会在项目根目录创建/读写 `.env` 文件。测试最后会清理这两个文件（`os.Remove`）。不要在这些文件中存放重要数据。

### `AddComment()` 方法

可以添加纯注释行到配置输出中：`cf.AddComment("标题", "注释行1", "注释行2")`。内部仍然创建 `ConfItem`，但 `Name` 为空，解析时会被跳过。

### `UpdateFile()` 注意事项

- 留空 `fpath` 时默认更新 `cf.files[0]`（第一个配置文件）
- 使用 `O_TRUNC` 打开，会**完全覆盖**文件内容
- 写入内容由 `cf.String()` 生成（当前值，非默认值）

### `JsonConf` 是独立组件

`jsonfile.go` 中的 `JsonConf` 与 `Conf` 类型无关，是一个简单的 JSON 文件读写工具。`filename` 参数不包含目录路径，目录由构造时的 `dirPath` 指定。

### 依赖

- `github.com/iotames/miniutils` — 仅用于 `IsPathExists()` 判断文件是否存在
- `golang.org/x/net` 和 `golang.org/x/text` 为间接依赖
