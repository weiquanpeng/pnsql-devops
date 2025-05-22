package example

import (
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlDBService struct{}

func (service *MysqlDBService) GetMysqlDBList() ([]example.MysqlDb, error) {
	var list []example.MysqlDb
	err := global.PVA_DB.Order("create_time").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (service *MysqlDBService) GetMysqlDbName(dbname string) (*example.MysqlDb, error) {
	var dbItem example.MysqlDb
	err := global.PVA_DB.
		Where("dbname = ? OR basename = ?", dbname, dbname).
		First(&dbItem).Error
	if err != nil {
		return nil, err
	}
	return &dbItem, nil
}
