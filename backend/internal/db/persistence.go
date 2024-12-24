package db

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	ErrorCannotInitConnection = "DB_CANNOT_INIT_DATABASE"
)

type PersistenceManager struct {
	config configs.DBConfig
	DB     *gorm.DB
}

func Init(config configs.DBConfig) *PersistenceManager {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
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
		db:     db,
		config: config,
	}
}

func (manager *PersistenceManager) IsConnected() bool {
	sqlDB, err := manager.db.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		return false
	}

	return true
}
