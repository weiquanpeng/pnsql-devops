package example

import (
	"gorm.io/gorm"
	"time"
)

// PgSQLStat 用于存储每条 SQL 语句的调用次数（Calls）、采样时间（CapturedAt）、以及计算的 QPS
type PgSQLStat struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Source     string    `gorm:"column:source;type:varchar(64);not null;comment:实例名;index:idx_source_captured_at" json:"source"`
	Query      string    `gorm:"type:text;not null;comment:SQL 语句" json:"query"`
	Calls      int64     `gorm:"not null;comment:调用次数" json:"calls"`
	QPS        float64   `gorm:"not null;default:0;comment:每秒请求数（计算得出）" json:"qps"`
	CapturedAt time.Time `gorm:"type:datetime;not null;index:idx_captured_at;comment:采样时间" json:"captured_at"`
}

func (s *PgSQLStat) TableName() string {
	return "pg_sql_stats"
}

func (s *PgSQLStat) BeforeCreate(tx *gorm.DB) error {
	s.CapturedAt = time.Now()
	return nil
}
