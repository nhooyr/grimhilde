package grimhilde_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"go.nhooyr.io/grimhilde/internal/grimhilde"
)

func TestRedirector(t *testing.T) {
	t.Parallel()

	rd := &grimhilde.Redirector{
		VCS: "git",
		VCSBaseURL: &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "/nhooyr",
		},
	}

	t.Run("goGet", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name string
			path string
			body string
		}{
			{
				name: "root",
				path: "/",
				body: `<meta name="go-import" content="go.bar.io git https://github.com/nhooyr">`,
			},
			{
				name: "package",
				path: "/grimhilde",
				body: `<meta name="go-import" content="go.bar.io/grimhilde git https://github.com/nhooyr/grimhilde">`,
			},
			{
				name: "subpackage",
				path: "/grimhilde/meow/bar",
				body: `<meta name="go-import" content="go.bar.io/grimhilde git https://github.com/nhooyr/grimhilde">`,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				resp := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "http://go.bar.io"+tc.path+"?go-get=1", nil)
				rd.ServeHTTP(resp, req)

				if resp.Body.String() != tc.body {
					t.Errorf("got body %q; expected body %q", resp.Body, tc.body)
				}

			})
		}
	})

	testCases := []struct {
		name     string
		path     string
		location string
	}{
		{
			name:     "root",
			path:     "/",
			location: "https://github.com/nhooyr",
		},
		{
			name:     "package",
			path:     "/grimhilde",
			location: "https://godoc.org/go.bar.io/grimhilde",
		},
		{
			name:     "subpackage",
			path:     "/grimhilde/bar/meow/foo",
			location: "https://godoc.org/go.bar.io/grimhilde",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http://go.bar.io"+tc.path, nil)
			rd.ServeHTTP(resp, req)

			gotLocation := resp.Header().Get("Location")
			if gotLocation != tc.location {
				t.Errorf("got location %q; expected %q", gotLocation, tc.location)
			}
		})
	}
}
