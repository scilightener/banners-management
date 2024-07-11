package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBannerUpdate_AsUser_Fail(t *testing.T) {
	e, tokenUser, tokenAdm := initTest(t)

	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(newCreateBannerDTO()).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := rawToInt64(v.Raw())

	e.PATCH("/banner/{id}", id).
		WithJSON(newUpdateBannerDTO()).
		WithHeader("Authorization", "Bearer "+tokenUser).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerUpdate_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)
	b := newCreateBannerDTO()
	updDTO := newUpdateBannerDTO()

	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := rawToInt64(v.Raw())

	e.PATCH("/banner/{id}", id).
		WithJSON(updDTO).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusOK)

	upd := e.GET("/user_banner").
		WithQuery("feature_id", *updDTO.FeatureID).WithQuery("tag_id", (*updDTO.TagIDs)[0]).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Raw()

	asrt := assert.New(t)
	asrt.Equal(*updDTO.Content.Title, upd["title"])
	asrt.Equal(*updDTO.Content.Text, upd["text"])
	asrt.Equal(*updDTO.Content.URL, upd["url"])
}

func TestBannerUpdate_BannerConflict(t *testing.T) {
	e, _, tokenAdm := initTest(t)
	b1 := newCreateBannerDTO()
	b2 := newCreateBannerDTO()
	updDTO := updateBannerDTO(&b1.FeatureID, &b1.TagIDs, nil)

	e.POST("/banner").
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
	id2 := rawToInt64(v2.Raw())

	e.PATCH("/banner/{id}", id2).
		WithJSON(updDTO).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().ContainsKey("error").
		Value("error").String().Length().Gt(0)
}
