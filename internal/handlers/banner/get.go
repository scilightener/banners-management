package banner

import (
	"errors"
	"log/slog"
	"net/http"

	"avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/service/banner"
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
	api.Response
}

func NewGetHandler(svc *banner.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.banner.get"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		p := r.URL.Query()
		var (
			fID, tID int64
			uLR      bool
			resErr   error
		)
		if err := api.ParseInt64(p.Get(featureID), featureID, &fID); err != nil {
			resErr = errors.Join(resErr, err)
		}
		if err := api.ParseInt64(p.Get(tagID), tagID, &tID); err != nil {
			resErr = errors.Join(resErr, err)
		}
		if err := api.ParseBool(p.Get(useLastRevision), useLastRevision, &uLR); err != nil {
			uLR = false // no error, parameter is optional. default is false
		}

		if resErr != nil {
			log.Error("failed to parse query params", sl.Err(resErr))
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(resErr.Error()), log)
			return
		}

		b, err := svc.BannerByFeatureTag(r.Context(), fID, tID, uLR, true)
		if errors.Is(err, banner.ErrNotActive) {
			jsn.EncodeResponse(w, http.StatusForbidden, api.ErrResponse(err.Error()), log)
			return
		} else if errors.Is(err, banner.ErrNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api.ErrResponse(err.Error()), log)
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}

		jsn.EncodeResponse(w, http.StatusOK, GetResponse{
			b.Title,
			b.Text,
			b.URL,
			api.OkResponse(),
		}, log)
	}
}
