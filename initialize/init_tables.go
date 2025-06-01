package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/dao/model/system"
	"github.com/xiaohongshu/PnSql/server/global"
	"github.com/xiaohongshu/PnSql/server/utils"
	"os"
)

func RegisterTables() {
	db := global.PVA_DB
	err := db.AutoMigrate(
		// 系统模块表
		system.Jwt{},
		// 业务模块表
		example.TaskConfig{},
		example.SubTaskConfig{},
		example.ProcessSubtaskConfig{},
		example.SysUser{},
		example.MysqlCluster{},
		example.MysqlDb{},
		example.MysqlInstance{},
		example.MysqlSession{},
		example.PgSession{},
	)
	if err != nil {
		fmt.Println("register table err:", err)
		os.Exit(0)
	}

	//初始化数据
	err = insertDefaultData()
	if err != nil {
		fmt.Println("insert default data err:", err)
		os.Exit(0)
	}
}

func insertDefaultData() error {
	//初始化 admin 用户
	var user_count int64
	global.PVA_DB.Model(&example.SysUser{}).Where("account = ?", "admin").Count(&user_count)
	if user_count == 0 {
		admin := example.SysUser{
			Account:  "admin",
			Password: utils.Md5("123123"),
			Enable:   1,
			Roles:    json.RawMessage(`[1, 2, 3]`),
		}
		result := global.PVA_DB.Create(&admin)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
