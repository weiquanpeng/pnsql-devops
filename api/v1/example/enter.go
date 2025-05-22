package example

import (
	"github.com/xiaohongshu/PnSql/server/dao/service"
)

type ExaApi struct {
	SysUserApi
	TaskConfigApi
	SubTaskConfigApi
	MysqlClusterApi
	MysqlDBApi
	MysqlInstanceApi
	MysqlSessionApi
}

var (
	exampleService = service.GroupApp.ExampleServer
)
