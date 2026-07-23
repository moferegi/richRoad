package service

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	clientService "github.com/flipped-aurora/gin-vue-admin/server/service/client"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// VideoSliceService HLS 视频切片服务
type VideoSliceService struct{}

// ffmpegCheckOnce 缓存 FFmpeg 检测结果
var (
	ffmpegAvailable   bool
	ffmpegCheckDone   bool
	ffmpegCheckMu     sync.Mutex
)

// CheckFfmpeg 检测 FFmpeg 是否可用
func (s *VideoSliceService) CheckFfmpeg() (bool, error) {
	ffmpegCheckMu.Lock()
	defer ffmpegCheckMu.Unlock()

	if ffmpegCheckDone {
		return ffmpegAvailable, nil
	}

	ffmpegCheckDone = true
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Warn("FFmpeg 不可用", zap.Error(err), zap.String("output", string(output)))
		ffmpegAvailable = false
		return false, fmt.Errorf("FFmpeg 未安装或不可用，请先安装 FFmpeg: %w", err)
	}

	global.GVA_LOG.Info("FFmpeg 可用", zap.String("version", strings.SplitN(string(output), "\n", 2)[0]))
	ffmpegAvailable = true
	return true, nil
}

// SliceResult HLS切片结果
type SliceResult struct {
	M3u8URL     string  `json:"m3u8Url"`
	Duration    float64 `json:"duration"`
	SegmentCount int    `json:"segmentCount"`
}

// SliceAndUpload 接收视频文件路径，切片并上传至默认云存储
// videoPath: 本地临时视频文件路径
// cloudPrefix: 云存储前缀，如 "hls/{episodeID}"
// globalKey: AES-128加密密钥(hex字符串, 32字符=16字节)，为空则不加密
// episodeID: 单集ID，用于生成密钥token
// 返回 m3u8 完整 URL、时长(秒)和错误
func (s *VideoSliceService) SliceAndUpload(videoPath string, cloudPrefix string, globalKey string, episodeID uint) (*SliceResult, error) {
	// 1. 检查 FFmpeg
	ok, err := s.CheckFfmpeg()
	if !ok {
		return nil, err
	}

	// 2. 在源视频同目录下创建时间戳子目录存放切片（如 D:/hi/2.mp4 → D:/hi/55555/）
	timestamp := time.Now().Format("20060102150405")
	sliceDir := filepath.Join(filepath.Dir(videoPath), timestamp)
	if err := os.MkdirAll(sliceDir, 0755); err != nil {
		return nil, fmt.Errorf("创建切片目录失败: %w", err)
	}
	// 切片文件保留本地，不自动清理

	// 3. 获取视频时长 (ffprobe)
	duration, err := s.probeDuration(videoPath)
	if err != nil {
		global.GVA_LOG.Warn("获取视频时长失败，使用默认值0", zap.Error(err))
		duration = 0
	}

	// 4. 准备加密密钥文件（如果配置了 global-key）
	var keyInfoPath string
	if globalKey = strings.TrimSpace(globalKey); len(globalKey) == 32 {
		keyBytes, hexErr := hex.DecodeString(globalKey)
		if hexErr != nil {
			global.GVA_LOG.Warn("HLS global-key hex解码失败，跳过加密", zap.Error(hexErr))
		} else if len(keyBytes) == 16 {
			// 写入 global.key 文件
			keyFilePath := filepath.Join(sliceDir, "global.key")
			if writeErr := os.WriteFile(keyFilePath, keyBytes, 0644); writeErr != nil {
				global.GVA_LOG.Warn("写入global.key失败，跳过加密", zap.Error(writeErr))
			} else {
				// 生成 JWT token，嵌入密钥URI
				token, tokenErr := s.GenerateHlsKeyToken(episodeID)
				if tokenErr != nil {
					global.GVA_LOG.Warn("生成HLS密钥token失败，跳过加密", zap.Error(tokenErr))
				} else {
					// 构造密钥URI：使用相对路径，前端会根据当前环境动态拼接后端地址
					keyURI := global.GVA_CONFIG.System.RouterPrefix + "/hlsKey?token=" + token
					// key_info 文件格式:
					// 第1行: 密钥URI (播放器通过此URL获取密钥)
					// 第2行: 本地密钥文件路径 (ffmpeg用此加密)
					// 第3行: IV (可选，留空则使用segment sequence)
					keyInfoContent := fmt.Sprintf("%s\n%s\n", keyURI, keyFilePath)
					keyInfoPath = filepath.Join(sliceDir, "key_info.txt")
					if writeErr := os.WriteFile(keyInfoPath, []byte(keyInfoContent), 0644); writeErr != nil {
						global.GVA_LOG.Warn("写入key_info文件失败，跳过加密", zap.Error(writeErr))
						keyInfoPath = ""
					} else {
						global.GVA_LOG.Info("HLS AES-128加密已启用", zap.Uint("episodeId", episodeID))
					}
				}
			}
		}
	}

	// 5. 执行 FFmpeg 切片
	m3u8Path := filepath.Join(sliceDir, "index.m3u8")
	segmentPattern := filepath.Join(sliceDir, "segment_%03d.ts")
	if err := s.runSlice(videoPath, m3u8Path, segmentPattern, keyInfoPath); err != nil {
		return nil, fmt.Errorf("FFmpeg 切片失败: %w", err)
	}

	// 6. 获取默认上传云配置
	cloudSvc := clientService.CloudStorageService{}
	defaultDomain, err := cloudSvc.GetDefaultUploadDomain()
	if err != nil {
		return nil, fmt.Errorf("查询默认上传云配置失败: %w", err)
	}
	if defaultDomain == nil {
		return nil, fmt.Errorf("未找到默认上传云配置，请在 ExternalLinkDomain 中设置一个默认云存储")
	}

	// 7. 获取云存储的公开访问域名
	baseURL := s.resolveCloudBaseURL(*defaultDomain)
	if baseURL == "" {
		// 无法获取云域名时回退到本地方式，由 SignLearningVideoURL 处理
		global.GVA_LOG.Warn("无法获取云存储域名，使用相对路径")
	}

	// 8. 收集并上传所有文件（.ts + .m3u8）
	files, err := os.ReadDir(sliceDir)
	if err != nil {
		return nil, fmt.Errorf("读取临时目录失败: %w", err)
	}

	var segmentCount int
	var uploadErrors []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发上传数
	sem := make(chan struct{}, 4)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
		// 跳过本地辅助文件：不传到CDN
		if fileName == "key_info.txt" || fileName == "enc.key" || fileName == "global.key" {
			continue
		}
		segmentCount++
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			localPath := filepath.Join(sliceDir, fname)
			cloudKey := cloudPrefix + "/" + timestamp + "/" + fname

			data, readErr := os.ReadFile(localPath)
			if readErr != nil {
				mu.Lock()
				uploadErrors = append(uploadErrors, fmt.Sprintf("读取 %s 失败: %v", fname, readErr))
				mu.Unlock()
				return
			}

			contentType := "video/mp2t"
			if strings.HasSuffix(fname, ".m3u8") {
				contentType = "application/vnd.apple.mpegurl"
			}

			_, upErr := cloudSvc.UploadBytesToCloud(*defaultDomain, data, cloudKey, contentType)
			if upErr != nil {
				mu.Lock()
				uploadErrors = append(uploadErrors, fmt.Sprintf("上传 %s 失败: %v", fname, upErr))
				mu.Unlock()
			}
		}(fileName)
	}
	wg.Wait()

	if len(uploadErrors) > 0 {
		global.GVA_LOG.Error("HLS 片段上传存在失败",
			zap.String("prefix", cloudPrefix),
			zap.Strings("errors", uploadErrors))
		return nil, fmt.Errorf("部分文件上传失败: %s", strings.Join(uploadErrors, "; "))
	}

	// 9. 构造 m3u8 完整 URL
	m3u8CloudKey := cloudPrefix + "/" + timestamp + "/index.m3u8"
	m3u8URL := m3u8CloudKey // 默认返回相对路径
	if baseURL != "" {
		m3u8URL = strings.TrimRight(baseURL, "/") + "/" + m3u8CloudKey
	}

	global.GVA_LOG.Info("HLS 切片上传完成",
		zap.String("m3u8", m3u8URL),
		zap.Float64("duration", duration),
		zap.Int("segments", segmentCount))

	return &SliceResult{
		M3u8URL:      m3u8URL,
		Duration:     duration,
		SegmentCount: segmentCount,
	}, nil
}

// HlsKeyClaims HLS密钥token的JWT声明
type HlsKeyClaims struct {
	EpisodeID uint   `json:"episodeId"`
	Purpose   string `json:"purpose"`
	jwt.RegisteredClaims
}

// GenerateHlsKeyToken 生成HLS密钥访问token（嵌入m3u8的#EXT-X-KEY URI中）
func (s *VideoSliceService) GenerateHlsKeyToken(episodeID uint) (string, error) {
	claims := HlsKeyClaims{
		EpisodeID: episodeID,
		Purpose:   "hls_key",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.GVA_CONFIG.JWT.SigningKey))
}

// ValidateHlsKeyToken 验证HLS密钥token
func (s *VideoSliceService) ValidateHlsKeyToken(tokenStr string) (*HlsKeyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &HlsKeyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.GVA_CONFIG.JWT.SigningKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token解析失败: %w", err)
	}
	if claims, ok := token.Claims.(*HlsKeyClaims); ok && token.Valid {
		if claims.Purpose != "hls_key" {
			return nil, fmt.Errorf("非法token用途")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("无效token")
}

// probeDuration 通过 ffprobe 获取视频时长（秒）
func (s *VideoSliceService) probeDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		videoPath,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe 执行失败: %w", err)
	}

	var info struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal(output, &info); err != nil {
		return 0, fmt.Errorf("解析 ffprobe 结果失败: %w", err)
	}

	var duration float64
	if _, err := fmt.Sscanf(info.Format.Duration, "%f", &duration); err != nil {
		return 0, fmt.Errorf("解析时长失败: %w", err)
	}

	return duration, nil
}

// runSlice 执行 FFmpeg HLS 切片
// keyInfoPath: AES-128 key_info 文件路径，为空则不加密
func (s *VideoSliceService) runSlice(inputPath, outputM3u8, segmentPattern, keyInfoPath string) error {
	args := []string{
		"-i", inputPath,
		"-c", "copy",                      // 码流拷贝，不重新编码
		"-bsf:v", "h264_mp4toannexb",       // H.264 AnnexB 格式转换（TS 容器需要）
		"-hls_time", "10",                  // 每段 10 秒
		"-hls_list_size", "0",              // 列出所有段
		"-hls_segment_filename", segmentPattern,
	}
	if keyInfoPath != "" {
		args = append(args, "-hls_key_info_file", keyInfoPath)
	}
	args = append(args, "-f", "hls", outputM3u8)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("FFmpeg 切片命令失败",
			zap.String("input", inputPath),
			zap.String("output", outputM3u8),
			zap.String("stderr", string(output)),
			zap.Error(err))
		return fmt.Errorf("ffmpeg 切片失败: %w\n%s", err, string(output))
	}

	return nil
}

// resolveCloudBaseURL 获取云存储的公开访问基础 URL
func (s *VideoSliceService) resolveCloudBaseURL(domain clientModel.ExternalLinkDomain) string {
	result := strings.TrimSpace(domain.Domain)
	if result == "" {
		result = strings.TrimSpace(domain.BaseURL)
	}
	return strings.TrimRight(result, "/")
}
