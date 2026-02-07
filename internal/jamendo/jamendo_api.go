package jamendo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TracksResponse struct {
	Results []Track
}

type Track struct {
	Name       string `json:"name"`
	ArtistName string `json:"artist_name"`
	Audio      string `json:"audio"`
}

type JamendoClient struct {
	httpClient *http.Client
	baseUrl    string
	clientId   string
}

func NewJamendoClient() JamendoClient {
	return JamendoClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl:  "https://api.jamendo.com/v3.0",
		clientId: "6e1ba05b",
	}
}

func (j *JamendoClient) get(path string) ([]byte, error) {

	req, err := http.NewRequest(
		"GET",
		j.baseUrl+path,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := j.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)

}

func (j *JamendoClient) GetTrack(tag string) (TracksResponse, error) {

	apiEndpoint := fmt.Sprintf("/tracks/?client_id=%s&format=jsonpretty&limit=5&fuzzytags=%s", j.clientId, tag)

	body, err := j.get(apiEndpoint)

	if err != nil {
		return TracksResponse{}, err
	}

	var output TracksResponse

	err = json.Unmarshal(body, &output)

	if err != nil {
		return TracksResponse{}, err
	}

	return output, nil

}
