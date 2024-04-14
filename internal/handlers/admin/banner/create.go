package banner

import (
	api2 "avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/models/dto/banner"
	"avito-test-task/internal/service"
	"errors"
	"log/slog"
	"net/http"
)

type CreateResponse struct {
	BannerID int64 `json:"banner_id,omitempty"`
	api2.Response
}

func NewCreateHandler(svc *service.Banner, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.banner.create"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api2.RequestIDKey, api2.RequestID(r)),
		)

		req := new(banner.CreateDTO)
		err := jsn.DecodeRequest(r, req, log)
		if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(err.Error()), log)
			return
		}

		id, err := svc.SaveBanner(r.Context(), *req)
		if errors.Is(err, service.ErrBannerAlreadyExists) {
			jsn.EncodeResponse(w, http.StatusConflict, api2.ErrResponse(err.Error()), log)
			return
		} else if validErr := new(service.ValidationError); errors.As(err, validErr) {
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(validErr.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusInternalServerError, api2.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusCreated, CreateResponse{id, api2.OkResponse()}, log)
	}
}
