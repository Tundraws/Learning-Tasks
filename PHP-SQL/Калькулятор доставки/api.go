package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	citiesAPI   = "https://exercise.develop.maximaster.ru/service/city/"
	deliveryAPI = "https://exercise.develop.maximaster.ru/service/delivery/"
	apiUser     = "cli"
	apiPassword = "12344321"
)

type DeliveryResponse struct {
	Price   int    `json:"price"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func createHTTPClient() *http.Client {
	// ТОЛЬКО ДЛЯ РАЗРАБОТКИ! Не использовать в production!
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
}

// func createRequest(url string) (*http.Request, error) {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.SetBasicAuth(apiUser, apiPassword)
// 	return req, nil
// }

func fetchCitiesFromAPI() ([]string, error) {
	req, err := http.NewRequest("GET", citiesAPI, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.SetBasicAuth(apiUser, apiPassword)

	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернуло статус %d", resp.StatusCode)
	}

	var cities []string
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %v", err)
	}

	return cities, nil
}
func calculateDelivery(city, weight string) (*DeliveryResponse, error) {
	// Проверяем входные параметры
	if city == "" || weight == "" {
		return nil, fmt.Errorf("город и вес обязательны для заполнения")
	}

	// Создаем URL с параметрами
	u, err := url.Parse(deliveryAPI)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания URL: %v", err)
	}

	q := u.Query()
	q.Set("city", city)
	q.Set("weight", weight)
	u.RawQuery = q.Encode()

	// Создаем запрос
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.SetBasicAuth(apiUser, apiPassword)

	// Выполняем запрос
	client := createHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API вернуло статус %d: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var result DeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %v", err)
	}

	// Проверяем статус ответа
	if result.Status == "error" {
		return nil, fmt.Errorf("ошибка API: %s", result.Message)
	}

	return &result, nil
}
