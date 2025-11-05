package example

import (
	"errors"
	"fmt"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
	"gorm.io/gorm"
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

func (service *MysqlInstanceService) GetAllIPAndPorts() ([]struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	VMName   string `gorm:"column:vmname" json:"vmname"`
}, error) {
	var result []struct {
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		VMName   string `gorm:"column:vmname" json:"vmname"`
	}
	err := global.PVA_DB.Model(&example.MysqlInstance{}).
		Select("ip, port, username, password, vmname").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *MysqlInstanceService) GetVMName() ([]string, error) {
	var instances []example.MysqlInstance
	var vmNames []string
	err := global.PVA_DB.Select("vmname").Find(&instances).Error
	if err != nil {
		return nil, err
	}
	for _, instance := range instances {
		vmNames = append(vmNames, instance.VMName)
	}
	return vmNames, nil
}

func (service *MysqlInstanceService) GetIPByVMName(vmName string) (string, error) {
	var instance example.MysqlInstance
	err := global.PVA_DB.Select("ip").Where("vmname = ?", vmName).First(&instance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("找不到VM名称为 %s 的实例", vmName)
		}
		return "", err
	}
	return instance.IP, nil
}
