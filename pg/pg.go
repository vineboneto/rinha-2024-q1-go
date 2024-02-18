package pg

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbContext *gorm.DB
	once      sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		initDB()
	})
	return dbContext
}

func initDB() {
	host := os.Getenv("DB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}
	dsn := fmt.Sprintf("user=postgres password=1234 dbname=rinha host=%s sslmode=disable", host)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Configurar as opções do pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	totalConn := 500

	// Configurar o número máximo de conexões no pool
	sqlDB.SetMaxOpenConns(totalConn)

	// Configurar o número máximo de conexões inativas no pool
	sqlDB.SetMaxIdleConns(totalConn)

	// Configurar o tempo de vida máximo das conexões no pool
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	for i := 0; i < totalConn; i++ {
		_, err := sqlDB.Exec("SELECT 1")
		if err != nil {
			log.Fatal(err)
		}
	}

	dbContext = db

}

func CloseDB() {
	sqlDB, err := dbContext.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}
