package example

import (
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlCluster struct {
	global.PvaModel
	Dbcluster string `json:"dbcluster" gorm:"not null;default:'';size:64;index:idx_clustername,unique;comment:物理集群名"`
	Env       string `json:"env" gorm:"not null;size:64;comment:环境"`
	Intype    string `json:"intype" gorm:"not null;default:'';size:64;comment:实例状态：online/discard/维护中"`
	Level     string `json:"level" gorm:"size:45;default:'S1';comment:集群级别"`
	Remark    string `json:"remark" gorm:"size:1000;default:'';comment:备注"`
}

func (m *MysqlCluster) TableName() string {
	return "mysql_cluster"
}
