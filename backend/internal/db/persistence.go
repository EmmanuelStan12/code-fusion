package db

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type PersistenceManager struct {
	config   configs.DBConfig
	DB       *gorm.DB
	entities []interface{}
}

func Init(config configs.DBConfig) *PersistenceManager {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)
	return &PersistenceManager{
		DB:     db,
		config: config,
	}
}

func (manager *PersistenceManager) IsConnected() bool {
	sqlDB, err := manager.DB.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		return false
	}

	return true
}

func (manager *PersistenceManager) RegisterEntity(entities ...interface{}) {
	manager.entities = append(manager.entities, entities...)
}

func (manager *PersistenceManager) Migrate() {
	err := manager.DB.AutoMigrate(manager.entities...)
	if err != nil {
		panic(err)
	}
}
