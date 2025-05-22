package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "github.com/xiaohongshu/PnSql/server/api/v1"
	"github.com/xiaohongshu/PnSql/server/common/filter"
	"github.com/xiaohongshu/PnSql/server/common/jwt"
	"github.com/xiaohongshu/PnSql/server/global"
)

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(filter.Cors())
	// 公共路由组 - 不需要认证
	publicGroup := r.Group("/api/v1/public")
	{
		// 登录相关
		publicGroup.GET("/login/get", v1.ApiGroupApp.SystemApiGroup.Login.Login)
		publicGroup.POST("/login/post", v1.ApiGroupApp.SystemApiGroup.Login.ToLogin)
		publicGroup.POST("/logout/post", v1.ApiGroupApp.SystemApiGroup.Logout.Logout)
	}

	// 私有路由组 - 需要JWT认证
	privateGroup := r.Group("/api/v1")
	privateGroup.Use(jwt.JwtAuth())
	{
		// 用户管理
		userGroup := privateGroup.Group("/sys/user")
		{
			userGroup.GET("/post/:id", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.GetById)
			userGroup.POST("/addUser", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.AddSysUser)
			userGroup.POST("/uptUser", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.UptSysUser)
			userGroup.POST("/uptFiled", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.UptSysUserField)
			userGroup.POST("/delUser", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.DelSysUser)
			userGroup.POST("/load", v1.ApiGroupApp.ExampleApiGroup.SysUserApi.LoadSysUserPage)
		}
		// 任务管理
		taskGroup := privateGroup.Group("/task")
		{
			taskGroup.GET("/list", v1.ApiGroupApp.ExampleApiGroup.TaskConfigApi.LoadTaskConfigPage)
			taskGroup.GET("/subtask/list", v1.ApiGroupApp.ExampleApiGroup.SubTaskConfigApi.GetSubTaskConfigData)
			taskGroup.POST("/mine/list", v1.ApiGroupApp.ExampleApiGroup.TaskConfigApi.LoadOwnerTaskConfigPage)
			taskGroup.POST("/approve/list", v1.ApiGroupApp.ExampleApiGroup.SubTaskConfigApi.LoadSubTaskConfigPage)
			taskGroup.POST("/log/list", v1.ApiGroupApp.ExampleApiGroup.TaskConfigApi.LoadOwnerTaskConfigPage)
			taskGroup.POST("/approve/uptStatus", v1.ApiGroupApp.ExampleApiGroup.SubTaskConfigApi.UptSubTaskConfigData)
		}
		mysqlGroup := privateGroup.Group("/mysql")
		{
			mysqlGroup.GET("/cluster", v1.ApiGroupApp.ExampleApiGroup.MysqlClusterApi.GetMysqlClusterList)
			mysqlGroup.GET("/clustername", v1.ApiGroupApp.ExampleApiGroup.MysqlClusterApi.GetMysqlClusterNameData)
			mysqlGroup.GET("/db", v1.ApiGroupApp.ExampleApiGroup.MysqlDBApi.GetMysqlDBList)
			mysqlGroup.GET("/dbname", v1.ApiGroupApp.ExampleApiGroup.MysqlDBApi.GetMysqlDbNameData)
			mysqlGroup.GET("/instance", v1.ApiGroupApp.ExampleApiGroup.MysqlInstanceApi.GetMysqlInstanceList)
			mysqlGroup.GET("/instance/list", v1.ApiGroupApp.ExampleApiGroup.MysqlInstanceApi.GetAllClusterAndVmNames)
			mysqlGroup.POST("/session", v1.ApiGroupApp.ExampleApiGroup.MysqlSessionApi.GetMysqlSessionList)
		}
	}

	// 启动服务
	stPort := global.P_cfg.System.Port
	if stPort == "" {
		stPort = "2379"
	}
	err := r.Run(fmt.Sprintf("0.0.0.0:%s", stPort))
	if err != nil {
		panic(fmt.Sprintf("start service error: %s", err.Error()))
	}
}
