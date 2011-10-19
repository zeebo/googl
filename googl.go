package googl

import (
	"url"
	"http"
	"os"
	"json"
	"bytes"
)

type Googl struct {
	key string
}

func New(key string) *Googl {
	return &Googl{key: key}
}

func (g *Googl) url() *url.URL {
	values := make(url.Values)
	values.Add("key", g.key)

	return &url.URL{
		Host:     "www.googleapis.com",
		Path:     "/urlshortener/v1/url",
		RawQuery: values.Encode(),
		Scheme:   "https",
	}
}

func (g *Googl) Shorten(u *url.URL) (*url.URL, os.Error) {
	buf := new(bytes.Buffer)
	var request struct {
		LongUrl string `json:"longUrl"`
	}
	request.LongUrl = u.String()

	enc := json.NewEncoder(buf)
	if err := enc.Encode(request); err != nil {
		return nil, err
	}

	r, err := http.Post(g.url().String(), "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var response struct {
		Kind    string
		Id      string
		LongUrl string
	}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&response); err != nil {
		return nil, err
	}

	return url.Parse(response.Id)
}
