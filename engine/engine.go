package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/engine/backend"
	"github.com/thoas/picfit/engine/config"
	"github.com/thoas/picfit/image"
	loggerpkg "github.com/thoas/picfit/logger"
)

type Engine struct {
	DefaultFormat  string
	DefaultQuality int
	Format         string
	backends       []*backendWrapper
	logger         *slog.Logger
}

type backendWrapper struct {
	backend   backend.Backend
	mimetypes []string
	weight    int
}

// New initializes an Engine
func New(cfg config.Config, logger *slog.Logger) *Engine {
	var b []*backendWrapper

	if cfg.Backends == nil {
		b = append(b, &backendWrapper{
			backend:   &backend.GoImage{},
			mimetypes: MimeTypes,
		})
	} else {
		if cfg.Backends.Gifsicle != nil {
			path := cfg.Backends.Gifsicle.Path
			if path == "" {
				path = "gifsicle"
			}

			if _, err := exec.LookPath(path); err == nil {
				b = append(b, &backendWrapper{
					backend:   &backend.Gifsicle{Path: path},
					mimetypes: cfg.Backends.Gifsicle.Mimetypes,
					weight:    cfg.Backends.Gifsicle.Weight,
				})
			}
		}
		if cfg.Backends.GoImage != nil {
			b = append(b, &backendWrapper{
				backend:   &backend.GoImage{},
				mimetypes: cfg.Backends.GoImage.Mimetypes,
				weight:    cfg.Backends.GoImage.Weight,
			})
		}
	}

	sort.Slice(b, func(i, j int) bool {
		return b[i].weight < b[j].weight
	})

	quality := config.DefaultQuality
	if cfg.Quality != 0 {
		quality = cfg.Quality
	}

	return &Engine{
		DefaultFormat:  cfg.DefaultFormat,
		DefaultQuality: quality,
		Format:         cfg.Format,
		backends:       b,
		logger:         logger,
	}
}

func (e Engine) String() string {
	backendNames := []string{}
	for _, backend := range e.backends {
		backendNames = append(backendNames, backend.backend.String())
	}

	return strings.Join(backendNames, " ")
}

func (e Engine) Transform(ctx context.Context, dst io.Writer, output *image.ImageFile, operations []EngineOperation) (*image.ImageFile, error) {
	var (
		err    error
		source = output.Stream
		start  = time.Now()
	)

	ct := output.ContentType()
	for i := range operations {
		isLast := i == len(operations)-1
		var target io.Writer

		// swith writer target
		// on last operation we write on dst
		// else we use a temp buffer
		if isLast {
			target = dst
		} else {
			target = &bytes.Buffer{}
		}
		output.Stream = source

		for j := range e.backends {
			if !slices.Contains(e.backends[j].mimetypes, ct) {
				continue
			}

			defer func() {
				loggerpkg.WithMemStats(e.logger).InfoContext(ctx, "Engine handled image",
					slog.String("image", output.Filepath),
					slog.String("backend", e.backends[j].backend.String()),
					slog.String("operation", operations[i].Operation.String()),
					slog.String("options", operations[i].Options.String()),
					slog.String("duration", time.Now().Sub(start).String()),
				)
			}()

			err = operate(ctx, target, e.backends[j].backend, output, operations[i].Operation, operations[i].Options)
			if err == nil {
				break
			}

			if !errors.Is(err, backend.MethodNotImplementedError) {
				return nil, err
			}
		}
		// is not last operations so we repass target to new source stream
		if !isLast {
			buf := target.(*bytes.Buffer)
			source = io.NopCloser(bytes.NewReader(buf.Bytes()))
		}
	}

	return output, err
}

func operate(ctx context.Context, dst io.Writer, b backend.Backend, img *image.ImageFile, operation Operation, options *backend.Options) error {
	switch operation {
	case Noop:
		_, err := io.Copy(dst, img.Stream)
		return err
	case Flip:
		return b.Flip(ctx, dst, img, options)
	case Rotate:
		return b.Rotate(ctx, dst, img, options)
	case Resize:
		return b.Resize(ctx, dst, img, options)
	case Thumbnail:
		return b.Thumbnail(ctx, dst, img, options)
	case Fit:
		return b.Fit(ctx, dst, img, options)
	case Flat:
		return b.Flat(ctx, dst, img, options)
	case Effect:
		return b.Effect(ctx, dst, img, options)
	default:
		return fmt.Errorf("operation not found for %s", operation)
	}
}
