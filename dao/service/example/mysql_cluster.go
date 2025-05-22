package example

import (
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlClusterService struct{}

func (service *MysqlClusterService) GetMysqlClusterList() ([]example.MysqlCluster, error) {
	var list []example.MysqlCluster
	err := global.PVA_DB.Order("create_time").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (service *MysqlClusterService) GetMysqlClusterName(clustername string) (*example.MysqlCluster, error) {
	var dbItem example.MysqlCluster
	err := global.PVA_DB.
		Where("dbcluster = ?", clustername).
		First(&dbItem).Error
	if err != nil {
		return nil, err
	}
	return &dbItem, nil
}
