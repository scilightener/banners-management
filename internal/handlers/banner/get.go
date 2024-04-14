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

const (
	featureID       = "feature_id"
	tagID           = "tag_id"
	useLastRevision = "use_last_revision"
)

type GetResponse struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
	api2.Response
}

func NewGetHandler(svc *service.Banner, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.banner.get"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api2.RequestIDKey, api2.RequestID(r)),
		)

		p := r.URL.Query()
		var (
			fID, tID int64
			uLR      bool
			resErr   error
		)
		if err := api2.ParseInt64(p.Get(featureID), featureID, &fID); err != nil {
			resErr = errors.Join(resErr, err)
		}
		if err := api2.ParseInt64(p.Get(tagID), tagID, &tID); err != nil {
			resErr = errors.Join(resErr, err)
		}
		if err := api2.ParseBool(p.Get(useLastRevision), useLastRevision, &uLR); err != nil {
			uLR = false // no error, parameter is optional. default is false
		}

		if resErr != nil {
			log.Error("failed to parse query params", sl.Err(resErr))
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(resErr.Error()), log)
			return
		}

		b, err := svc.BannerByFeatureTag(r.Context(), fID, tID, 1, 0, uLR, true)
		if errors.Is(err, service.ErrBannerNotActive) {
			jsn.EncodeResponse(w, http.StatusForbidden, api2.ErrResponse(err.Error()), log)
			return
		} else if errors.Is(err, service.ErrBannerNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api2.ErrResponse(err.Error()), log)
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api2.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusOK, GetResponse{
			b.Title,
			b.Text,
			b.URL,
			api2.OkResponse(),
		}, log)
	}
}
