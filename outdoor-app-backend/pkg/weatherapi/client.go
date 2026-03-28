package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiKey = "4d27c704eb5cd677f01d612c5edf3c40"
	apiUrl = "http://apis.juhe.cn/simpleWeather/query"
)

// JuheWeatherResponse 聚合数据 API 返回格式
type JuheWeatherResponse struct {
	ErrorCode int    `json:"error_code"`
	Reason    string `json:"reason"`
	Result    *struct {
		City   string `json:"city"`
		Future []struct {
			Date        string `json:"date"`
			Temperature string `json:"temperature"`
			Weather     string `json:"weather"` // 如：雷阵雨、大雪
			Direct      string `json:"direct"`
		} `json:"future"`
	} `json:"result"`
}

// FetchFromJuhe 向聚合数据发起真实的 HTTP 请求，获取未来 5 天天气
func FetchFromJuhe(city string) (*JuheWeatherResponse, error) {
	encodedCity := url.QueryEscape(city)
	reqUrl := fmt.Sprintf("%s?city=%s&key=%s", apiUrl, encodedCity, apiKey)

	// 1. 发起请求，设置 5 秒超时防止卡死
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 2. 解析 JSON
	var weatherRes JuheWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherRes); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	// 3. 校验业务状态码 (ErrorCode 0 表示成功)
	if weatherRes.ErrorCode != 0 || weatherRes.Result == nil {
		return nil, fmt.Errorf("API 业务错误: %s", weatherRes.Reason)
	}

	return &weatherRes, nil
}

// IsBadWeather 核心预警规则：判断该天气描述是否属于“恶劣天气”
func IsBadWeather(weatherDesc string) bool {
	badKeywords := []string{"雨", "雪", "雷", "暴", "冰雹", "雾霾", "沙尘", "台风"}
	for _, kw := range badKeywords {
		if strings.Contains(weatherDesc, kw) {
			return true
		}
	}
	return false
}
