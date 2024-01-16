package evolution_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/freeverseio/laos-universal-node/internal/core/processor/evolution"
	"github.com/freeverseio/laos-universal-node/internal/core/processor/evolution/mock"
)

func TestLatestFinalizedBlockHash(t *testing.T) {
	t.Parallel()
	t.Run("Successful request to Laos parachain", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := mock.NewMockHTTPClient(ctrl)

		laosHTTP := evolution.NewLaosHTTP(mockHTTPClient, "http://caladan.com/own")

		// Create a sample JSON payload
		payload := []byte(`{"jsonrpc":"2.0","result":"sample-result","id":1}`)

		// Create a new httptest.ResponseRecorder
		recorder := httptest.NewRecorder()

		// Write the payload to the response body
		recorder.Body.Write(payload)
		// Write Status Code to the response body
		recorder.WriteHeader(http.StatusOK)

		// Create a new http.Response using the response recorder
		response := recorder.Result()

		defer func() {
			if err := response.Body.Close(); err != nil {
				t.Fatal("unexpected error closing response body")
			}
		}()

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		_, err := laosHTTP.LatestFinalizedBlockHash()
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
	})
	t.Run("unexpected status code from request to Laos parachain", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := mock.NewMockHTTPClient(ctrl)

		url := "http://caladan.com/own"

		laosHTTP := evolution.NewLaosHTTP(mockHTTPClient, "http://caladan.com/own")

		// Create a sample JSON payload

		// Create a new httptest.ResponseRecorder
		recorder := httptest.NewRecorder()

		// Write Status Code to the response body
		recorder.WriteHeader(http.StatusBadRequest)

		// Create a new http.Response using the response recorder
		response := recorder.Result()

		defer func() {
			if err := response.Body.Close(); err != nil {
				t.Fatal("unexpected error closing response body")
			}
		}()

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		expectedErr := fmt.Errorf("error in request to %s, got status code: %d", url, http.StatusBadRequest)
		_, err := laosHTTP.LatestFinalizedBlockHash()
		if err.Error() != expectedErr.Error() {
			t.Fatalf(`got error "%s", expected "%s"`, err.Error(), expectedErr.Error())
		}
	})
}

func TestBlockNumber(t *testing.T) {
	t.Parallel()
	t.Run("Successful request to Laos parachain", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := mock.NewMockHTTPClient(ctrl)

		laosHTTP := evolution.NewLaosHTTP(mockHTTPClient, "http://caladan.com/own")

		body := evolution.ChainGetBlock{
			JSONRPC: "2.0",
			Result: struct {
				Block evolution.Block `json:"block"`
			}{
				Block: evolution.Block{
					Header: evolution.BlockHeader{
						ParentHash:     "0xb0a5b16695d82b00c9f17013e530e46193cf1d31c9c5dbad26f6826d6bc5bcd9",
						Number:         "0x8413a",
						StateRoot:      "0xbb08bb4e145254a7255cbbd00f58eec4913e4651d8fa06534c5e0600f49dda97",
						ExtrinsicsRoot: "0x56d19c2f963be566164daf83d233d0081fa7052c91cd1e3bee84ac118ca0c034",
					},
				},
			},
			ID: 1,
		}

		var err error
		var payload []byte
		// Create a sample JSON payload
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("unexpected error marshaling response body: %s", err.Error())
		}

		// Create a new httptest.ResponseRecorder
		recorder := httptest.NewRecorder()

		// Write the payload to the response body
		recorder.Body.Write(payload)
		// Write Status Code to the response body
		recorder.WriteHeader(http.StatusOK)

		// Create a new http.Response using the response recorder
		response := recorder.Result()

		defer func() {
			if err = response.Body.Close(); err != nil {
				t.Fatal("unexpected error closing response body")
			}
		}()

		mockHTTPClient.EXPECT().Do(gomock.Any()).Return(response, nil).Times(1)

		_, err = laosHTTP.BlockNumber("0x95207a95aaf6c516017758f2fd4b7e173fb5a3fb56d3b0cdc0044cd0a9553f38")
		if err != nil {
			t.Fatalf("got %T, expected nil error", err)
		}
	})
}
