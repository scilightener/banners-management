package banner

import (
	"banners-management/internal/lib/api"
	"banners-management/internal/lib/api/jsn"
	bannerdto "banners-management/internal/model/dto/banner"
	"banners-management/internal/service"
	bannersvc "banners-management/internal/service/banner"
	"errors"
	"log/slog"
	"net/http"
)

type CreateResponse struct {
	BannerID int64 `json:"banner_id,omitempty"`
	api.Response
}

func NewCreateHandler(svc *bannersvc.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.admin.banner.create"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		req := new(bannerdto.CreateDTO)
		err := jsn.DecodeRequest(r, req, log)
		if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}

		id, err := svc.SaveBanner(r.Context(), *req)
		if errors.Is(err, bannersvc.ErrAlreadyExists) {
			jsn.EncodeResponse(w, http.StatusConflict, api.ErrResponse(err.Error()), log)
			return
		} else if validErr := new(service.ValidationError); errors.As(err, validErr) {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(validErr.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusInternalServerError, api.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusCreated, CreateResponse{id, api.OkResponse()}, log)
	}
}
