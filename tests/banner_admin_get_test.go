package tests

import (
	"net/http"
	"testing"
)

func TestBannerAdminGet_AsUser_Forbidden(t *testing.T) {
	e, tokenUser, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.GET("/banner", id).
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).
		WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenUser).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerAdminGet_BannerNotActive_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	b.IsActive = false
	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.GET("/banner", id).
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).
		WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusOK)
}

func TestBannerAdminGet_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.GET("/banner", id).
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).
		WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusOK)
}
