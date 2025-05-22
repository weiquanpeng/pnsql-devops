package example

import "time"

type MysqlSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Clustername   string    `json:"cluster_name" gorm:"not null;size:64;comment:对应cluster表dbcluster"`
	VMName        string    `gorm:"column:vmname;type:varchar(64);not null;unique" json:"vmname"`
	User          string    `json:"user" gorm:"column:USER;type:varchar(32);not null;default:'';comment:用户名"`
	Host          string    `json:"host" gorm:"column:HOST;type:varchar(261);not null;default:'';comment:主机"`
	DB            string    `json:"db" gorm:"column:DB;type:varchar(64);default:null;comment:数据库"`
	Command       string    `json:"command" gorm:"column:COMMAND;type:varchar(16);not null;default:'';comment:命令"`
	Time          int       `json:"time" gorm:"column:TIME;type:int;not null;default:0;comment:持续时间"`
	State         string    `json:"state" gorm:"column:STATE;type:varchar(64);default:null;comment:状态"`
	Info          string    `json:"info" gorm:"column:INFO;type:longtext;comment:详细信息"`
	ExecutionPlan string    `json:"execution_plan" gorm:"column:execution_plan;type:text;comment:执行计划"`
	ExecutedTime  time.Time `json:"executed_time" gorm:"column:executed_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:执行的时间点"`
}

func (s *MysqlSession) TableName() string {
	return "mysql_session"
}
