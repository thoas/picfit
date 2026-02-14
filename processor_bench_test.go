package picfit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/store"
	"github.com/thoas/picfit/tests"
)

func BenchmarkProcessor_ProcessContext(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")

	// Prepare common data
	params := map[string]any{
		"url": u,
		"w":   50,
		"h":   50,
		"op":  "resize",
	}
	key := "test-key"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ProcessContext_WithCache(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	cfg.KVStore = &store.Config{
		Type: "cache",
	}
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url": u,
		"w":   50,
		"h":   50,
		"op":  "resize",
	}
	key := "test-key-cached"

	// Warm up cache
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	c.Request = req
	c.Set("key", key)
	c.Set("parameters", params)
	c.Set("url", u)
	_, _ = processor.ProcessContext(c)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ProcessContext_Thumbnail(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url": u,
		"w":   50,
		"h":   50,
		"op":  "thumbnail",
	}
	key := "test-key-thumbnail"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ProcessContext_Rotate(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url": u,
		"op":  "rotate",
		"deg": 90,
	}
	key := "test-key-rotate"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ProcessContext_Flip(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url": u,
		"op":  "flip",
		"pos": "h",
	}
	key := "test-key-flip"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ProcessContext_Fit(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url": u,
		"w":   50,
		"h":   50,
		"op":  "fit",
	}
	key := "test-key-fit"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessor_ShardFilename(b *testing.B) {
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Shard.Width = 1
	cfg.Shard.Depth = 2
	cfg.Debug = false
	cfg.Logger.Level = "error"

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	filename := "018e7d45412eb26c7cf06c42c8fff633"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = processor.ShardFilename(filename)
	}
}

func BenchmarkProcessor_FileExists(b *testing.B) {
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = processor.FileExists(ctx, "avatar.png")
	}
}

func BenchmarkProcessor_ProcessContext_Blur(b *testing.B) {
	ts := tests.NewImageServer()
	defer ts.Close()

	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Debug = false
	cfg.Logger.Level = "error"
	gin.SetMode(gin.ReleaseMode)

	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		b.Fatal(err)
	}

	u, _ := url.Parse(ts.URL + "/original.jpg")
	params := map[string]any{
		"url":    u,
		"op":     "effect",
		"filter": "blur",
	}
	key := "test-key-rotate"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
		c.Request = req

		c.Set("key", key)
		c.Set("parameters", params)
		c.Set("url", u)

		_, err := processor.ProcessContext(c)
		if err != nil {
			b.Fatal(err)
		}
	}
}
