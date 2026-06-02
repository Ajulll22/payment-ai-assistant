package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ajulll22/payment-ai-assistant/pkg/logger"
	"github.com/Ajulll22/payment-ai-assistant/pkg/security"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

type SQLConfig struct {
	AppKey string

	User        string
	PasswordEnc string
	Host        string
	Port        string
	Name        string
	Timeout     int

	LogDir          string
	LogMaxFile      int
	FallbackLogFile *os.File
}

func SQLConnect(cfg SQLConfig) (db *gorm.DB, err error) {
	clear_password := security.Decrypt(cfg.PasswordEnc, cfg.AppKey)

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, clear_password, cfg.Host, cfg.Port, cfg.Name)

	timeout := time.Duration(cfg.Timeout) * time.Second
	return waitForSQLServer(dsn, timeout, cfg)
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func waitForSQLServer(dsn string, timeout time.Duration, cfg SQLConfig) (db *gorm.DB, err error) {
	start := time.Now()
	gormLogger := logger.GetGormLogger(cfg.LogDir, cfg.Name, cfg.LogMaxFile, gormLog.Info, cfg.FallbackLogFile)

	for {
		// Coba koneksi ke SQL Server
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger:      gormLogger,
			PrepareStmt: true,
		})
		if err == nil {
			sqlDB, _ := db.DB()

			// Coba Ping
			if err := sqlDB.Ping(); err == nil {
				log.Println("SQL Server is ready!")
				return db, nil
			}
		}

		// Jika waktu tunggu habis, kembalikan error
		if time.Since(start) > timeout {
			return db, fmt.Errorf("timeout waiting for SQL Server to be ready")
		}

		// Tunggu beberapa saat sebelum mencoba lagi
		log.Println("Waiting for SQL Server...")
		time.Sleep(5 * time.Second)
	}
}
