package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBannerAdminGet_AsUser_Forbidden(t *testing.T) {
	e, tokenUser, tokenAdm := initTest(t)

	b := newCreateBannerDTO()
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

	b := newCreateBannerDTO()
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

	b := newCreateBannerDTO()
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

func TestBannerAdminGet_MultipleBanners(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b1 := newCreateBannerDTO()
	b2 := createBannerDTO(b1.FeatureID, getNextTagIDs(2), true)

	v1 := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b1).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	v2 := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b2).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id1 := int64(v1.Raw().(float64))
	id2 := int64(v2.Raw().(float64))

	resp := e.GET("/banner").
		WithMaxRetries(5).
		WithQuery("feature_id", b1.FeatureID).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Array()

	require.Equal(t, int64(2), int64(resp.Length().Raw()))
	r1 := int64(resp.Value(0).Object().Raw()["banner_id"].(float64))
	r2 := int64(resp.Value(1).Object().Raw()["banner_id"].(float64))
	require.True(t, r1 == id1 && r2 == id2)
}
