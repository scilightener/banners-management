package routes

import (
	"avito-test-task/internal/app/routes/middleware"
	adm "avito-test-task/internal/handlers/admin/banner"
	"avito-test-task/internal/handlers/banner"
	"avito-test-task/internal/lib/jwt"
	"avito-test-task/internal/service"
	"log/slog"
	"net/http"
)

// New creates a new router with all the middlewares.
func New(logger *slog.Logger, manager *jwt.Manager, bannerSvc *service.Banner) http.Handler {
	healthRouter := http.NewServeMux()
	healthRouter.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	router := http.NewServeMux()
	router.Handle("GET /user_banner", banner.NewGetHandler(bannerSvc, logger))

	mw := middleware.Chain(
		middleware.NewRecovererMiddleware(logger),
		middleware.RequestIDMiddleware,
		middleware.NewLoggingMiddleware(logger),
		middleware.ContentTypeJSONMiddleware,
		middleware.NewAuthorizationMiddleware(logger, manager),
	)

	admRouter := http.NewServeMux()
	admRouter.Handle("GET /banner", adm.NewGetHandler(bannerSvc, logger))
	admRouter.Handle("POST /banner", adm.NewCreateHandler(bannerSvc, logger))
	admRouter.Handle("PATCH /banner/{id}", adm.NewUpdateHandler(bannerSvc, logger))
	admRouter.Handle("DELETE /banner/{id}", adm.NewDeleteHandler(bannerSvc, logger))

	router.Handle("/", middleware.EnsureAdmin(admRouter, logger))

	mainRouter := http.NewServeMux()
	mainRouter.Handle("GET /health", healthRouter)
	mainRouter.Handle("/", mw(router))

	return mainRouter
}
