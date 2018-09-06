package grimhilde

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Redirector implements a redirector from vanity Go import paths
// to their real source for both humans and `go get`.
type Redirector struct {
	VCS        string
	VCSBaseURL *url.URL
}

// ServeHTTP implements http.Handler.
func (rd *Redirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public")

	if r.URL.Query().Get("go-get") == "1" {
		rd.redirectGoGet(w, r)
	} else {
		rd.redirect(w, r)
	}
}

// See https://golang.org/cmd/go/#hdr-Remote_import_paths.
func (rd *Redirector) goGetImportTag(r *http.Request) string {
	repoName := leadingPathElement(r.URL.Path)
	repoImport := rd.vanityImport(r.Host, repoName)
	vcsURL := rd.vcsURL(repoName)
	return fmt.Sprintf(`<meta name="go-import" content="%v %v %v">`, repoImport, rd.VCS, vcsURL)
}

func (rd *Redirector) vanityImport(host, repoName string) string {
	return path.Join(host, repoName)
}

func (rd *Redirector) vcsURL(repoName string) string {
	vcsURL := *rd.VCSBaseURL
	vcsURL.Path = path.Join(vcsURL.Path, repoName)
	return vcsURL.String()
}

func (rd *Redirector) redirectGoGet(w http.ResponseWriter, r *http.Request) {
	tag := rd.goGetImportTag(r)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, tag)
}

func (rd *Redirector) redirect(w http.ResponseWriter, r *http.Request) {
	repoName := leadingPathElement(r.URL.Path)
	vcsURL := rd.vcsURL(repoName)
	// We want to send a StatusTemporaryRedirect and not a StatusSeeOther
	// because we want the method to stay the same and we do not want a StatusFound
	// because we want to be explicit about the method remaining the same.
	http.Redirect(w, r, vcsURL, http.StatusTemporaryRedirect)
}

// leadingPathElement returns the leading element of the path without the root slash.
func leadingPathElement(p string) string {
	// The path is always absolute unless some middleware edited it.
	p = strings.TrimPrefix(p, "/")

	i := strings.Index(p, "/")
	if i < 0 {
		return p
	}
	return p[:i]
}
