package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Erro        bool   `json:"erro"`
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

	// Validar formato do CEP (8 dígitos)
	if !isValidCEP(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	// Buscar localização pelo CEP
	location, err := getLocationByCEP(cep)
	if err != nil {
		if err.Error() == "CEP not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
		}
		return
	}

	// Buscar temperatura pela localização
	tempC, err := getTemperature(location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "error fetching weather data"})
		return
	}

	// Converter temperaturas
	tempF := celsiusToFahrenheit(tempC)
	tempK := celsiusToKelvin(tempC)

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
	if viaCEP.Erro {
		return "", fmt.Errorf("CEP not found")
	}

	// Retorna a cidade e estado
	return fmt.Sprintf("%s,%s", viaCEP.Localidade, viaCEP.UF), nil
}

func getTemperature(location string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		// Para testes locais, você pode usar uma chave de API padrão
		// Em produção, sempre use variáveis de ambiente
		log.Println("Warning: WEATHER_API_KEY not set")
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API error: status %d", resp.StatusCode)
	}

	var weatherAPI WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPI); err != nil {
		return 0, err
	}

	return weatherAPI.Current.TempC, nil
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
