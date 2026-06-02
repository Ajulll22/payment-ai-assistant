package database


type SQLConfig struct {
	User        string
	PasswordEnc string
	Host        string
	Port        string
	Name        string
	Timeout     int

	LogDir          string
	LogName         string
	MaxLogDays      int
	FallbackLogFile *os.File
}

func SQLConnect(cfg SQLConfig) (db *gorm.DB, err error) {
	clear_password := security.Decrypt(cfg.PasswordEnc, "62277ecdae08d9e813ab17a4ec2db8c58db38e398617824a2ef035c64d3da4be")

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
	gormLogger := logger.GetGormLogger(cfg.LogDir, cfg.LogName, cfg.MaxLogDays, gormLog.Info, cfg.FallbackLogFile)

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
