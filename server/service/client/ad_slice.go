package client

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// AdSliceService 广告视频切片上传服务
type AdSliceService struct{}

// SliceResult 切片结果
type AdSliceResult struct {
	MediaUrl     string  `json:"mediaUrl"`
	Duration     float64 `json:"duration"`
	SegmentCount int     `json:"segmentCount"`
}

// SliceAndUploadAdVideo 接收视频文件，切片并上传至默认云存储
// 参照 englishLearningVideo 的 SliceAndUpload 方法，但不加密
func (s *AdSliceService) SliceAndUploadAdVideo(fileHeader *multipart.FileHeader, cloudPrefix string) (*AdSliceResult, error) {
	// 1. 保存上传文件到临时目录
	uploadsDir := filepath.Join(global.GVA_CONFIG.Local.StorePath, "uploads", "ad_video")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	ext := filepath.Ext(fileHeader.Filename)
	tempFileName := fmt.Sprintf("ad_%d%s", time.Now().UnixNano(), ext)
	tempFilePath := filepath.Join(uploadsDir, tempFileName)

	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(tempFilePath)
	if err != nil {
		return nil, fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer dst.Close()

	if _, err = dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 2. 检查 FFmpeg
	if err := s.CheckFfmpeg(); err != nil {
		return nil, err
	}

	// 3. 创建切片目录
	timestamp := time.Now().Format("20060102150405")
	sliceDir := filepath.Join(filepath.Dir(tempFilePath), "slice_"+timestamp)
	if err := os.MkdirAll(sliceDir, 0755); err != nil {
		return nil, fmt.Errorf("创建切片目录失败: %w", err)
	}

	// 4. 获取视频时长
	duration, err := s.probeDuration(tempFilePath)
	if err != nil {
		global.GVA_LOG.Warn("获取广告视频时长失败，使用默认值0", zap.Error(err))
		duration = 0
	}

	// 5. 执行 FFmpeg 切片（不加密）
	m3u8Path := filepath.Join(sliceDir, "index.m3u8")
	segmentPattern := filepath.Join(sliceDir, "segment_%03d.ts")
	if err := s.runSlice(tempFilePath, m3u8Path, segmentPattern); err != nil {
		return nil, fmt.Errorf("FFmpeg 切片失败: %w", err)
	}

	// 6. 获取默认上传云配置
	cloudSvc := CloudStorageService{}
	defaultDomain, err := cloudSvc.GetDefaultUploadDomain()
	if err != nil {
		return nil, fmt.Errorf("查询默认上传云配置失败: %w", err)
	}
	if defaultDomain == nil {
		return nil, fmt.Errorf("未找到默认上传云配置，请在 ExternalLinkDomain 中设置一个默认云存储")
	}

	// 7. 收集并上传所有文件
	files, err := os.ReadDir(sliceDir)
	if err != nil {
		return nil, fmt.Errorf("读取切片目录失败: %w", err)
	}

	var segmentCount int
	var uploadErrors []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
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
		global.GVA_LOG.Error("广告视频切片上传存在失败",
			zap.String("prefix", cloudPrefix),
			zap.Strings("errors", uploadErrors))
		return nil, fmt.Errorf("部分文件上传失败: %s", strings.Join(uploadErrors, "; "))
	}

	// 8. 构造 m3u8 相对路径
	m3u8URL := cloudPrefix + "/" + timestamp + "/index.m3u8"

	global.GVA_LOG.Info("广告视频切片上传完成",
		zap.String("m3u8", m3u8URL),
		zap.Float64("duration", duration),
		zap.Int("segments", segmentCount))

	return &AdSliceResult{
		MediaUrl:     m3u8URL,
		Duration:     duration,
		SegmentCount: segmentCount,
	}, nil
}

// CheckFfmpeg 检查 FFmpeg 是否可用
func (s *AdSliceService) CheckFfmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Warn("FFmpeg 不可用", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("FFmpeg 未安装或不可用，请先安装 FFmpeg: %w", err)
	}
	return nil
}

func (s *AdSliceService) probeDuration(videoPath string) (float64, error) {
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

func (s *AdSliceService) runSlice(inputPath, outputM3u8, segmentPattern string) error {
	args := []string{
		"-i", inputPath,
		"-c", "copy",
		"-bsf:v", "h264_mp4toannexb",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentPattern,
		"-f", "hls", outputM3u8,
	}

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
