package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fghwett/heytap/config"
)

var path = flag.String("path", "./config.yml", "配置文件地址")

func main() {
	flag.Parse()

	_, err := config.Init(*path)
	if err != nil {
		fmt.Printf("读取配置文件失败 err: %s", err)
		os.Exit(-1)
	}

}
