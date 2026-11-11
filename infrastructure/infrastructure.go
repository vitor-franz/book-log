package infrastructure

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitor-franz/book-log/infrastructure/router/middleware"
	settings "github.com/vitor-franz/book-log/infrastructure/settings/impl"
	"github.com/vitor-franz/book-log/infrastructure/setup"
	"github.com/vitor-franz/book-log/modules/auth"
	"github.com/vitor-franz/book-log/utils/logger"
)

func Start() error {
	settingsPath := flag.String("settings", "settings.toml", "path to the TOML settings file")
	flag.Parse()

	settingsLoader, err := settings.NewTOMLSettingsLoader(*settingsPath)
	if err != nil {
		return err
	}

	runtimeSettings := settingsLoader.GetRuntimeSettings()
	loggerSettings := settingsLoader.GetLoggerSettings()

	handler, err := logger.CreateLogHandler(loggerSettings.OutputPath, runtimeSettings.Environment)
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(handler))

	jwtSettings := settingsLoader.GetJWTSettings()

	db, err := setup.OpenDatabase(settingsLoader)
	if err != nil {
		return err
	}
	defer db.Close()

	router := mux.NewRouter()
	router.Use(middleware.NewRootMiddleware())

	apiRouter := router.PathPrefix("/api").Subrouter()

	authRepository := auth.NewRepository(db)

	authUseCase := auth.NewUseCase(authRepository, jwtSettings.Secret)

	auth.NewModule(authUseCase, jwtSettings.InsecureCookie).Setup(router, apiRouter)

	slog.Info("Starting HTTP server", slog.String("address", ":8080"))

	return http.ListenAndServe(":8080", router)
}
