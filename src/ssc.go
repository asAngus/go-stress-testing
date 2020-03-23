package main

import (
	"flag"
	"ssc"
)

// 编译可执行文件
//go:generate go build ssc.go
func main() {
	var path string
	flag.StringVar(&path, "p", "", "配置文件路径")
	flag.Parse()
	ssc.Collect(path)
}
