package utils

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	ipSearcher     *xdb.Searcher
	ipSearcherOnce sync.Once
	ipSearcherErr  error
)

// InitIPSearcher 初始化 IP 归属地查询（全量加载到内存，查询最快）
func InitIPSearcher(dbPath string) error {
	ipSearcherOnce.Do(func() {
		cBuff, err := xdb.LoadContentFromFile(dbPath)
		if err != nil {
			ipSearcherErr = fmt.Errorf("load ip2region xdb failed: %w", err)
			return
		}
		header, err := xdb.LoadHeaderFromBuff(cBuff)
		if err != nil {
			ipSearcherErr = fmt.Errorf("load ip2region header failed: %w", err)
			return
		}
		version, err := xdb.VersionFromHeader(header)
		if err != nil {
			ipSearcherErr = fmt.Errorf("detect ip2region version failed: %w", err)
			return
		}
		ipSearcher, ipSearcherErr = xdb.NewWithBuffer(version, cBuff)
	})
	return ipSearcherErr
}

// GetIPLocation 查询 IP 归属地，返回格式化后的地址字符串
// ip2region 原始格式: "国家|区域|省份|城市|ISP"
func GetIPLocation(ip string) string {
	if ipSearcher == nil {
		return ""
	}
	// 跳过本地回环地址
	if ip == "::1" || ip == "127.0.0.1" || ip == "" {
		return "本机"
	}
	region, err := ipSearcher.SearchByStr(ip)
	if err != nil {
		return ""
	}
	return formatRegion(region)
}

// formatRegion 格式化 ip2region 返回的原始地区字符串
// 原始: "中国|0|北京|北京|联通" → "中国 北京"
// 原始: "美国|0|加利福尼亚|洛杉矶|0" → "美国 加利福尼亚 洛杉矶"
func formatRegion(region string) string {
	parts := strings.Split(region, "|")
	if len(parts) < 5 {
		return region
	}
	country := parts[0]
	province := parts[2]
	city := parts[3]

	var result []string
	if country != "0" && country != "" {
		result = append(result, country)
	}
	if province != "0" && province != "" && province != country {
		result = append(result, province)
	}
	if city != "0" && city != "" && city != province {
		result = append(result, city)
	}
	if len(result) == 0 {
		return "未知"
	}
	return strings.Join(result, " ")
}
