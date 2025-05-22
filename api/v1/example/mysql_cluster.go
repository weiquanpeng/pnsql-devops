package example

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"gorm.io/gorm"
)

type ClusterWithInstances struct {
	Cluster   *example.MysqlCluster   `json:"cluster"`
	Instances []example.MysqlInstance `json:"instances"`
}

type MysqlClusterApi struct{}

func (api *MysqlClusterApi) GetMysqlClusterList(c *gin.Context) {
	MysqlClusterList, err := exampleService.MysqlClusterService.GetMysqlClusterList()
	if err != nil {
		response.FailWithMessage("查询集群列表失败", c)
		return
	}
	response.OkWithData(MysqlClusterList, c)
}

func (api *MysqlClusterApi) GetMysqlClusterNameData(c *gin.Context) {
	cluster := c.Query("cluster")
	if cluster == "" {
		response.FailWithMessage("参数 cluster 不能为空", c)
		return
	}
	clusterItem, err := exampleService.MysqlClusterService.GetMysqlClusterName(cluster)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailWithMessage("集群不存在", c)
		} else {
			response.FailWithMessage("查询集群失败", c)
		}
		return
	}
	instanceList, err := exampleService.MysqlInstanceService.GetClusterInstanceList(clusterItem.Dbcluster)
	if err != nil {
		response.FailWithMessage("查询实例列表失败", c)
		return
	}
	res := ClusterWithInstances{
		Cluster:   clusterItem,
		Instances: instanceList,
	}
	response.OkWithData(res, c)
}
