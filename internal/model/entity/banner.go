package entity

import "time"

// Banner is a banner domain entity.
type Banner struct {
	ID        int64
	Title     string
	Text      string
	URL       string
	FeatureID int64
	IsActive  bool
	TagIDs    []int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBanner returns a new Banner instance.
func NewBanner(
	title string,
	text string,
	url string,
	featureID int64,
	isActive bool,
	tagIDs []int64,
) *Banner {
	now := time.Now()
	return &Banner{
		Title:     title,
		Text:      text,
		URL:       url,
		FeatureID: featureID,
		IsActive:  isActive,
		TagIDs:    tagIDs,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdatableBanner is a banner domain entity, that's being used to update a main Banner entity.
// Pointer parameters indicate that they're optional, and are not considered during update.
type UpdatableBanner struct {
	ID        int64
	Title     *string
	Text      *string
	URL       *string
	FeatureID *int64
	IsActive  *bool
	TagIDs    *[]int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
