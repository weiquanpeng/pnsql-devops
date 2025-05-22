package main

import (
	"github.com/xiaohongshu/PnSql/server/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitRouter()
}
