package a

import (
	"log/slog"
	"go.uber.org/zap"
)


func example() {
	slog.Info("Server started") // want "log message should start with a lowercase letter.*"
	zap.Info("Connection lost") // want "log message should start with a lowercase letter.*"
	
	slog.Info("server started")

	slog.Info("запуск сервера") // want "log message should contain only English words.*"
	slog.Info("error: ошибка")  // want "log message should contain only English words.*"

	slog.Info("done.")          // want "log message should not end with punctuation.*"
	slog.Info("failed!")        // want "log message should not end with punctuation.*"
	slog.Info("rocket 🚀 launch") // want "log message should not contain emojis.*"

	slog.Info("user password: 123") // want "log message contains potential sensitive data.*"
	slog.Info("api_key is set")     // want "log message contains potential sensitive data.*"
	slog.Info("testforcustomkeyword is test") // want "log message contains potential sensitive data.*"
}