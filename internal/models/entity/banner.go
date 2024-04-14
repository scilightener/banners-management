package entity

import "time"

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
