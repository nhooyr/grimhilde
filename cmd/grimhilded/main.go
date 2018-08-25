package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/url"

	"github.com/nhooyr/grimhilde/internal/grimhilde"
	"google.golang.org/appengine"
)

type config struct {
	VCS        string `json:"vcs"`
	VCSBaseURL string `json:"vcs_base_url"`
}

func main() {
	stdlog.SetFlags(0)

	rd, err := redirector()
	if err != nil {
		stdlog.Fatalf("failed to create redirector: %v", err)
	}

	http.Handle("/", rd)

	appengine.Main()
}

func redirector() (http.Handler, error) {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.json: %v", err)
	}

	var c config
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.json into redirector: %v", err)
	}

	vcsBaseURL, err := url.Parse(c.VCSBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vcs_base_url: %v", err)
	}

	rd := &grimhilde.Redirector{
		VCSBaseURL: vcsBaseURL,
		VCS:        c.VCS,
	}

	return rd, nil
}
