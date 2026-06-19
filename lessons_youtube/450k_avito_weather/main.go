package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// имеется сторонний сервис погоды (его имитация это фнкция WeatherForecast)
// сторонний сервис работает за секунду, что для нас долго
// На наш сервис идет большая нагрузка. Как доработать текущую реализацию?
// 1. Предложить и реализовать решение
// 2. Дополнительрое задание: если будет несколько городов
// Доработать задачу с учетом 2 го пункта
//
// func WeatherForecast() int {
//	time.Sleep(1 * time.Second)
//	return rand.Intn(70) - 30
//}

//func main() {
//  http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, `{"temperature":%d}`, WeatherForecast())
//	})
//	if err := http.ListenAndServe(":3333", nil); err != nil {
//		panic(err)
//	}
//}
//
//

type Data struct {
	Temperature int
	mu          sync.RWMutex
}

func NewData(interval time.Duration) *Data {
	ticker := time.NewTicker(interval)
	newData := &Data{}

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			newData.UpdateTemprature()
		}
	}()
	return newData
}

func (d *Data) UpdateTemprature() {
	tmp := WeatherForecast()
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Temperature = tmp
}

func (d *Data) GetTemperature() int {
	d.mu.RLock()
	defer d.mu.Unlock()
	return d.Temperature
}

func WeatherForecast() int {
	time.Sleep(1 * time.Second)
	return rand.Intn(70) - 30
}

func main() {
	data := NewData(1 * time.Minute)

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"temperature":%d}`, data.GetTemperature())
	})
	if err := http.ListenAndServe(":3333", nil); err != nil {
		panic(err)
	}
}
