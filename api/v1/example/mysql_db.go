package example

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"gorm.io/gorm"
)

type MysqlDBApi struct{}

func (api *MysqlDBApi) GetMysqlDBList(c *gin.Context) {
	MysqlDBList, err := exampleService.MysqlDBService.GetMysqlDBList()
	if err != nil {
		response.FailWithMessage("查询DB列表失败", c)
		return
	}
	response.OkWithData(MysqlDBList, c)
}

type DbWithInstances struct {
	Db        *example.MysqlDb        `json:"db"`
	Instances []example.MysqlInstance `json:"instances"`
}

func (api *MysqlDBApi) GetMysqlDbNameData(c *gin.Context) {
	db := c.Query("db")
	if db == "" {
		response.FailWithMessage("参数 db 不能为空", c)
		return
	}
	dbItem, err := exampleService.MysqlDBService.GetMysqlDbName(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailWithMessage("库名不存在", c)
		} else {
			response.FailWithMessage("查询DB失败", c)
		}
		return
	}
	instanceList, err := exampleService.MysqlInstanceService.GetClusterInstanceList(dbItem.Clustername)
	if err != nil {
		response.FailWithMessage("查询实例列表失败", c)
		return
	}
	res := DbWithInstances{
		Db:        dbItem,
		Instances: instanceList,
	}
	response.OkWithData(res, c)
}
