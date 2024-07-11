package tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	"avito-test-task/internal/models/dto/banner"
)

func TestBannerCreate_AsUser_Fail(t *testing.T) {
	e, tokenUsr, _ := initTest(t)

	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(banner.CreateDTO{
			TagIDs:    getNextTagIDs(2),
			FeatureID: getNextFeatureID(),
			Content: banner.CreateContent{
				Title: gofakeit.Word(),
				Text:  gofakeit.Word(),
				URL:   gofakeit.URL(),
			},
			IsActive: true,
		}).
		WithHeader("Authorization", "Bearer "+tokenUsr).
		Expect().
		Status(http.StatusForbidden)
}

func TestBannerCreate_Successful(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(banner.CreateDTO{
			TagIDs:    getNextTagIDs(2),
			FeatureID: getNextFeatureID(),
			Content: banner.CreateContent{
				Title: gofakeit.Word(),
				Text:  gofakeit.Word(),
				URL:   gofakeit.URL(),
			},
			IsActive: true,
		}).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("banner_id").
		Value("banner_id").Number()
}

// TODO: add conflict by feature & tags separately
func TestBannerCreate_BannerAlreadyExists(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	b := getCreateBannerDTO()
	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("banner_id").
		Value("banner_id").Number()

	e.POST("/banner").
		WithMaxRetries(5).
		WithJSON(b).
		WithHeader("Authorization", "Bearer "+tokenAdm).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().ContainsKey("error").
		Value("error").String().Length().Gt(0)
}

func TestBannerCreate_InvalidData_FailCases(t *testing.T) {
	e, _, tokenAdm := initTest(t)

	testCases := []struct {
		name           string
		dto            banner.CreateDTO
		wrongParamName string
	}{
		{
			name: "Empty TagIDs",
			dto: banner.CreateDTO{
				TagIDs:    []int64{},
				FeatureID: getNextFeatureID(),
				Content: banner.CreateContent{
					Title: gofakeit.Word(),
					Text:  gofakeit.Word(),
					URL:   gofakeit.URL(),
				},
				IsActive: true,
			},
			wrongParamName: "TagIDs",
		},
		{
			name: "Empty FeatureID",
			dto: banner.CreateDTO{
				TagIDs:    getNextTagIDs(2),
				FeatureID: 0,
				Content: banner.CreateContent{
					Title: gofakeit.Word(),
					Text:  gofakeit.Word(),
					URL:   gofakeit.URL(),
				},
				IsActive: true,
			},
			wrongParamName: "FeatureID",
		},
		{
			name: "Empty Content.title",
			dto: banner.CreateDTO{
				TagIDs:    getNextTagIDs(2),
				FeatureID: getNextFeatureID(),
				Content: banner.CreateContent{
					Title: "",
					Text:  gofakeit.Word(),
					URL:   gofakeit.URL(),
				},
				IsActive: true,
			},
			wrongParamName: "Title",
		},
		{
			name: "Empty Content.text",
			dto: banner.CreateDTO{
				TagIDs:    getNextTagIDs(2),
				FeatureID: getNextFeatureID(),
				Content: banner.CreateContent{
					Title: gofakeit.Word(),
					Text:  "",
					URL:   gofakeit.URL(),
				},
				IsActive: true,
			},
			wrongParamName: "Text",
		},
		{
			name: "Invalid Content.url",
			dto: banner.CreateDTO{
				TagIDs:    getNextTagIDs(2),
				FeatureID: getNextFeatureID(),
				Content: banner.CreateContent{
					Title: gofakeit.Word(),
					Text:  gofakeit.Word(),
					URL:   "invalid_url",
				},
				IsActive: true,
			},
			wrongParamName: "URL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			e.POST("/banner").
				WithMaxRetries(5).
				WithJSON(tc.dto).
				WithHeader("Authorization", "Bearer "+tokenAdm).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().ContainsKey("error").
				Value("error").String().Contains(tc.wrongParamName)
		})
	}
}
