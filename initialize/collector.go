// initialize/collector.go
// 创建用户和授权：
// create user ro_query with login password '9enlmubivmpj_9cl';
// grant pg_monitor to ro_query;

package initialize

import (
	"context"
	"fmt"
	"github.com/xiaohongshu/PnSql/server/global"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/dao/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SimpleCollector 结构体
type SimpleCollector struct {
	interval   time.Duration
	stopChan   chan struct{}
	maxWorkers int
	running    bool
	mutex      sync.Mutex
}

// InitSimpleCollector 初始化采集器
func InitSimpleCollector(interval time.Duration) *SimpleCollector {
	return &SimpleCollector{
		interval:   interval,
		stopChan:   make(chan struct{}),
		maxWorkers: 3, // 默认3个并发
		running:    false,
	}
}

// Start 启动采集器
func (c *SimpleCollector) Start() {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return
	}
	c.running = true
	c.mutex.Unlock()

	go func() {
		for {
			select {
			case <-c.stopChan:
				c.mutex.Lock()
				c.running = false
				c.mutex.Unlock()
				log.Println("采集器已停止")
				return
			default:
				c.collectPostgresStats()
				time.Sleep(c.interval)
			}
		}
	}()
	log.Printf("PostgreSQL监控采集器已启动，采集间隔: %v", c.interval)
}

// Stop 停止采集器
func (c *SimpleCollector) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.running {
		close(c.stopChan)
		c.running = false
	}
}

// IsRunning 检查采集器是否在运行
func (c *SimpleCollector) IsRunning() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.running
}

// collectPostgresStats 并发采集统计
func (c *SimpleCollector) collectPostgresStats() {
	instances, err := service.GroupApp.ExampleServer.MysqlInstanceService.GetAllIPAndPorts()
	if err != nil {
		log.Printf("获取实例列表失败: %v", err)
		return
	}

	if len(instances) == 0 {
		log.Println("没有可用的PostgreSQL实例")
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, c.maxWorkers)

	for _, instance := range instances {
		wg.Add(1)
		sem <- struct{}{}

		go func(ip string, port int) {
			defer wg.Done()
			defer func() { <-sem }()
			c.monitorPostgresInstance(ip, port)
		}(instance.IP, instance.Port)
	}

	wg.Wait()
}

// 定义查询结果结构体
type sessionRecord struct {
	PID             int        `gorm:"column:pid"`
	Datname         string     `gorm:"column:datname"`
	ApplicationName string     `gorm:"column:application_name"`
	XactStart       *time.Time `gorm:"column:xact_start"`
	OpenTiming      string     `gorm:"column:open_timing"`
	SQLTiming       string     `gorm:"column:sql_timing"`
	Query           string     `gorm:"column:query"`
	State           string     `gorm:"column:state"`
	WaitEvent       string     `gorm:"column:wait_event"`
}

// monitorPostgresInstance 监控单个实例
func (c *SimpleCollector) monitorPostgresInstance(ip string, port int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s "+
			"sslmode=disable application_name=PG_Monitor timezone=Asia/Shanghai",
		ip, port, global.P_cfg.Postgresql.Username, global.P_cfg.Postgresql.Password, global.P_cfg.Postgresql.Dbname,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("CST", 8*60*60))
		},
	})
	if err != nil {
		log.Printf("[%s:%d] 连接失败: %v", ip, port, err)
		return
	}

	defer func() {
		if rawDB, err := db.DB(); err == nil && rawDB != nil {
			rawDB.Close()
		}
	}()

	var activeSessions []sessionRecord

	if err := db.WithContext(ctx).Raw(`
    SELECT 
        pid,
        datname,
        application_name,
        CASE 
            WHEN xact_start IS NULL THEN NULL 
            ELSE xact_start AT TIME ZONE 'Asia/Shanghai' 
        END AS xact_start,        
        (now() - xact_start) AS open_timing,
        (now() - query_start) AS sql_timing,
        query,
        state,
        COALESCE(wait_event, 'NULL') AS wait_event
    FROM pg_stat_activity
    WHERE state != 'idle'
    AND pid != pg_backend_pid()
    AND now() - query_start IS NOT NULL
`).Scan(&activeSessions).Error; err != nil {
		log.Printf("[%s:%d] 查询失败: %v", ip, port, err)
		return
	}

	// 转换为PgSession并保存
	var pgSessions []example.PgSession
	for _, s := range activeSessions {
		// 转换时区到CST
		var xactStartCST *time.Time
		if s.XactStart != nil {
			t := s.XactStart.In(time.FixedZone("CST", 8*60*60))
			xactStartCST = &t
		}

		pgSessions = append(pgSessions, example.PgSession{
			Source:          ip,
			PID:             int32(s.PID),
			Datname:         s.Datname,
			ApplicationName: s.ApplicationName,
			XactStart:       xactStartCST,
			OpenTiming:      s.OpenTiming,
			SQLTiming:       s.SQLTiming,
			Query:           s.Query,
			State:           s.State,
			WaitEvent:       s.WaitEvent,
		})
	}

	if err := service.GroupApp.ExampleServer.PgSessionService.BatchInsertSessions(pgSessions); err != nil {
		log.Printf("[%s:%d] 保存会话数据失败: %v", ip, port, err)
	}

	// 输出日志
	c.logSessions(ip, port, activeSessions)
}

// logSessions 输出会话日志
func (c *SimpleCollector) logSessions(ip string, port int, sessions []sessionRecord) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n[%s:%d] 活跃查询（共%d条）:\n", ip, port, len(sessions)))

	for _, s := range sessions {
		xactStartStr := "NULL"
		if s.XactStart != nil {
			xactStartStr = s.XactStart.In(time.FixedZone("CST", 8*60*60)).Format("2006-01-02 15:04:05")
		}

		query := s.Query
		if len(query) > 300 {
			query = query[:300] + "..."
		}

		sb.WriteString(fmt.Sprintf("  PID:%-6d DB:%-12s App:%-20s XactStart:%-20s OpenTime:%-15s SQLTime:%-15s State:%-10s WaitEvent:%-15s | %s\n",
			s.PID,
			s.Datname,
			s.ApplicationName,
			xactStartStr,
			s.OpenTiming,
			s.SQLTiming,
			s.State,
			s.WaitEvent,
			query,
		))
	}

	if len(sessions) == 0 {
		sb.WriteString(fmt.Sprintf("[%s:%d] 无活跃查询\n", ip, port))
	}

	log.Print(sb.String())
}
