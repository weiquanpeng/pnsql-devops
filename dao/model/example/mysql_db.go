package example

import (
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlDb struct {
	global.PvaModel
	Basename    string `json:"basename" gorm:"not null;default:'0';size:255;comment:标注不同环境的同一个db，解决同名db问题，basename+env唯一"`
	Dbname      string `json:"dbname" gorm:"not null;size:64;comment:同一集群不同分片的同名db只插一条数据"`
	Env         string `json:"env" gorm:"not null;size:64;comment:环境"`
	Clustername string `json:"cluster_name" gorm:"not null;size:64;comment:对应cluster表dbcluster"`
	Dbowner     string `json:"dbowner" gorm:"not null;default:'';size:256;comment:dbowner"`
	Level       string `json:"level" gorm:"not null;default:'S2';size:32;comment:数据库级别：S0/S1/S2/S3"`
	Intype      string `json:"intype" gorm:"not null;default:'online';size:32;comment:online/discard/维护中"`
	Remark      string `json:"remark" gorm:"size:1000;default:'';comment:备注"`
}

func (m *MysqlDb) TableName() string {
	return "mysql_db"
}
