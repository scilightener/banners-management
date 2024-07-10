package banner

import (
	"avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/logger/sl"
	bannerdto "avito-test-task/internal/models/dto/banner"
	"avito-test-task/internal/service"
	bannersvc "avito-test-task/internal/service/banner"
	"errors"
	"log/slog"
	"net/http"
)

func NewUpdateHandler(svc *bannersvc.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.admin.banner.update"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		var id int64
		err := api.ParseInt64(r.PathValue("id"), "id", &id)
		if err != nil {
			log.Error("failed to parse id", sl.Err(err))
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}
		req := new(bannerdto.UpdateDTO)
		err = jsn.DecodeRequest(r, req, log)
		if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}

		err = svc.UpdateBanner(r.Context(), id, *req)
		if validErr := new(service.ValidationError); errors.As(err, validErr) {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(validErr.Error()), log)
			return
		} else if errors.Is(err, bannersvc.ErrNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api.ErrResponse(err.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusInternalServerError, api.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusOK, api.OkResponse(), log)
	}
}
