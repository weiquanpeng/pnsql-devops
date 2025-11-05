package example

import (
	"github.com/xiaohongshu/PnSql/server/global"
)

type MysqlInstance struct {
	global.PvaModel
	IP          string `gorm:"column:ip;type:varchar(128);not null" json:"ip"`
	Port        string `gorm:"column:port;type:varchar(64);not null" json:"port"`
	Username    string `gorm:"column:username;type:varchar(256);not null" json:"username"`
	Password    string `gorm:"column:password;type:varchar(128);not null" json:"password"`
	VMName      string `gorm:"column:vmname;type:varchar(128);not null;unique" json:"vmname"`
	Env         string `gorm:"column:env;type:varchar(64);not null" json:"env"`
	ClusterName string `gorm:"column:clustername;type:varchar(64);not null" json:"clustername"`
	Role        string `gorm:"column:role;type:varchar(64);not null" json:"role"`
	DNS         string `gorm:"column:dns;type:varchar(256)" json:"dns,omitempty"`
	DBVersion   string `gorm:"column:dbversion;type:varchar(64);not null" json:"dbversion"`
	InType      string `gorm:"column:intype;type:varchar(10);not null;default:''" json:"intype"`
	Source      string `gorm:"column:source;type:varchar(64);not null;default:'自建'" json:"source"`
	Shared      string `gorm:"column:shared;type:varchar(10);not null;default:''" json:"shared"`
}

func (MysqlInstance) TableName() string {
	return "mysql_instance"
}
