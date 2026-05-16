package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
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
		candidates := buildIPDBCandidates(dbPath)
		var lastErr error

		for _, path := range candidates {
			cBuff, err := xdb.LoadContentFromFile(path)
			if err != nil {
				lastErr = err
				continue
			}

			header, err := xdb.LoadHeaderFromBuff(cBuff)
			if err != nil {
				lastErr = err
				continue
			}

			version, err := xdb.VersionFromHeader(header)
			if err != nil {
				lastErr = err
				continue
			}

			searcher, err := xdb.NewWithBuffer(version, cBuff)
			if err != nil {
				lastErr = err
				continue
			}

			ipSearcher = searcher
			ipSearcherErr = nil
			return
		}

		if lastErr != nil {
			ipSearcherErr = fmt.Errorf("load ip2region xdb failed, tried %v, last error: %w", candidates, lastErr)
			return
		}

		ipSearcherErr = fmt.Errorf("load ip2region xdb failed: no candidate path found")
	})
	return ipSearcherErr
}

// GetIPLocation 查询 IP 归属地，返回格式化后的地址字符串
// ip2region 原始格式: "国家|区域|省份|城市|ISP"
func GetIPLocation(ip string) string {
	normalizedIP := normalizeIP(ip)
	if normalizedIP == "" {
		return "未知"
	}

	parsed := net.ParseIP(normalizedIP)
	if parsed == nil {
		return "未知"
	}

	if parsed.IsLoopback() {
		return "本机"
	}

	if ipSearcher == nil {
		if parsed.To4() == nil {
			return "IPv6"
		}
		return "未知"
	}

	region, err := ipSearcher.SearchByStr(normalizedIP)
	if err != nil {
		if parsed.To4() == nil {
			return "IPv6"
		}
		return "未知"
	}

	result := formatRegion(region)
	if strings.TrimSpace(result) == "" {
		if parsed.To4() == nil {
			return "IPv6"
		}
		return "未知"
	}
	return result
}

func buildIPDBCandidates(dbPath string) []string {
	unique := map[string]struct{}{}
	add := func(path string) {
		if path == "" {
			return
		}
		clean := filepath.Clean(path)
		if _, exists := unique[clean]; exists {
			return
		}
		unique[clean] = struct{}{}
	}

	add(dbPath)

	if wd, err := os.Getwd(); err == nil {
		add(filepath.Join(wd, dbPath))
		add(filepath.Join(wd, "resource", "ip2region", "ip2region.xdb"))
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		add(filepath.Join(exeDir, dbPath))
		add(filepath.Join(exeDir, "resource", "ip2region", "ip2region.xdb"))
	}

	result := make([]string, 0, len(unique))
	for path := range unique {
		result = append(result, path)
	}
	return result
}

func normalizeIP(ip string) string {
	v := strings.TrimSpace(ip)
	if v == "" {
		return ""
	}

	if strings.Contains(v, ",") {
		parts := strings.Split(v, ",")
		v = strings.TrimSpace(parts[0])
	}

	if host, _, err := net.SplitHostPort(v); err == nil {
		v = host
	}

	return strings.Trim(v, "[]")
}

// formatRegion 格式化 ip2region 返回的原始地区字符串
// 原始: "中国|0|北京|北京|联通" → "中国 北京"
// 原始: "美国|0|加利福尼亚|洛杉矶|0" → "美国 加利福尼亚 洛杉矶"
func formatRegion(region string) string {
	parts := strings.Split(region, "|")
	if len(parts) < 5 {
		return strings.TrimSpace(region)
	}

	country := regionPart(parts[0])
	area := regionPart(parts[1])
	province := regionPart(parts[2])
	city := regionPart(parts[3])

	// 某些海外IP会把运营商/机构名写到“城市”位，这里做过滤。
	if isLikelyOrgName(city) {
		city = ""
	}

	result := make([]string, 0, 3)
	appendUniqueRegionPart(&result, country)
	appendUniqueRegionPart(&result, province)
	appendUniqueRegionPart(&result, city)

	// 当省市缺失时，回退展示“区域”(parts[1])，例如："中国 香港特别行政区"。
	if len(result) <= 1 {
		appendUniqueRegionPart(&result, area)
	}

	if len(result) == 0 {
		return "未知"
	}
	return strings.Join(result, " ")
}

func regionPart(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || v == "0" {
		return ""
	}
	return v
}

func appendUniqueRegionPart(dst *[]string, v string) {
	if v == "" {
		return
	}
	for _, existed := range *dst {
		if existed == v {
			return
		}
	}
	*dst = append(*dst, v)
}

func isLikelyOrgName(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	lower := strings.ToLower(v)
	keywords := []string{
		"communication", "communications", "telecom", "telecommunications", "network", "networks",
		"datacenter", "hosting", "cloud", "cdn", "isp", "llc", "ltd", "inc", "corp", "company", "co.",
		"通信", "电信", "联通", "移动", "宽带", "网络", "科技", "公司", "数据中心", "云计算", "运营商",
	}
	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}
