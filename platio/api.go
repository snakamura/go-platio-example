package platio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type API struct {
	client        httpClient
	collectionUrl string
	authorization string
}

func NewAPI(collectionUrl string, authorization string) *API {
	return &API{
		client:        &http.Client{},
		collectionUrl: collectionUrl,
		authorization: authorization,
	}
}

func (api *API) GetLatestRecords(count int) ([]Record, error) {
	url := fmt.Sprintf("%s/records?limit=%d", api.collectionUrl, count)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := api.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var records []Record
	if err = json.NewDecoder(res.Body).Decode(&records); err != nil {
		return nil, err
	}

	return records, nil
}

func (api *API) UpdateRecord(id RecordId, values *Values) error {
	body, err := json.Marshal(struct {
		Values  *Values `json:"values"`
		Replace bool    `json:"replace"`
	}{values, false})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", api.collectionUrl+"/records/"+id, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = api.sendRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (api *API) sendRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", api.authorization)
	req.Header.Add("Content-Type", "application/json")
	res, err := api.client.Do(req)
	if err != nil {
		return res, err
	} else if res.StatusCode >= http.StatusBadRequest {
		return nil, ErrorResponse{res.Status}
	} else {
		return res, err
	}
}

type ErrorResponse struct {
	Status string
}

func (e ErrorResponse) Error() string {
	return e.Status
}
