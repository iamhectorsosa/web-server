package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		statusCode   int
		contentType  string
		want         ErrorResponse
	}{
		{
			name:         "responds with bad request error",
			errorMessage: "Oops! Bad Request",
			statusCode:   http.StatusBadRequest,
			contentType:  "application/json",
			want:         ErrorResponse{Error: "Oops! Bad Request"},
		},
		{
			name:         "responds with not found error",
			errorMessage: "Not Found",
			statusCode:   http.StatusNotFound,
			contentType:  "application/json",
			want:         ErrorResponse{Error: "Not Found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			respondWithError(response, tt.statusCode, tt.errorMessage)

			var got ErrorResponse
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Fatalf("error decoding JSON response: %v", err)
			}

			AssertResponseBody(t, got, tt.want)
			AssertResponseCode(t, response.Code, tt.statusCode)
			AssertResponseHeader(t, response.Header().Get("Content-Type"), tt.contentType)
		})
	}
}

func TestRespondWithJSON(t *testing.T) {
	type Payload struct {
		Body string `json:"body"`
	}

	tests := []struct {
		name        string
		statusCode  int
		contentType string
		want        Payload
	}{
		{
			name:        "responds with a successful request",
			statusCode:  http.StatusOK,
			contentType: "application/json",
			want:        Payload{Body: "Hello World"},
		},
		{
			name:        "responds with a successful request for a resource created",
			statusCode:  http.StatusCreated,
			contentType: "application/json",
			want:        Payload{Body: "New Resource"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			respondWithJSON(response, tt.statusCode, tt.want)

			var got Payload
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Fatalf("error decoding JSON response: %v", err)
			}

			AssertResponseBody(t, got, tt.want)
			AssertResponseCode(t, response.Code, tt.statusCode)
			AssertResponseHeader(t, response.Header().Get("Content-Type"), tt.contentType)
		})
	}
}

func AssertResponseBody(t testing.TB, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Response Body, got: %+v, want: %+v", got, want)
	}
}

func AssertResponseCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("Response Code, got: %d, want: %d", got, want)
	}
}

func AssertResponseHeader(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("Response Header, got: %s, want: %s", got, want)
	}
}
