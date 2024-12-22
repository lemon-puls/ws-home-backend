package maputils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"ws-home-backend/config"
)

// RegeoResponse 逆地理编码响应结构体
type RegeoResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
			Township string `json:"township"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

// AddressInfo 地址信息
type AddressInfo struct {
	FormattedAddress string `json:"formatted_address"` // 完整地址
	Province         string `json:"province"`          // 省份
	City             string `json:"city"`              // 城市
	District         string `json:"district"`          // 区县
	Township         string `json:"township"`          // 乡镇
}

// GetAddressFromLocation 根据经纬度获取地址信息
func GetAddressFromLocation(longitude, latitude string) (*AddressInfo, error) {
	if longitude == "" || latitude == "" {
		return nil, fmt.Errorf("经纬度不能为空")
	}

	// 高德侧要求 经纬度不得超过 6 位小数
	// 检查并处理经度的小数位数
	if strings.Contains(longitude, ".") {
		parts := strings.Split(longitude, ".")
		if len(parts[1]) > 6 {
			longitude = parts[0] + "." + parts[1][:6]
		}
	}

	// 检查并处理纬度的小数位数
	if strings.Contains(latitude, ".") {
		parts := strings.Split(latitude, ".")
		if len(parts[1]) > 6 {
			latitude = parts[0] + "." + parts[1][:6]
		}
	}

	// 构建请求URL
	location := fmt.Sprintf("%s,%s", longitude, latitude)
	url := fmt.Sprintf("%s?key=%s&location=%s&extensions=base&output=json",
		config.Conf.AmapConfig.RegeoURL,
		config.Conf.AmapConfig.Key,
		location)

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求高德地图 API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 解析JSON响应
	var regeoResp RegeoResponse
	if err := json.Unmarshal(body, &regeoResp); err != nil {
		return nil, fmt.Errorf("解析响应内容失败: %v", err)
	}

	// 检查响应状态
	if regeoResp.Status != "1" {
		return nil, fmt.Errorf("高德地图API返回错误: %s", regeoResp.Info)
	}

	// 构建返回结果
	addressInfo := &AddressInfo{
		FormattedAddress: regeoResp.Regeocode.FormattedAddress,
		Province:         regeoResp.Regeocode.AddressComponent.Province,
		City:             regeoResp.Regeocode.AddressComponent.City,
		District:         regeoResp.Regeocode.AddressComponent.District,
		Township:         regeoResp.Regeocode.AddressComponent.Township,
	}

	return addressInfo, nil
}
