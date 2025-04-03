package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/config"
	"github.com/LeoUraltsev/HouseService/internal/gen"
	"github.com/LeoUraltsev/HouseService/internal/handlers"
	"github.com/LeoUraltsev/HouseService/internal/jwt"
	mv "github.com/LeoUraltsev/HouseService/internal/middleware"
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
	houseService := service.NewHouseService(db, j, log)
	flatService := service.NewFlatService(db, log)

	authMV := mv.Middleware{
		JWT: j,
	}

	r := gen.HandlerWithOptions(
		handlers.New(houseService, flatService, authService, log),
		gen.ChiServerOptions{
			Middlewares: []gen.MiddlewareFunc{
				middleware.RequestID,
				authMV.AuthMW,
			},
		},
	)

	http.ListenAndServe(cfg.HttpAddress, r)
	return nil
}
