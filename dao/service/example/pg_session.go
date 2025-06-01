package example

import (
	"fmt"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
)

type PgSessionService struct{}

func (service *PgSessionService) GetPgSessionList(ip string, startTime string, stopTime string) ([]example.PgSession, error) {
	var list []example.PgSession

	// 构建基础查询
	query := global.PVA_DB.Where("source = ?", ip)

	// 使用字符串直接构建SQL条件
	if startTime != "" && stopTime != "" {
		query = query.Where("captured_at BETWEEN ? AND ?", startTime, stopTime)
	} else if startTime != "" {
		query = query.Where("captured_at >= ?", startTime)
	} else if stopTime != "" {
		query = query.Where("captured_at <= ?", stopTime)
	}

	// 执行查询
	err := query.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 在PgSessionService中添加以下方法
func (service *PgSessionService) BatchInsertSessions(sessions []example.PgSession) error {
	if len(sessions) == 0 {
		return nil
	}

	// 使用批量插入提高性能
	err := global.PVA_DB.CreateInBatches(sessions, 100).Error
	if err != nil {
		return fmt.Errorf("批量插入会话失败: %w", err)
	}
	return nil
}
