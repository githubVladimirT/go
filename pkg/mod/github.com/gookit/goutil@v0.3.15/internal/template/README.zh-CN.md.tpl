# Go Util

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/goutil?style=flat-square)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/goutil)](https://github.com/gookit/goutil)
[![GoDoc](https://godoc.org/github.com/gookit/goutil?status.svg)](https://pkg.go.dev/github.com/gookit/goutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/goutil)](https://goreportcard.com/report/github.com/gookit/goutil)
[![Unit-Tests](https://github.com/gookit/goutil/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/goutil/actions)
[![Coverage Status](https://coveralls.io/repos/github/gookit/goutil/badge.svg?branch=master)](https://coveralls.io/github/gookit/goutil?branch=master)

Go一些常用的工具函数收集、实现和整理

- `arrutil` array/slice 相关操作的函数工具包
- `dump`  简单的变量打印工具，打印 slice, map 会自动换行显示每个元素，同时会显示打印调用位置
- `cliutil` CLI 的一些工具函数包
- `envutil` ENV 信息获取判断工具包
- `fmtutil` format data tool
- `fsutil` 文件系统操作相关的工具函数包
- `jsonutil` JSON util
- `maputil` map 相关操作的函数工具包
- `mathutil` int/number 相关操作的函数工具包
- `strutil` string 相关操作的函数工具包
- `sysutil` system 相关操作的函数工具包
- `testutil` test help 相关操作的函数工具包

> **[EN README](README.md)**

## GoDoc

- [Godoc for github](https://pkg.go.dev/github.com/gookit/goutil)

## Util Packages
{{pgkFuncs}}
## Code Check & Testing

```bash
gofmt -w -l ./
golint ./...
go test ./...
```

## Gookit packages

  - [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
  - [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP
  - [gookit/gcli](https://github.com/gookit/gcli) Build CLI application, tool library, running CLI commands
  - [gookit/slog](https://github.com/gookit/slog) Lightweight, easy to extend, configurable logging library written in Go
  - [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
  - [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
  - [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
  - [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
  - [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
  - [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
  - [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
  - More, please see https://github.com/gookit

## License

[MIT](LICENSE)
