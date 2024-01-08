// City Geo - 中国城市经纬度数据.
// Copyright (c) 2020-present, b3log.org
//
// Lute HTTP is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.
// 2024-01-08 revised func main() and query(). 

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

const baiduAK = "XSVaW6UooxiXEFlaBOGDXFmIARffS5Oo"

func main() {
	logFile, err := os.OpenFile("query.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	generateCity()

	data, err := os.ReadFile("city.txt")
	if nil != err {
		log.Fatal(err)
	}

	var output []map[string]interface{}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		country, province, city := parts[0], parts[1], parts[2]
		area := ""
		if len(parts) == 4 {
			area = parts[3]
		}
		lat, lng := query(country, province, city, area)
		resultLine := map[string]interface{}{
			"country":  country,
			"province": province,
			"city":     city,
			"area":     area,
			"lat":      lat,
			"lng":      lng,
		}
		output = append(output, resultLine)
		log.Printf("query result %+v", resultLine)
		//fmt.Printf("query result %+v\n", resultLine)
	}

	resultData, err := json.MarshalIndent(output, "", "  ")
	if nil != err {
		log.Fatal(err)
	}

	if err := os.WriteFile("data.json", resultData, 0644); nil != err {
		log.Fatal(err)
	}

	log.Println("Program completed successfully")

}

func query(country, province, city, area string) (latitude, longitude string) {
	maxRetries := 3
	for retries := 0; retries < maxRetries; retries++ {
		api := "http://api.map.baidu.com/geocoding/v3/?address=" + province + city + area + "&ret_coordtype=gcj02ll&output=json&ak=" + baiduAK
		response, data, errors := gorequest.New().Get(api).EndBytes()

		// 检查是否有错误
		if errors != nil {
			fmt.Printf("Attempt %d: Request error for city [%s]: %+v\n", retries+1, city, errors)
			log.Printf("Attempt %d: Request error for city [%s]: %+v", retries+1, city, errors)
			continue
		}

		// 检查是否收到响应
		if response == nil {
			log.Printf("Attempt %d: No response received for city [%s]", retries+1, city)
			continue
		}

		//if errors == nil && response != nil && response.StatusCode == 200 {
		// 检查响应状态码
		if response.StatusCode == 302 {
			log.Fatalf("Received status code 302 for city [%s], terminating program.", city)
		} else if response.StatusCode != 200 {
			log.Printf("Attempt %d: Unexpected status code for city [%s]: %d", retries+1, city, response.StatusCode)
			time.Sleep(1 * time.Second)
			continue
		}

		// 解析响应数据
		responseData := map[string]interface{}{}
		if err := json.Unmarshal(data, &responseData); err != nil {
			log.Printf("Attempt %d: city [%s] unmarshal failed [%+v]", retries+1, city, err)
			time.Sleep(1 * time.Second)
			continue
		}

		status, ok := responseData["status"].(float64)
		if !ok || status != 0.0 {
			log.Printf("Attempt %d: city [%s] result response data [%+v]", retries+1, city, responseData)
			continue
		}

		result, ok := responseData["result"].(map[string]interface{})
		if !ok {
			log.Printf("Attempt %d: city [%s] result type assertion failed", retries+1, city)
			continue
		}

		location, ok := result["location"].(map[string]interface{})
		if !ok {
			log.Printf("Attempt %d: city [%s] location type assertion failed", retries+1, city)
			continue
		}

		lat, ok := location["lat"].(float64)
		if !ok {
			log.Printf("Attempt %d: city [%s] latitude type assertion failed", retries+1, city)
			continue
		}

		lng, ok := location["lng"].(float64)
		if !ok {
			log.Printf("Attempt %d: city [%s] longitude type assertion failed", retries+1, city)
			continue
		}

		latitude, longitude = fmt.Sprint(lat), fmt.Sprint(lng)
		return
	}

	log.Printf("Failed to query after %d attempts for city [%s]", maxRetries, city)
	return "", ""
}

func generateCity() {
	provincesData, err := ioutil.ReadFile("provinces.json")
	if nil != err {
		log.Fatal(err)
	}
	citiesData, err := ioutil.ReadFile("cities.json")
	if nil != err {
		log.Fatal(err)
	}
	areasData, err := ioutil.ReadFile("areas.json")
	if nil != err {
		log.Fatal(err)
	}

	var provinces []map[string]interface{}
	if err := json.Unmarshal(provincesData, &provinces); nil != err {
		log.Fatal(err)
	}

	var cities = []map[string]interface{}{}
	if err := json.Unmarshal(citiesData, &cities); nil != err {
		log.Fatal(err)
	}

	areas := []map[string]interface{}{}
	if err := json.Unmarshal(areasData, &areas); nil != err {
		log.Fatal(err)
	}

	var lines string
	for _, province := range provinces {
		provinceName := province["name"].(string)
		provinceCode := province["code"].(string)
		selectCities := getCities(provinceCode, cities)
		for _, city := range selectCities {
			cityName := city["name"].(string)
			cityCode := city["code"].(string)
			lines += "中国\t" + provinceName + "\t" + cityName + "\n"
			selectAreas := getAreas(cityCode, areas)
			for _, area := range selectAreas {
				areaName := area["name"].(string)

				lines += "中国\t" + provinceName + "\t" + cityName + "\t" + areaName + "\n"
			}
		}
	}

	ioutil.WriteFile("city.txt", []byte(lines), 0644)
}

func getCities(provinceCode string, cities []map[string]interface{}) (ret []map[string]interface{}) {
	for _, city := range cities {
		if city["provinceCode"].(string) == provinceCode {
			ret = append(ret, city)
		}
	}
	return
}

func getAreas(cityCode string, areas []map[string]interface{}) (ret []map[string]interface{}) {
	for _, area := range areas {
		if area["cityCode"].(string) == cityCode {
			ret = append(ret, area)
		}
	}
	return
}
