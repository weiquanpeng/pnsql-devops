package example

import (
	"gorm.io/gorm"
	"time"
)

type PgSession struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	VMName          string     `gorm:"column:vmname;type:varchar(128);not null;comment:虚拟机名称;index:idx_vmname_captured_at" json:"vmname"`
	Source          string     `gorm:"column:source;type:varchar(64);not null;comment:实例名" json:"source"`
	PID             int32      `gorm:"column:pid;not null" json:"pid"`
	Datname         string     `gorm:"column:datname;type:varchar(64);default:'';comment:数据库名" json:"datname"`
	ApplicationName string     `gorm:"column:application_name;type:varchar(64);default:'';comment:应用名称" json:"application_name"`
	XactStart       *time.Time `gorm:"column:xact_start;type:datetime;comment:事务开始时间" json:"xact_start"`
	OpenTiming      string     `gorm:"column:open_timing;type:varchar(32);comment:事务持续时间" json:"open_timing"`
	SQLTiming       string     `gorm:"column:sql_timing;type:varchar(32);not null;comment:sql持续时间" json:"sql_timing"`
	Query           string     `gorm:"column:query;type:text;comment:执行的SQL语句" json:"query"`
	State           string     `gorm:"column:state;type:varchar(64);default:'';comment:会话状态" json:"state"`
	WaitEvent       string     `gorm:"column:wait_event;type:varchar(64);default:'';comment:等待事件类型" json:"wait_event"`
	CapturedAt      time.Time  `gorm:"column:captured_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:记录捕获时间;index:idx_vmname_captured_at" json:"captured_at"`
}

func (s *PgSession) TableName() string {
	return "pg_session"
}

func (s *PgSession) BeforeCreate(tx *gorm.DB) error {
	s.CapturedAt = time.Now()
	return nil
}
