package getexternal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

// CleanURL creates url for bitbucket API server
func CleanURL(bbu, wp, q string) (rel string) {
	u, err := url.Parse(bbu)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	p := strings.Split(u.Path, "/")

	newu := &url.URL{
		Scheme:   "https",
		Host:     "bitbucket.org",
		Path:     "api/2.0/repositories/" + p[1] + "/" + p[2] + "/" + wp,
		RawQuery: q,
	}

	return newu.String()
}

// GetBbRepo retrieves star count from watchers and last activity(not last commit).
func GetBbRepo(url string) (int, string) {
	w, p, u := "fields=size", "watchers", `fields=updated_on`

	type bb struct {
		Stars   int    `json:"size"`
		Updated string `json:"updated_on"`
		//Node_id int `json:"id"`
	}
	thisbb := bb{}

	resw, err := http.Get(CleanURL(url, p, w))
	if err != nil {
		log.Fatal().Stack().Err(err)
	}
	bodyw, err := ioutil.ReadAll(resw.Body)
	resw.Body.Close()
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	err = json.Unmarshal(bodyw, &thisbb)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}

	res, err := http.Get(CleanURL(url, "", u))
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}

	err = json.Unmarshal(body, &thisbb)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}

	return thisbb.Stars, strings.Split(thisbb.Updated, "T")[0]
}
