package db

import (
	"sync"

	"github.com/woxQAQ/gim/internal/models"
	"github.com/woxQAQ/gim/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// Config 数据库配置
type Config struct {
	DatabasePath string // SQLite数据库文件路径
}

// Init 初始化数据库连接并执行迁移
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		// 连接SQLite数据库
		var dl logger.Logger
		dl, err = logger.NewLogger(logger.DomainDatabase, nil, zap.AddCallerSkip(1))
		if err != nil {
			return
		}
		gormLogger := NewGormLogger(dl)
		instance, err = gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			return
		}

		// 自动迁移数据库表
		if err = migrateDB(); err != nil {
			return
		}
	})

	return err
}

// migrateDB 执行数据库迁移
func migrateDB() error {
	// 在这里添加需要迁移的模型
	models := []interface{}{
		&models.User{},
		&models.Message{},
		&models.MessageAttachment{},
		// 在此处添加其他需要迁移的模型
	}

	// 执行迁移
	for _, model := range models {
		if err := instance.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return instance
}

// Close 关闭数据库连接
func Close() error {
	if instance != nil {
		db, err := instance.DB()
		if err != nil {
			return err
		}
		return db.Close()
	}
	return nil
}
