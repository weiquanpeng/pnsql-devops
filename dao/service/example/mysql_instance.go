package example

import (
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlInstanceService struct{}

type InstanceBriefInfo struct {
	Clusternames []string `json:"clusternames"`
	Vmnames      []string `json:"vmnames"`
}

func (service *MysqlInstanceService) GetMysqlInstanceList() ([]example.MysqlInstance, error) {
	var list []example.MysqlInstance
	err := global.PVA_DB.Order("create_time").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (service *MysqlInstanceService) GetClusterInstanceList(clustername string) ([]example.MysqlInstance, error) {
	var list []example.MysqlInstance
	err := global.PVA_DB.Where("clustername = ?", clustername).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (service *MysqlInstanceService) GetAllInstanceBriefInfo() ([]example.MysqlInstance, error) {
	var list []example.MysqlInstance
	err := global.PVA_DB.Select("ip, port, vmname, clustername").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (service *MysqlInstanceService) GetAllClusterAndVmNames() (*InstanceBriefInfo, error) {
	var clusterNames []string
	var vmNames []string
	err := global.PVA_DB.Model(&example.MysqlInstance{}).
		Distinct("clustername").
		Pluck("clustername", &clusterNames).Error
	if err != nil {
		return nil, err
	}
	err = global.PVA_DB.Model(&example.MysqlInstance{}).
		Distinct("vmname").
		Pluck("vmname", &vmNames).Error
	if err != nil {
		return nil, err
	}
	return &InstanceBriefInfo{
		Clusternames: clusterNames,
		Vmnames:      vmNames,
	}, nil
}
