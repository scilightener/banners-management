package banner

import (
	api2 "avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/models/dto/banner"
	"avito-test-task/internal/service"
	"errors"
	"log/slog"
	"net/http"
)

func NewUpdateHandler(svc *service.Banner, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.banner.update"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api2.RequestIDKey, api2.RequestID(r)),
		)

		var id int64
		err := api2.ParseInt64(r.PathValue("id"), "id", &id)
		if err != nil {
			log.Error("failed to parse id", sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(err.Error()), log)
			return
		}
		req := new(banner.UpdateDTO)
		err = jsn.DecodeRequest(r, req, log)
		if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(err.Error()), log)
			return
		}

		err = svc.UpdateBanner(r.Context(), id, *req)
		if validErr := new(service.ValidationError); errors.As(err, validErr) {
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(validErr.Error()), log)
			return
		} else if errors.Is(err, service.ErrBannerNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api2.ErrResponse(err.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusInternalServerError, api2.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusOK, api2.OkResponse(), log)
	}
}
