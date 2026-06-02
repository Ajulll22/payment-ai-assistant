package logger

type LogType string

const (
	LogTypeApp      LogType = "app_log"
	LogTypeEmail    LogType = "email_log"
	LogTypeFirebase LogType = "firebase_log"
	LogTypeSMS      LogType = "sms_log"
	LogTypeWhatsapp LogType = "whatsapp_log"
	LogTypeOAuth    LogType = "oauth_log"
	LogTypeMQTT     LogType = "mqtt_log"
	LogTypeWorker   LogType = "worker_log"
)
