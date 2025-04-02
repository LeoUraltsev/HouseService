package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/config"
	"github.com/LeoUraltsev/HouseService/internal/gen"
	"github.com/LeoUraltsev/HouseService/internal/handlers"
	"github.com/LeoUraltsev/HouseService/internal/jwt"
	"github.com/LeoUraltsev/HouseService/internal/service"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/LeoUraltsev/HouseService/internal/storage/postgres"
)

func Run(log *slog.Logger, cfg *config.Config) error {

	db, err := postgres.New(context.Background(), cfg.PostgresURLConnection, log)
	if err != nil {
		return err
	}
	defer db.Close()

	j := jwt.New(cfg.JWTDuration, cfg.JWTSecret)

	authService := service.NewAuthService(db, j, log)

	r := gen.HandlerWithOptions(
		handlers.New(nil, db, authService, log),
		gen.ChiServerOptions{
			Middlewares: []gen.MiddlewareFunc{
				middleware.RequestID,
			},
		},
	)

	http.ListenAndServe("localhost:10000", r)
	return nil
}
