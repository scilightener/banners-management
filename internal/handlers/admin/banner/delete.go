package banner

import (
	api2 "avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/service"
	"errors"
	"log/slog"
	"net/http"
)

func NewDeleteHandler(svc *service.Banner, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.banner.delete"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api2.RequestIDKey, api2.RequestID(r)),
		)

		var id int64
		err := api2.ParseInt64(r.PathValue("id"), "id", &id)
		if err != nil {
			log.Error("failed to parse query params", sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(err.Error()), log)
			return
		}

		err = svc.DeleteBanner(r.Context(), id)
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

		jsn.EncodeResponse(w, http.StatusNoContent, api2.OkResponse(), log)
	}
}
