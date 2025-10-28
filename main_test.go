package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{"Valid CEP with 8 digits", "01310100", true},
		{"Valid CEP with hyphen", "01310-100", true},
		{"Invalid CEP with 7 digits", "0131010", false},
		{"Invalid CEP with 9 digits", "013101000", false},
		{"Invalid CEP with letters", "0131010a", false},
		{"Empty CEP", "", false},
		{"CEP with spaces", "01310 100", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCEP(tt.cep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{25, 77},
	}

	for _, tt := range tests {
		result := celsiusToFahrenheit(tt.celsius)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 273},
		{-273, 0},
		{25, 298},
		{100, 373},
	}

	for _, tt := range tests {
		result := celsiusToKelvin(tt.celsius)
		assert.Equal(t, tt.expected, result)
	}
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedMsg    string
	}{
		{"CEP with letters", "0131010a", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"CEP too short", "0131010", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"CEP too long", "013101000", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"Empty CEP", "", http.StatusUnprocessableEntity, "invalid zipcode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/weather/"+tt.cep, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(weatherHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

func TestWeatherHandler_CEPNotFound(t *testing.T) {
	// CEP válido no formato mas que não existe
	req, err := http.NewRequest("GET", "/weather/99999999", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", response.Message)
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

// Teste de integração - requer WEATHER_API_KEY configurada
func TestWeatherHandler_ValidCEP_Integration(t *testing.T) {
	// Este teste só roda se a variável de ambiente estiver configurada
	// Usa um CEP válido de São Paulo
	t.Skip("Integration test - requires WEATHER_API_KEY")

	req, err := http.NewRequest("GET", "/weather/01310100", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response WeatherResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	
	// Verifica se as temperaturas foram calculadas corretamente
	expectedF := response.TempC*1.8 + 32
	expectedK := response.TempC + 273
	
	assert.Equal(t, expectedF, response.TempF)
	assert.Equal(t, expectedK, response.TempK)
}
