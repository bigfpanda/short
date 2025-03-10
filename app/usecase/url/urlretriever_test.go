package url

import (
	"short/app/entity"
	"short/app/usecase/repo"
	"testing"
	"time"

	"github.com/byliuyang/app/mdtest"
)

type urlMap = map[string]entity.URL

func TestUrlRetriever_GetUrlAfter(t *testing.T) {
	now := time.Now()
	before := now.Add(-5 * time.Second)
	after := now.Add(5 * time.Second)

	testCases := []struct {
		name        string
		urls        urlMap
		alias       string
		expiringAt  time.Time
		hasErr      bool
		expectedURL entity.URL
	}{
		{
			name:        "alias not found",
			urls:        urlMap{},
			alias:       "220uFicCJj",
			expiringAt:  now,
			hasErr:      true,
			expectedURL: entity.URL{},
		},
		{
			name: "url expired",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &before,
				},
			},
			alias:       "220uFicCJj",
			expiringAt:  now,
			hasErr:      true,
			expectedURL: entity.URL{},
		},
		{
			name: "url never expire",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: nil,
				},
			},
			alias:      "220uFicCJj",
			expiringAt: now,
			hasErr:     false,
			expectedURL: entity.URL{
				Alias:    "220uFicCJj",
				ExpireAt: nil,
			},
		},
		{
			name: "unexpired url found",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &after,
				},
			},
			alias:      "220uFicCJj",
			expiringAt: now,
			hasErr:     false,
			expectedURL: entity.URL{
				Alias:    "220uFicCJj",
				ExpireAt: &after,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeRepo := repo.NewURLFake(testCase.urls)
			retriever := NewRetrieverPersist(&fakeRepo)
			tracer := mdtest.NewTracerFake()
			fakeTrace := tracer.BeginTrace("GetUrlAfter")
			url, err := retriever.GetAfter(fakeTrace, testCase.alias, testCase.expiringAt)

			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
				return
			}
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURL, url)
		})
	}
}

func TestUrlRetriever_GetUrl(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name        string
		urls        urlMap
		alias       string
		hasErr      bool
		expectedURL entity.URL
	}{
		{
			name:        "alias not found",
			urls:        urlMap{},
			alias:       "220uFicCJj",
			hasErr:      true,
			expectedURL: entity.URL{},
		},
		{
			name: "valid url found",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &now,
				},
			},
			alias:  "220uFicCJj",
			hasErr: false,
			expectedURL: entity.URL{
				Alias:    "220uFicCJj",
				ExpireAt: &now,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeRepo := repo.NewURLFake(testCase.urls)
			retriever := NewRetrieverPersist(&fakeRepo)
			tracer := mdtest.NewTracerFake()
			fakeTrace := tracer.BeginTrace("GetUrl")
			url, err := retriever.Get(fakeTrace, testCase.alias)

			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
				return
			}
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURL, url)
		})
	}
}
