package platio

func NewTestAPI(client httpClient, collectionUrl string, authorization string) *API {
	return &API{
		client:        client,
		collectionUrl: collectionUrl,
		authorization: authorization,
	}
}

var SendRequestTest = (*API).sendRequest
