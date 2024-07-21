package banner

import "avito-test-task/internal/model/entity"

// CreateDTO is expected to be received as a create banner request.
type CreateDTO struct {
	TagIDs    []int64       `json:"tag_ids" validate:"required,gt=0,dive"`
	FeatureID int64         `json:"feature_id" validate:"required"`
	Content   CreateContent `json:"content" validate:"required"`
	IsActive  bool          `json:"is_active"`
}

// CreateContent contains information about banner that's being created.
type CreateContent struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}

// ToModel returns a new entity.Banner constructed from CreateDTO.
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
