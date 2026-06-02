package constant

type ContextKey string

const (
	UserIDKey        ContextKey = "user_id"
	UsernameKey      ContextKey = "username_token"
	NewTokenKey      ContextKey = "new_token"
	LoggerKey        ContextKey = "logger"
	DBTransactionKey ContextKey = "tx"
	IPAddressKey     ContextKey = "ip_address"
)
