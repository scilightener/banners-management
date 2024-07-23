package auth

import (
	"banners-management/internal/lib/api"
	"banners-management/internal/lib/api/jsn"
	"banners-management/internal/lib/jwt"
	"banners-management/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type Response struct {
	Token string `json:"token"`
}

func NewAuthHandler(j *jwt.Manager, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.auth.auth"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		p := r.URL.Query()
		role := p.Get("role")
		if role == "" {
			log.Info("role param not specified")
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse("role param not specified"), log)
			return
		}

		token, err := j.GenerateToken(role)
		if err != nil {
			log.Info("failed to generate token", sl.Err(err))
			jsn.EncodeResponse(w, http.StatusInternalServerError, api.ErrResponse("failed to generate token"), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusOK, Response{Token: token}, log)
	}
}
