package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ViaCEPResponse struct {
	Cep         string      `json:"cep"`
	Logradouro  string      `json:"logradouro"`
	Complemento string      `json:"complemento"`
	Bairro      string      `json:"bairro"`
	Localidade  string      `json:"localidade"`
	UF          string      `json:"uf"`
	Erro        interface{} `json:"erro,omitempty"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/weather/", weatherHandler)
	http.HandleFunc("/", healthHandler)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrair CEP da URL
	path := strings.TrimPrefix(r.URL.Path, "/weather/")
	cep := strings.TrimSpace(path)
	
	log.Printf("Received request for CEP: %s", cep)

	// Validar formato do CEP (8 dígitos)
	if !isValidCEP(cep) {
		log.Printf("Invalid CEP format: %s", cep)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	// Buscar localização pelo CEP
	location, err := getLocationByCEP(cep)
	if err != nil {
		if err.Error() == "CEP not found" {
			log.Printf("CEP not found: %s", cep)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})
		} else {
			log.Printf("ERROR: Failed to get location for CEP %s: %v", cep, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
		}
		return
	}
	
	log.Printf("Found location for CEP %s: %s", cep, location)

	// Buscar temperatura pela localização
	tempC, err := getTemperature(location)
	if err != nil {
		log.Printf("ERROR: Failed to get temperature for location '%s': %v", location, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "error fetching weather data"})
		return
	}

	// Converter temperaturas
	tempF := celsiusToFahrenheit(tempC)
	tempK := celsiusToKelvin(tempC)

	log.Printf("Successfully processed CEP %s: %.1f°C, %.1f°F, %.1f°K", cep, tempC, tempF, tempK)

	// Retornar resposta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	})
}

func isValidCEP(cep string) bool {
	// Remove hífens se houver
	cep = strings.ReplaceAll(cep, "-", "")
	
	// Verifica se tem exatamente 8 dígitos
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func getLocationByCEP(cep string) (string, error) {
	// Remove hífens do CEP
	cep = strings.ReplaceAll(cep, "-", "")

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP not found")
	}

	var viaCEP ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
		return "", err
	}

	// ViaCEP retorna um campo "erro": true quando o CEP não existe
	// O campo pode ser bool ou string, então verificamos também se a localidade está vazia
	if viaCEP.Erro != nil || viaCEP.Localidade == "" {
		return "", fmt.Errorf("CEP not found")
	}

	// Retorna a cidade e estado
	return fmt.Sprintf("%s,%s", viaCEP.Localidade, viaCEP.UF), nil
}

func getTemperature(location string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("ERROR: WEATHER_API_KEY not set")
		return 0, fmt.Errorf("weather API key not configured")
	}

	// URL encode da localização para evitar problemas com caracteres especiais
	encodedLocation := url.QueryEscape(location)
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, encodedLocation)
	log.Printf("Fetching weather for location: %s", location)
	
	resp, err := http.Get(weatherURL)
	if err != nil {
		log.Printf("ERROR: Failed to fetch weather data: %v", err)
		return 0, fmt.Errorf("failed to connect to weather API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: Weather API returned status %d for location: %s", resp.StatusCode, location)
		
		// Tentar ler o corpo da resposta para mais detalhes
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			log.Printf("Weather API error details: %+v", errorResp)
		}
		
		return 0, fmt.Errorf("weather API error: status %d", resp.StatusCode)
	}

	var weatherAPI WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPI); err != nil {
		log.Printf("ERROR: Failed to decode weather API response: %v", err)
		return 0, fmt.Errorf("failed to parse weather data: %v", err)
	}

	log.Printf("Successfully fetched temperature for %s: %.1f°C", location, weatherAPI.Current.TempC)
	return weatherAPI.Current.TempC, nil
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
