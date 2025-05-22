package example

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"time"
)

type MysqlSessionApi struct{}

type SessionRequest struct {
	ClusterName  string `json:"clustername"`
	InstanceName string `json:"instancename"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

func (api *MysqlSessionApi) GetMysqlSessionList(c *gin.Context) {
	const timeLayout = "2006-01-02 15:04:05" // 定义时间格式

	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数格式错误", c)
		return
	}

	// 时间解析处理
	var startTime, endTime time.Time
	var err error

	if req.StartTime != "" {
		startTime, err = time.ParseInLocation(timeLayout, req.StartTime, time.Local)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("开始时间格式错误，需要格式: %s", timeLayout), c)
			return
		}
	}

	if req.EndTime != "" {
		endTime, err = time.ParseInLocation(timeLayout, req.EndTime, time.Local)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("结束时间格式错误，需要格式: %s", timeLayout), c)
			return
		}
	}

	// 调用服务层
	MysqlSessionList, err := exampleService.MysqlSessionService.GetMysqlSessionListByFilter(
		req.ClusterName,
		req.InstanceName,
		startTime,
		endTime,
	)

	if err != nil {
		response.FailWithMessage("查询会话列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(MysqlSessionList, c)
}
