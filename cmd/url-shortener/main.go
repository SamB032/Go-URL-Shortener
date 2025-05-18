package main

import (
	"log/slog"
	"os"
	"strings"
	"context"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
	server "github.com/SamB032/Go-URL-Shortener/internal/httpServer"
)

type EnvironmentVariables struct {
	ServerPort       string
	LoggingLevel     string
	PostgresHost     string
	PostgresPort     string
	PostgresPassword string
	PostgresUser     string
	PostgresDBName   string
	TemplatesDir     string
}

func getEnvironmentVariables() *EnvironmentVariables {
	return &EnvironmentVariables{
		ServerPort:       os.Getenv("SERVER_PORT"),
		LoggingLevel:     os.Getenv("LOGGING_LEVEL"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDBName:   os.Getenv("POSTGRES_DB"),
		TemplatesDir:     os.Getenv("TEMPLATES_DIR"),
	}
}

func setupLogger(loggingLevel string) *slog.Logger {
	if len(loggingLevel) == 0 {
		loggingLevel = "INFO"
	}

	// Convert the LOGGING_LEVEL to uppercase to make it case-insensitive
	level := strings.ToUpper(loggingLevel)

	// Set the appropriate log level based on the environment variable
	var logLevel slog.Level
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		// Default to INFO if the LOGGING_LEVEL is invalid
		logLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

func main() {
	environmentVariables := getEnvironmentVariables()

	// Initialise Logger
	logger := setupLogger(environmentVariables.LoggingLevel)

	// Setup tracer
	tp, errTp := initTracer()
	if errTp != nil {
		logger.Error("Failed to initialise tracer", slog.String("error", errTp.Error()))
	}
	defer tp.Shutdown(context.Background())


	// Initialise Database
	database, dbErr := database.ConnectToDatabase(
		environmentVariables.PostgresHost,
		environmentVariables.PostgresPort,
		environmentVariables.PostgresUser,
		environmentVariables.PostgresPassword,
		environmentVariables.PostgresDBName,
		logger,
	)
	if dbErr != nil {
		os.Exit(1)
	}

	// Start server
	server := server.NewServer(
		environmentVariables.ServerPort,
		logger,
		database,
		environmentVariables.TemplatesDir,
	)

	err := server.Start(environmentVariables.ServerPort)
	if err != nil {
		logger.Error("Failed to start HTTP server",
			slog.String("serverPort", environmentVariables.ServerPort),
			slog.String("error", err.Error()),
		)
	}
}
