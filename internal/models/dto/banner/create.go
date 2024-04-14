package banner

import "avito-test-task/internal/models/entity"

type CreateDTO struct {
	TagIDs    []int64       `json:"tag_ids" validate:"required,gt=0,dive"`
	FeatureID int64         `json:"feature_id" validate:"required"`
	Content   CreateContent `json:"content" validate:"required"`
	IsActive  bool          `json:"is_active"`
}

type CreateContent struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}

func (d CreateDTO) ToModel() *entity.Banner {
	return entity.NewBanner(
		d.Content.Title,
		d.Content.Text,
		d.Content.URL,
		d.FeatureID,
		d.IsActive,
		d.TagIDs,
	)
}
