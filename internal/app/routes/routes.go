package routes

import (
	"log/slog"
	"net/http"

	"avito-test-task/internal/app/routes/middleware"
	adm "avito-test-task/internal/handlers/admin/banner"
	bannerhndl "avito-test-task/internal/handlers/banner"
	"avito-test-task/internal/lib/jwt"
	bannersvc "avito-test-task/internal/service/banner"
)

// New creates a new router with all the middlewares.
func New(logger *slog.Logger, manager *jwt.Manager, bannerSvc *bannersvc.Service) http.Handler {
	healthRouter := http.NewServeMux()
	healthRouter.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	router := http.NewServeMux()
	router.Handle("GET /user_banner", bannerhndl.NewGetHandler(bannerSvc, logger))

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
	admRouter.Handle("DELETE /banner", adm.NewDeleteByFeatureTagHandler(bannerSvc, logger))

	router.Handle("/", middleware.EnsureAdmin(admRouter, logger))

	mainRouter := http.NewServeMux()
	mainRouter.Handle("GET /health", healthRouter)
	mainRouter.Handle("/", mw(router))

	return mainRouter
}
