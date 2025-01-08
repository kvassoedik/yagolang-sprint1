package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"final/handler"
)

func TestCalculateHandler(t *testing.T) {
	handler := http.HandlerFunc(handler.CalculateHandler)

	tests := []struct {
		name           string
		input          map[string]string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Valid Expression",
			input: map[string]string{
				"expression": "(2+3)*4",
			},
			expectedStatus: 200,
			expectedBody:   map[string]string{"result": "20"},
		},
		{
			name: "Invalid Expression",
			input: map[string]string{
				"expression": "2+abc",
			},
			expectedStatus: 422,
			expectedBody:   map[string]string{"error": "Expression is not valid"},
		},
		{
			name: "Division by Zero",
			input: map[string]string{
				"expression": "10/0",
			},
			expectedStatus: 500,
			expectedBody:   map[string]string{"error": "Internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// создание JSON payload
			payload, _ := json.Marshal(tt.input)

			// создать запрос
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// записать ответ
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// проверить статус ответа
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// проверить тело ответа
			body, _ := ioutil.ReadAll(rr.Body)
			var response map[string]string
			json.Unmarshal(body, &response)

			for key, value := range tt.expectedBody {
				if response[key] != value {
					t.Errorf("Handler returned wrong body: got %v want %v", response[key], value)
				}
			}
		})
	}
}
