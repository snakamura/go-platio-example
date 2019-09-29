package platio_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"go-platio-example/mock_platio"
	. "go-platio-example/platio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const authorization = "auth"

func TestLatestRecord(t *testing.T) {
	t.Run("should get the latest record", func(t *testing.T) {
		records := `[
          {
            "id": "r11111111111111111111111111",
            "values": {
                "cd33ed98": {
                    "type": "String",
                    "value": "abc"
                },
                "ce0f2361": {
                    "type": "Number",
                    "value": 30
                }
            }
          }
        ]`
		expected := &Record{
			Id: "r11111111111111111111111111",
			Values: Values{
				Name: &StringValue{"abc"},
				Age:  &NumberValue{30},
			},
		}

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.String() != "/records?limit=1" {
				t.Error("Invalid request path", req.URL)
			}
			if req.Header.Get("Authorization") != authorization {
				t.Error("Invalid Authorization header", req.Header)
			}
			res.WriteHeader(200)
			res.Write([]byte(records))
		}))
		defer testServer.Close()

		api := NewAPI(testServer.URL, authorization)
		record, err := api.GetLatestRecord()
		if err != nil || !reflect.DeepEqual(record, expected) {
			t.Error("It should return the latest record", err, record)
		}
	})

	t.Run("should return nil if no records", func(t *testing.T) {
		records := `[]`
		var expected *Record = nil

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(200)
			res.Write([]byte(records))
		}))
		defer testServer.Close()

		api := NewAPI(testServer.URL, authorization)
		record, err := api.GetLatestRecord()
		if err != nil || !reflect.DeepEqual(record, expected) {
			t.Error("It should return nil if no records", err, record)
		}
	})

	t.Run("should treat an error response as an error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(400)
			res.Write([]byte(""))
		}))
		defer testServer.Close()

		api := NewAPI(testServer.URL, authorization)
		record, err := api.GetLatestRecord()
		if err == nil || err.Error() != "400 Bad Request" || record != nil {
			t.Error("It should return an error on an error response", err, record)
		}
	})
}

func TestUpdateRecord(t *testing.T) {
	t.Run("should update the specified record", func(t *testing.T) {
		id := "r11111111111111111111111111"
		values := &Values{
			Age: &NumberValue{50},
		}

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.String() != "/records/"+id {
				t.Error("Invalid request path", req.URL)
			}
			if req.Header.Get("Authorization") != authorization {
				t.Error("Invalid Authorization header", req.Header)
			}
			body, _ := ioutil.ReadAll(req.Body)
			expected := `{"values":{"ce0f2361":{"type":"Number","value":50}},"replace":false}`
			if string(body) != expected {
				t.Error("Invalid request body", string(body))
			}
			res.WriteHeader(200)
			res.Write([]byte(`{}`))
		}))
		defer testServer.Close()

		api := NewAPI(testServer.URL, authorization)
		if err := api.UpdateRecord(id, values); err != nil {
			t.Error("It should update the record", err)
		}
	})
}

func TestSendRequest(t *testing.T) {
	t.Run("should set proper headers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mock_platio.NewMockhttpClient(ctrl)
		mockClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 200}, nil)

		api := NewTestAPI(mockClient, "https://api.plat.io/xyz", "auth")

		req, _ := http.NewRequest("GET", "test", nil)
		res, _ := SendRequestTest(api, req)

		if res.StatusCode != 200 {
			t.Error("It should return a response")
		}

		if req.Header.Get("Authorization") != "auth" {
			t.Error("It should set Authorization header")
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("It should set Content-Type header")
		}
	})

	t.Run("should turn an error if error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mock_platio.NewMockhttpClient(ctrl)
		mockClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("error"))

		api := NewTestAPI(mockClient, "https://api.plat.io/xyz", "auth")

		req, _ := http.NewRequest("GET", "test", nil)
		_, err := SendRequestTest(api, req)

		if err.Error() != "error" {
			t.Error("It should return an error", err)
		}
	})

	t.Run("should turn an error response to an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mock_platio.NewMockhttpClient(ctrl)
		mockClient.EXPECT().Do(gomock.Any()).Return(&http.Response{Status: "400 Bad Request", StatusCode: 400}, nil)

		api := NewTestAPI(mockClient, "https://api.plat.io/xyz", "auth")

		req, _ := http.NewRequest("GET", "test", nil)
		_, err := SendRequestTest(api, req)

		if err.Error() != "400 Bad Request" {
			t.Error("It should return an error", err)
		}
	})
}
