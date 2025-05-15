package restclient

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RestClient struct {
	Client *http.Client
}

func New() *RestClient {
	return &RestClient{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *RestClient) Get(url string, result interface{}) error {
	resp, err := r.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	return json.Unmarshal(bodyBytes, &result)
}
