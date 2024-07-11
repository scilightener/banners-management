package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBannerUserGet_NotAuthed_Unauthorized(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect()

	e.GET("/user_banner").
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).WithQuery("tag_id", b.TagIDs[0]).
		Expect().
		Status(http.StatusUnauthorized)
}

func TestBannerUserGet_BannerNotActive(t *testing.T) {
	e, tokenUsr, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	b.IsActive = false
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect()

	e.GET("/user_banner").
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenUsr).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerUserGet_Successful(t *testing.T) {
	e, tokenUsr, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect()

	resp := e.GET("/user_banner").
		WithMaxRetries(5).
		WithQuery("feature_id", b.FeatureID).WithQuery("tag_id", b.TagIDs[0]).
		WithHeader("Authorization", "Bearer "+tokenUsr).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Raw()
	assert := assert.New(t)
	_, ok := resp["title"].(string)
	assert.True(ok)
	_, ok = resp["text"].(string)
	assert.True(ok)
	_, ok = resp["url"].(string)
	assert.True(ok)
}

func TestBannerUserGet_NotFound(t *testing.T) {
	e, tokenUsr, _ := initTest(t)

	e.GET("/user_banner").
		WithMaxRetries(5).
		WithQuery("feature_id", 0).WithQuery("tag_id", 0).
		WithHeader("Authorization", "Bearer "+tokenUsr).
		Expect().
		Status(http.StatusNotFound)
}

func TestBannerUserGet_InvalidData_FailCases(t *testing.T) {
	e, tokenUsr, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect()

	testCases := []struct {
		name           string
		featureID      interface{}
		tagID          interface{}
		expectedStatus int
	}{
		{
			name:           "InvalidFeatureID",
			featureID:      "invalid",
			tagID:          b.TagIDs[0],
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidTagID",
			featureID:      b.FeatureID,
			tagID:          "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			e.GET("/user_banner").
				WithMaxRetries(5).
				WithQuery("feature_id", tc.featureID).WithQuery("tag_id", tc.tagID).
				WithHeader("Authorization", "Bearer "+tokenUsr).
				Expect().
				Status(tc.expectedStatus)
		})
	}
}
