package tests

import (
	"net/http"
	"testing"
)

func TestBannerDelete_AsUser_Fail(t *testing.T) {
	e, tokenUser, tokenAdm := initTest(t)

	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(getCreateBannerDTO()).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.DELETE("/banner/{id}", id).
		WithMaxRetries(5).
		WithHeader("Authorization", "Bearer "+tokenUser).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerDelete_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()

	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.DELETE("/banner/{id}", id).
		WithMaxRetries(5).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusNoContent)

	e.GET("/user_banner").
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).
		WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusNotFound)
}

func TestBannerDelete_NotFound(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	e.DELETE("/banner/{id}", 100000000).
		WithMaxRetries(5).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusNotFound)
}
