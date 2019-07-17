package sentrygin

import (
	"context"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

const valuesKey = "sentry"

type handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

type Options struct {
	// Repanic configures whether Sentry should repanic after recovery, in most cases it should be set to true,
	// as gin.Default includes it's own Recovery middleware what handles http responses.
	Repanic bool
	// WaitForDelivery configures whether you want to block the request before moving forward with the response.
	// Because Gin's default `Recovery` handler doesn't restart the application,
	// it's safe to either skip this option or set it to `false`.
	WaitForDelivery bool
	// Timeout for the event delivery requests.
	Timeout time.Duration
}

// New returns a function that satisfies gin.HandlerFunc interface
// It can be used with Use() methods.
func New(options Options) gin.HandlerFunc {
	handler := handler{
		repanic:         false,
		timeout:         time.Second * 2,
		waitForDelivery: false,
	}

	if options.Repanic {
		handler.repanic = true
	}

	if options.WaitForDelivery {
		handler.waitForDelivery = true
	}

	return handler.handle
}

func (h *handler) handle(ctx *gin.Context) {
	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetRequest(sentry.Request{}.FromHTTPRequest(ctx.Request))
	ctx.Set(valuesKey, hub)
	defer h.recoverWithSentry(hub, ctx.Request)
	ctx.Next()
}

func (h *handler) recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
		if eventID != nil && h.waitForDelivery {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}

// GetHubFromContext retrieves attached *sentry.Hub instance from gin.Context.
func GetHubFromContext(ctx *gin.Context) *sentry.Hub {
	if hub, ok := ctx.Get(valuesKey); ok {
		if hub, ok := hub.(*sentry.Hub); ok {
			return hub
		}
	}
	return nil
}
