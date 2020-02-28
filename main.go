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
	data, err := ioutil.ReadFile("city.txt")
	if nil != err {
		log.Fatal(err)
	}

	var output []map[string]interface{}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if 3 != len(parts) {
			continue
		}
		country, province, city := parts[0], parts[1], parts[2]
		province = strings.ReplaceAll(province, "省", "")
		city = strings.ReplaceAll(city, "市", "")
		lat, lng := query(country, province, city)
		resultLine := map[string]interface{}{
			"country":  country,
			"province": province,
			"city":     city,
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

func query(country, province, city string) (latitude, longitude string) {
	api := "http://api.map.baidu.com/geocoding/v3/?address=" + province + city + "&output=json&ak=" + baiduAK
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
