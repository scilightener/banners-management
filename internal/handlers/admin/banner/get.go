package banner

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/model/entity"
	"avito-test-task/internal/service/banner"
)

const (
	featureID = "feature_id"
	tagID     = "tag_id"
	limit     = "limit"
	offset    = "offset"
)

type GetResponse []GetResponseItem

type GetResponseItem struct {
	BannerID  int64   `json:"banner_id"`
	TagIDs    []int64 `json:"tag_ids"`
	FeatureID int64   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		URL   string `json:"url"`
	} `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ri *GetResponseItem) fromEntity(b *entity.Banner) {
	ri.BannerID = b.ID
	ri.TagIDs = b.TagIDs
	ri.FeatureID = b.FeatureID
	ri.Content.Title = b.Title
	ri.Content.Text = b.Text
	ri.Content.URL = b.URL
	ri.IsActive = b.IsActive
	ri.CreatedAt = b.CreatedAt
	ri.UpdatedAt = b.UpdatedAt
}

func NewGetHandler(svc *banner.Service, log *slog.Logger) http.HandlerFunc {
	const comp = "handlers.admin.banner.get"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("comp", comp),
			slog.String(api.RequestIDKey, api.RequestID(r)),
		)

		p := r.URL.Query()
		var (
			fID, tID = new(int64), new(int64)
			li, off  = new(int), new(int)
			uLR      = true
		)
		err := api.ParseInt64(p.Get(featureID), featureID, fID)
		if err != nil {
			fID = nil
		}
		err = api.ParseInt64(p.Get(tagID), tagID, tID)
		if err != nil {
			tID = nil
		}
		err = api.ParseInt(p.Get(limit), limit, li)
		if err != nil {
			li = nil
		}
		err = api.ParseInt(p.Get(offset), offset, off)
		if err != nil {
			off = nil
		}

		bs, err := svc.BannersByFeatureTag(r.Context(), fID, tID, li, off, &uLR)
		if errors.Is(err, banner.ErrNotFound) {
			jsn.EncodeResponse(w, http.StatusNotFound, api.ErrResponse(err.Error()), log)
			return
		} else if err != nil {
			jsn.EncodeResponse(w, http.StatusBadRequest, api.ErrResponse(err.Error()), log)
			return
		}

		resp := make([]GetResponseItem, len(bs))
		for i, b := range bs {
			var ri GetResponseItem
			ri.fromEntity(b)
			resp[i] = ri
		}
		jsn.EncodeResponse(w, http.StatusOK, GetResponse(resp), log)
	}
}
