package example

import (
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
	"time"
)

type MysqlSessionService struct{}

func (service *MysqlSessionService) GetMysqlSessionListByFilter(clustername, vmname string, startTime, endTime time.Time) ([]example.MysqlSession, error) {
	var list []example.MysqlSession
	db := global.PVA_DB.Model(&example.MysqlSession{})

	if clustername != "" {
		db = db.Where("clustername = ?", clustername)
	}
	if vmname != "" {
		db = db.Where("vmname = ?", vmname)
	}
	if !startTime.IsZero() && !endTime.IsZero() {
		db = db.Where("executed_time BETWEEN ? AND ?", startTime, endTime)
	}

	err := db.Order("TIME DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
