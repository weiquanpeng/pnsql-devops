package example

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"github.com/xiaohongshu/PnSql/server/global"
)

// 定义请求结构体（时间保持为字符串）
type PgSessionRequest struct {
	VmName    string `json:"vmname" binding:"required"`
	StartTime string `json:"startTime,omitempty"`
	StopTime  string `json:"stopTime,omitempty"`
}

type PgSessionApi struct{}

func (api *PgSessionApi) GetPgSessionList(c *gin.Context) {
	// 绑定JSON请求体
	var req PgSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Errorf("参数绑定失败: %v", err)
		response.FailWithMessage("请求参数错误: "+err.Error(), c)
		return
	}
	IP, err := exampleService.MysqlInstanceService.GetIPByVMName(req.VmName)
	sessions, err := exampleService.PgSessionService.GetPgSessionList(
		IP,
		req.StartTime, // 直接传字符串
		req.StopTime,  // 直接传字符串
	)
	if err != nil {
		//global.Logger.Errorf("查询PostgreSQL会话失败: %v", err)
		response.FailWithMessage("查询会话失败: "+err.Error(), c)
		return
	}
	response.OkWithData(sessions, c)
}
