package httpserver

import (
	"context"
	"fmt"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"
)

func init() {
	ioc.Registry(New, configuration.NewConf, NewRequestLoggerMiddleware)
	ioc.RegistryAtEnd(startAtEnd, New, observability.NewObservability)
}

type Server struct {
	Manager  *fuego.Server
	conf     configuration.Conf
	listener net.Listener
}

func New(conf configuration.Conf, requestLoggerMiddleware RequestLoggerMiddleware) Server {
	s := fuego.NewServer(fuego.WithAddr(":" + conf.PORT))
	server := Server{
		Manager: s,
		conf:    conf,
	}
	fmt.Println("Starting server...")
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Second*2)
		defer shutdownCancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Failed to shutdown:", err)
		}
		cancel()
	}()
	fuego.WithoutLogger()(s)

	fuego.Use(s, requestLoggerMiddleware)
	fuego.Use(s, NewIdempotencyKeyMiddleware())
	server.healthCheck()
	return server
}

func startAtEnd(e Server, obs observability.Observability) error {
	obs.Logger.Info(
		"http server started",
		"port", e.conf.PORT,
		"service", e.conf.PROJECT_NAME,
		"version", e.conf.VERSION,
	)
	return e.Manager.Run()
}

// SetListener sets a custom listener for the server
func (s *Server) SetListener(listener net.Listener) {
	fuego.WithListener(listener)(s.Manager)
}

func (s Server) healthCheck() error {
	h, err := health.New(
		health.WithComponent(health.Component{
			Name:    s.conf.PROJECT_NAME,
			Version: s.conf.VERSION,
		}), health.WithSystemInfo())
	if err != nil {
		return err
	}
	fuego.GetStd(s.Manager,
		"/health",
		h.Handler().ServeHTTP,
		option.Summary("healthCheck"))
	return nil
}

func WrapPostStd(s Server, path string, f func(w http.ResponseWriter, r *http.Request)) {
	fuego.PostStd(s.Manager, path, f)
}

type RequestLoggerMiddleware func(http.Handler) http.Handler

func init() {
	ioc.Registry(NewRequestLoggerMiddleware, observability.NewObservability)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func NewRequestLoggerMiddleware(obs observability.Observability) RequestLoggerMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx, span := obs.Tracer.Start(r.Context(), "http.request")
			defer span.End()

			sw := &statusWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			defer func() {
				fields := []any{
					"method", r.Method,
					"path", r.URL.Path,
					"remote_ip", clientIP(r),
					"status", sw.status,
					"duration_ms", time.Since(start).Milliseconds(),
				}

				if sw.status >= 500 {
					obs.Logger.ErrorContext(ctx, "http request failed", fields...)
				} else {
					obs.Logger.InfoContext(ctx, "http request", fields...)
				}
			}()

			next.ServeHTTP(sw, r.WithContext(ctx))
		})
	}
}

func clientIP(r *http.Request) string {
	// 1️⃣ Forwarded (RFC 7239)
	if fwd := r.Header.Get("Forwarded"); fwd != "" {
		return fwd
	}

	// 2️⃣ X-Forwarded-For (el primero es el cliente original)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 3️⃣ X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// 4️⃣ Fallback
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

type IdempotencyKeyMiddleware func(http.Handler) http.Handler

func NewIdempotencyKeyMiddleware() IdempotencyKeyMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if idempotencyKey := r.Header.Get("Idempotency-Key"); idempotencyKey != "" {
				ctx = sharedcontext.WithIdempotencyKey(ctx, idempotencyKey)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
