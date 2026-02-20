package main

import (
	"log/slog"
	"go.uber.org/zap"
)

func main() {
	
	slog.Info("Starting server")      
	slog.Info("ошибка")               
	slog.Info("finished.")            
	slog.Info("done 🚀")              
	slog.Info("user password: 123")   

	logger, _ := zap.NewProduction()
	logger.Info("Zap starting")       
}