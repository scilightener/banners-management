package banner

import "avito-test-task/internal/model/entity"

// UpdateDTO is expected to be received as an update banner request.
// Pointer parameters are optional.
type UpdateDTO struct {
	TagIDs    *[]int64       `json:"tag_ids"`
	FeatureID *int64         `json:"feature_id"`
	Content   *UpdateContent `json:"content"`
	IsActive  *bool          `json:"is_active"`
}

// UpdateContent contains information about banner that's being updated.
// Pointer parameters are optional.
type UpdateContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	URL   *string `json:"url"`
}

// ToModel returns a new entity.UpdatableBanner constructed from UpdateDTO.
func (d UpdateDTO) ToModel(id int64) *entity.UpdatableBanner {
	var (
		title *string
		text  *string
		url   *string
	)
	if d.Content != nil {
		text = d.Content.Text
		title = d.Content.Title
		url = d.Content.URL
	}
	return &entity.UpdatableBanner{
		ID:        id,
		Title:     title,
		Text:      text,
		URL:       url,
		FeatureID: d.FeatureID,
		IsActive:  d.IsActive,
		TagIDs:    d.TagIDs,
	}
}
