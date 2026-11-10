package main

import (
	"log/slog"

	"github.com/vitor-franz/book-log/modules/auth"
)

func main() {
	err := start()
	if err != nil {
		panic(err)
	}
}

func start() error {
	authRepository := auth.NewRepository()

	authUseCase := auth.NewUseCase(authRepository)

	auth.NewModule(authUseCase)

	slog.Info("started application")

	return nil
}
