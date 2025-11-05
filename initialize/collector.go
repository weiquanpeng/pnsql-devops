// initialize/collector.go
// 创建用户和授权：
// create user ro_query with password '9enlmubivmpj_9cl';
// grant pg_monitor to ro_query;

package initialize

import (
	"context"
	"fmt"
	"github.com/xiaohongshu/PnSql/server/dao/model/example"
	"github.com/xiaohongshu/PnSql/server/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"

	"github.com/xiaohongshu/PnSql/server/dao/service"
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
		maxWorkers: 6, // 默认3个并发
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
	global.Logger.Info("PostgreSQL监控采集器已启动，采集间隔: ", c.interval)
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

func (c *SimpleCollector) collectPostgresStats() {
	instances, err := service.GroupApp.ExampleServer.MysqlInstanceService.GetAllIPAndPorts()
	if err != nil {
		global.Logger.Error("获取实例列表失败: ", err)
		return
	}

	if len(instances) == 0 {
		global.Logger.Info("本次采集完成，无可用实例")
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, c.maxWorkers)

	// 添加成功采集计数器
	successCount := 0
	var counterMutex sync.Mutex

	for _, instance := range instances {
		wg.Add(1)
		sem <- struct{}{}

		go func(ip string, port int, user, pwd string, vmname string) {
			defer wg.Done()
			defer func() { <-sem }()

			if c.monitorPostgresInstance(ip, port, user, pwd, vmname) {
				counterMutex.Lock()
				successCount++
				counterMutex.Unlock()
			} else {
				global.Logger.Warnf("采集失败: %s:%d", ip, port)
			}

		}(instance.IP, instance.Port, instance.Username, instance.Password, instance.VMName)
	}

	wg.Wait()

	// 只在有成功采集时记录日志
	if successCount > 0 {
		global.Logger.Infof("本次采集完成，成功采集 %d 台实例", successCount)
	}
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

func (c *SimpleCollector) monitorPostgresInstance(ip string, port int, user string, password string, vmname string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s "+
			"sslmode=require application_name=PG_Monitor timezone=Asia/Shanghai",
		ip, port, user, password, global.P_cfg.Postgresql.Dbname,
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
		global.Logger.Error("postgres连接失败: ", ip, port, err)
		return false
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
            client_addr as application_name,
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
		global.Logger.Error("采集失败: ", ip, port, err)
		return false
	}

	// 即使没有会话也算采集成功
	if len(activeSessions) == 0 {
		return true
	}

	var pgSessions []example.PgSession
	for _, s := range activeSessions {
		var xactStartCST *time.Time
		if s.XactStart != nil {
			t := s.XactStart.In(time.FixedZone("CST", 8*60*60))
			xactStartCST = &t
		}
		pgSessions = append(pgSessions, example.PgSession{
			Source:          ip,
			VMName:          vmname,
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
		global.Logger.Error("保存会话数据失败: ", ip, port, err)
		return false
	}

	return true
}
