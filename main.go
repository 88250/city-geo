// City Geo - 中国城市经纬度数据.
// Copyright (c) 2020-present, b3log.org
//
// Lute HTTP is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"strings"
)

const baiduAK = "XSVaW6UooxiXEFlaBOGDXFmIARffS5Oo"

func main() {
	generateCity()

	data, err := ioutil.ReadFile("city.txt")
	if nil != err {
		log.Fatal(err)
	}

	var output []map[string]interface{}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if 3 > len(parts) {
			continue
		}

		country, province, city := parts[0], parts[1], parts[2]
		area := ""
		if 4 == len(parts) {
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
	}

	resultData, err := json.MarshalIndent(output, "", "  ")
	if nil != err {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("data.json", resultData, 0644); nil != err {
		log.Fatal(err)
	}
	log.Println("completed")
}

func query(country, province, city, area string) (latitude, longitude string) {
	api := "http://api.map.baidu.com/geocoding/v3/?address=" + province + city + area + "&output=json&ak=" + baiduAK
	response, data, errors := gorequest.New().Get(api).EndBytes()
	if nil != errors {
		log.Printf("city [%s] result failed [%+v]", city, errors)
	}
	if 200 != response.StatusCode {
		log.Printf("city [%s] result response [%+v]", city, response)
	}

	responseData := map[string]interface{}{}
	if err := json.Unmarshal(data, &responseData); nil != err {
		log.Printf("city [%s] unmarshal failed [%+v]", city, responseData)
		return
	}

	if status := responseData["status"]; 0.0 != status {
		log.Printf("city [%s] result response data [%+v]", city, responseData)
		return
	}

	result := responseData["result"].(map[string]interface{})
	location := result["location"].(map[string]interface{})
	lat := location["lat"].(float64)
	lng := location["lng"].(float64)
	latitude, longitude = fmt.Sprint(lat), fmt.Sprint(lng)
	return
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
