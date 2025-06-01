package example

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
)

type MysqlInstanceApi struct{}

func (api *MysqlInstanceApi) GetMysqlInstanceList(c *gin.Context) {
	MysqlInstanceList, err := exampleService.MysqlInstanceService.GetMysqlInstanceList()
	if err != nil {
		response.FailWithMessage("查询实例列表失败", c)
		return
	}
	response.OkWithData(MysqlInstanceList, c)
}

func (api *MysqlInstanceApi) GetAllVmNames(c *gin.Context) {
	result, err := exampleService.MysqlInstanceService.GetVMName()
	if err != nil {
		response.FailWithMessage("查询所有实例失败", c)
		return
	}
	response.OkWithData(result, c)
}

func (api *MysqlInstanceApi) GetAllClusterAndVmNames(c *gin.Context) {
	result, err := exampleService.MysqlInstanceService.GetAllClusterAndVmNames()
	if err != nil {
		response.FailWithMessage("查询集群和实例失败", c)
		return
	}
	response.OkWithData(result, c)
}
