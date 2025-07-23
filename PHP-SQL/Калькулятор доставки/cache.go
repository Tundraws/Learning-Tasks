package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var citiesCacheFile = "cities.cache"

func initCityCache() {
	// Проверяем нужно ли обновить кэш
	if shouldRefreshCache() {
		if err := updateCityCache(); err != nil {
			log.Printf("Не удалось обновить кэш городов: %v", err)
		}
	}
}

func shouldRefreshCache() bool {
	info, err := os.Stat(citiesCacheFile)
	if os.IsNotExist(err) {
		return true
	}

	// Если файл старше одного дня
	return time.Since(info.ModTime()) > 24*time.Hour
}

func updateCityCache() error {
	cities, err := fetchCitiesFromAPI()
	if err != nil {
		return fmt.Errorf("ошибка получения городов: %v", err)
	}

	data, err := json.Marshal(cities)
	if err != nil {
		return fmt.Errorf("ошибка кодирования городов: %v", err)
	}

	err = ioutil.WriteFile(citiesCacheFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка сохранения кэша: %v", err)
	}

	return nil
}

func getCities() ([]string, error) {
	// Пытаемся получить города из API
	cities, err := fetchCitiesFromAPI()
	if err != nil {
		log.Printf("Failed to fetch cities from API: %v", err)

		// Пробуем загрузить из кэша
		cachedCities, cacheErr := loadCitiesFromCache()
		if cacheErr != nil {
			log.Printf("Failed to load cities from cache: %v", cacheErr)
			// Возвращаем минимальный набор городов
			return []string{"Москва", "Санкт-Петербург", "Тула"}, nil
		}
		return cachedCities, nil
	}

	// Сохраняем в кэш
	if err := saveCitiesToCache(cities); err != nil {
		log.Printf("Failed to save cities to cache: %v", err)
	}

	return cities, nil
}

func loadCitiesFromCache() ([]string, error) {
	data, err := ioutil.ReadFile(citiesCacheFile)
	if err != nil {
		return nil, err
	}

	var cities []string
	if err := json.Unmarshal(data, &cities); err != nil {
		return nil, err
	}

	return cities, nil
}

func saveCitiesToCache(cities []string) error {
	data, err := json.Marshal(cities)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(citiesCacheFile, data, 0644)
}
