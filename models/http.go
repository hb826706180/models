package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HttpGet got请求
func HttpGet(url1 string) string {
	// 创建请求对象
	request, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return ""
	}

	// 设置请求头
	request.Header.Set("User-Agent", "Mozilla/5.0 (...) Gecko/20100101 Firefox/68.0")

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return ""
	}
	// 闭包函数
	defer response.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}

	return string(body)
}

func HttpPost(url1 string, data map[string]interface{}, headers map[string]string) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("序列化请求数据失败: %v", err)
	}

	req, err := http.NewRequest("POST", url1, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

//func HttpPost(url1 string, data map[string]string, headers map[string]string) (string, error) {
//	// 创建表单数据
//	form := url.Values{}
//	for key, value := range data {
//		form.Set(key, value)
//	}
//
//	// 创建请求对象
//	request, err := http.NewRequest("POST", url1, nil)
//	if err != nil {
//		return "", fmt.Errorf("创建请求失败: %v", err)
//	}
//
//	// 设置请求头
//	for key, value := range headers {
//		request.Header.Set(key, value)
//	}
//
//	// 设置请求体
//	request.PostForm = form
//
//	// 发送请求
//	client := &http.Client{}
//	response, err := client.Do(request)
//	if err != nil {
//		return "", fmt.Errorf("发送请求失败: %v", err)
//	}
//	defer response.Body.Close()
//
//	// 读取响应的内容
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return "", fmt.Errorf("读取响应失败: %v", err)
//	}
//	return string(body), nil
//}
