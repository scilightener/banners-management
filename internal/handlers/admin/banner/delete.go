package banner

import (
	"banners-management/internal/lib/api"
	"banners-management/internal/lib/api/jsn"
	"banners-management/internal/lib/logger/sl"
	"banners-management/internal/service"
	"banners-management/internal/service/banner"
	"errors"
	"log/slog"
	"net/http"
)

func NewDeleteHandler(svc *banner.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.admin.banner.delete"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		var id int64
		err := api.ParseInt64(r.PathValue("id"), "id", &id)
		if err != nil {
			log.Info("failed to parse query params", sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}

		err = svc.DeleteBanner(r.Context(), id)
		if validErr := new(service.ValidationError); errors.As(err, validErr) {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(validErr.Error()), log)
			return
		} else if errors.Is(err, banner.ErrNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api.ErrResponse(err.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusInternalServerError, api.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusNoContent, api.OkResponse(), log)
	}
}
