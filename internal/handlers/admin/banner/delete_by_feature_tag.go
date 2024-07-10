package banner

import (
	"avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/service"
	"avito-test-task/internal/service/banner"
	"errors"
	"log/slog"
	"net/http"
)

func NewDeleteByFeatureTagHandler(svc *banner.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.admin.banner.delete_by_feature_tag"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		p := r.URL.Query()
		var fID, tID *int64
		err := api.ParseInt64(p.Get(featureID), featureID, fID)
		if err != nil {
			fID = nil
		}
		err = api.ParseInt64(p.Get(tagID), tagID, tID)
		if err != nil {
			tID = nil
		}

		err = svc.DeleteBannerByFeatureTag(r.Context(), fID, tID)
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
