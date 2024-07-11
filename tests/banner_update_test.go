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
		WithJSON(getCreateBannerDTO()).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	e.PATCH("/banner/{id}", id).
		WithJSON(getUpdateBannerDTO()).
		WithHeader("Authorization", "Bearer "+tokenUser).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerUpdate_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()

	v := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id := int64(v.Raw().(float64))

	updDTO := getUpdateBannerDTO()

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

	assert := assert.New(t)
	assert.Equal(*updDTO.Content.Title, upd["title"])
	assert.Equal(*updDTO.Content.Text, upd["text"])
	assert.Equal(*updDTO.Content.URL, upd["url"])
}
