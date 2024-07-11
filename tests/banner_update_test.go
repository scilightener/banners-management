package tests

import (
	"avito-test-task/internal/models/dto/banner"
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

// TODO: add conflict by feature & tags separately !!! TRIGGER DOESN'T WORK ON UPDATE !!!
func TestBannerUpdate_NewBannerConflict(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b1 := getCreateBannerDTO()
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b1).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	b2 := getCreateBannerDTO()
	v2 := e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b2).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		JSON().Object().Value("banner_id")
	id2 := int64(v2.Raw().(float64))

	updDTO := banner.UpdateDTO{
		TagIDs:    &b1.TagIDs,
		FeatureID: &b1.FeatureID,
	}

	e.PATCH("/banner/{id}", id2).
		WithJSON(updDTO).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().ContainsKey("error").
		Value("error").String().Length().Gt(0)
}
