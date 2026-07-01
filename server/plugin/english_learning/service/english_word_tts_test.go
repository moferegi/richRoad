package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	clientModel "github.com/flipped-aurora/gin-vue-admin/server/model/client"
	requestModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/english_learning/model/request"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupEnglishWordTTSTestDB(t *testing.T) {
	t.Helper()

	originalDB := global.GVA_DB
	t.Cleanup(func() {
		global.GVA_DB = originalDB
	})

	db, err := gorm.Open(sqlite.Open("file:english_word_tts_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}

	if err = db.AutoMigrate(&clientModel.SysConfig{}); err != nil {
		t.Fatalf("automigrate failed: %v", err)
	}

	global.GVA_DB = db
}

func upsertTTSTestConfig(t *testing.T, key, value string) {
	t.Helper()

	var cfg clientModel.SysConfig
	err := global.GVA_DB.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		cfg = clientModel.SysConfig{
			ConfigKey:   key,
			ConfigValue: value,
			ConfigName:  key,
			ConfigGroup: "english_learning",
		}
		if createErr := global.GVA_DB.Create(&cfg).Error; createErr != nil {
			t.Fatalf("create config failed: %v", createErr)
		}
		return
	}

	if updateErr := global.GVA_DB.Model(&clientModel.SysConfig{}).Where("id = ?", cfg.ID).Update("config_value", value).Error; updateErr != nil {
		t.Fatalf("update config failed: %v", updateErr)
	}
}

func TestExtractTTSAudioURL(t *testing.T) {
	testCases := []struct {
		name     string
		payload  string
		expected string
	}{
		{
			name:     "top-level url",
			payload:  `{"url":"https://example.com/top.mp3"}`,
			expected: "https://example.com/top.mp3",
		},
		{
			name:     "nested data audioUrl",
			payload:  `{"data":{"audioUrl":"https://example.com/nested.mp3"}}`,
			expected: "https://example.com/nested.mp3",
		},
		{
			name:     "raw url text",
			payload:  `https://example.com/raw.mp3`,
			expected: "https://example.com/raw.mp3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractTTSAudioURL([]byte(tc.payload))
			if actual != tc.expected {
				t.Fatalf("expected=%s got=%s", tc.expected, actual)
			}
		})
	}
}

func TestRequestTTSAudioURL_Success(t *testing.T) {
	setupEnglishWordTTSTestDB(t)

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method POST, got=%s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"audioUrl":"https://example.com/generated.mp3"}}`))
	}))
	defer ttsServer.Close()

	upsertTTSTestConfig(t, "learning_tts_enabled", "true")
	upsertTTSTestConfig(t, "learning_tts_provider_url", ttsServer.URL)
	upsertTTSTestConfig(t, "learning_tts_timeout_ms", "3000")

	svc := EnglishWordService{}
	audioURL, err := svc.requestTTSAudioURL("destiny", "en-US-JennyNeural")
	if err != nil {
		t.Fatalf("request tts failed: %v", err)
	}
	if audioURL != "https://example.com/generated.mp3" {
		t.Fatalf("unexpected audioURL: %s", audioURL)
	}
}

func TestRequestTTSAudioURL_Disabled(t *testing.T) {
	setupEnglishWordTTSTestDB(t)
	upsertTTSTestConfig(t, "learning_tts_enabled", "false")

	svc := EnglishWordService{}
	_, err := svc.requestTTSAudioURL("destiny", "en-US-JennyNeural")
	if !errors.Is(err, errLearningTTSDisabled) {
		t.Fatalf("expected errLearningTTSDisabled, got=%v", err)
	}
}

func TestRequestTTSAudioURL_HTTPStatusError(t *testing.T) {
	setupEnglishWordTTSTestDB(t)

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"msg":"upstream error"}`))
	}))
	defer ttsServer.Close()

	upsertTTSTestConfig(t, "learning_tts_enabled", "true")
	upsertTTSTestConfig(t, "learning_tts_provider_url", ttsServer.URL)

	svc := EnglishWordService{}
	_, err := svc.requestTTSAudioURL("destiny", "en-US-JennyNeural")
	if err == nil {
		t.Fatal("expected status error, got nil")
	}
	if !strings.Contains(err.Error(), "TTS服务响应异常") {
		t.Fatalf("expected status error message, got=%v", err)
	}
}

func TestRequestTTSAudioURL_MissingAudioURL(t *testing.T) {
	setupEnglishWordTTSTestDB(t)

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"abc"}}`))
	}))
	defer ttsServer.Close()

	upsertTTSTestConfig(t, "learning_tts_enabled", "true")
	upsertTTSTestConfig(t, "learning_tts_provider_url", ttsServer.URL)

	svc := EnglishWordService{}
	_, err := svc.requestTTSAudioURL("destiny", "en-US-JennyNeural")
	if !errors.Is(err, errLearningTTSMissingAudioURL) {
		t.Fatalf("expected errLearningTTSMissingAudioURL, got=%v", err)
	}
}

func TestPreflightTTS_Success(t *testing.T) {
	setupEnglishWordTTSTestDB(t)

	ttsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"audioUrl":"https://example.com/preflight.mp3"}}`))
	}))
	defer ttsServer.Close()

	upsertTTSTestConfig(t, "learning_tts_enabled", "true")
	upsertTTSTestConfig(t, "learning_tts_provider_url", ttsServer.URL)
	upsertTTSTestConfig(t, "learning_tts_timeout_ms", "3000")

	svc := EnglishWordService{}
	result, err := svc.PreflightTTS(requestModel.TTSPreflightReq{Word: "future"})
	if err != nil {
		t.Fatalf("expected preflight success, got=%v", err)
	}
	if result.Word != "future" {
		t.Fatalf("expected word=future, got=%s", result.Word)
	}
	if result.AudioUS == "" || result.AudioUK == "" {
		t.Fatalf("expected both audio urls, got us=%s uk=%s", result.AudioUS, result.AudioUK)
	}
}

func TestPreflightTTS_BothDisabled(t *testing.T) {
	setupEnglishWordTTSTestDB(t)

	checkUS := false
	checkUK := false
	svc := EnglishWordService{}
	_, err := svc.PreflightTTS(requestModel.TTSPreflightReq{CheckUS: &checkUS, CheckUK: &checkUK})
	if err == nil {
		t.Fatal("expected preflight error when both checks are disabled")
	}
	if !strings.Contains(err.Error(), "请至少选择一种发音进行预检") {
		t.Fatalf("unexpected error: %v", err)
	}
}
