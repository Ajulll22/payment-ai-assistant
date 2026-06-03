package constant

import (
	"os"
	"strconv"
)

type Config struct {
	App               AppConfig
	AssistantSQL      AssistantConfig
	AssistantAnalytic AssistantConfig
	DBApp             DBConfig
	DBDWH             DBConfig
	Log               LogConfig
}

func GetEnv() *Config {
	dbAppTimeout, _ := strconv.Atoi(os.Getenv("DB_APP_TIMEOUT"))
	dbAppLogMaxFile, _ := strconv.Atoi(os.Getenv("DB_APP_LOG_MAX_FILE"))
	dbDWHTimeout, _ := strconv.Atoi(os.Getenv("DB_DWH_TIMEOUT"))
	dbDWHLogMaxFile, _ := strconv.Atoi(os.Getenv("DB_DWH_LOG_MAX_FILE"))

	logRequestStatus, _ := strconv.ParseBool(os.Getenv("LOG_REQUEST_STATUS"))
	logDebugStatus, _ := strconv.ParseBool(os.Getenv("LOG_DEBUG_STATUS"))
	logMaxFile, _ := strconv.Atoi(os.Getenv("LOG_MAX_FILE"))

	return &Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: os.Getenv("APP_PORT"),
			Key:  os.Getenv("APP_KEY"),
		},
		AssistantSQL: AssistantConfig{
			Model: os.Getenv("ASSISTANT_SQL_MODEL"),
			URL:   os.Getenv("ASSISTANT_SQL_URL"),
		},
		AssistantAnalytic: AssistantConfig{
			Model: os.Getenv("ASSISTANT_ANALYTIC_MODEL"),
			URL:   os.Getenv("ASSISTANT_ANALYTIC_URL"),
		},
		DBApp: DBConfig{
			User:       os.Getenv("DB_APP_USER"),
			Password:   os.Getenv("DB_APP_PASSWORD"),
			Name:       os.Getenv("DB_APP_NAME"),
			Host:       os.Getenv("DB_APP_HOST"),
			Port:       os.Getenv("DB_APP_PORT"),
			Timeout:    dbAppTimeout,
			LogMaxFile: dbAppLogMaxFile,
		},
		DBDWH: DBConfig{
			User:       os.Getenv("DB_DWH_USER"),
			Password:   os.Getenv("DB_DWH_PASSWORD"),
			Name:       os.Getenv("DB_DWH_NAME"),
			Host:       os.Getenv("DB_DWH_HOST"),
			Port:       os.Getenv("DB_DWH_PORT"),
			Timeout:    dbDWHTimeout,
			LogMaxFile: dbDWHLogMaxFile,
		},
		Log: LogConfig{
			Path:          os.Getenv("LOG_PATH"),
			DebugFilename: os.Getenv("LOG_DEBUG_FILENAME"),
			ErrorFilename: os.Getenv("LOG_ERROR_FILENAME"),
			RequestStatus: logRequestStatus,
			DebugStatus:   logDebugStatus,
			MaxFile:       logMaxFile,
		},
	}
}

type AppConfig struct {
	Name string
	Port string
	Key  string
}

type AssistantConfig struct {
	Model string
	URL   string
}

type DBConfig struct {
	User       string
	Password   string
	Name       string
	Host       string
	Port       string
	Timeout    int
	LogMaxFile int
}

type LogConfig struct {
	Path          string
	RequestStatus bool
	DebugStatus   bool
	MaxFile       int
	DebugFilename string
	ErrorFilename string
}
