package coordinates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Структуры для парсинга JSON ответа
type GeoResponse struct {
	Results          []City  `json:"results"`
	GenerationTimeMs float64 `json:"generationtime_ms"`
}

type City struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CountryCode string  `json:"country_code"`
	Admin1ID    int     `json:"admin1_id"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryID   int     `json:"country_id"`
	Country     string  `json:"country"`
	Admin1      string  `json:"admin1"`
}

// Функция для поиска города
func SearchCity(cityName string) (*GeoResponse, error) {
	// Кодируем параметры запроса
	baseURL := "https://geocoding-api.open-meteo.com/v1/search"
	params := url.Values{}
	params.Add("name", cityName)
	params.Add("count", "1")
	params.Add("language", "ru")
	params.Add("format", "json")

	fullURL := baseURL + "?" + params.Encode()

	// Выполняем HTTP GET запрос
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Парсим JSON
	var geoResponse GeoResponse
	if err := json.Unmarshal(body, &geoResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return &geoResponse, nil
}
