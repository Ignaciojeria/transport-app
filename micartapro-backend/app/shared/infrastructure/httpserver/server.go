package httpserver

import (
	"context"
	"fmt"
	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"
	"github.com/rs/cors"
)

func init() {
	ioc.Register(New)
	ioc.RegisterAtEnd(startAtEnd)
}

type Server struct {
	Manager  *fuego.Server
	conf     configuration.Conf
	listener net.Listener
}

func New(conf configuration.Conf, requestLoggerMiddleware RequestLoggerMiddleware) Server {
	// CORS debe ir antes del logger para manejar preflight requests
	corsMiddleware := NewCORSMiddleware()

	s := fuego.NewServer(
		fuego.WithAddr(":"+conf.PORT),
		fuego.WithGlobalMiddlewares(
			corsMiddleware,
			requestLoggerMiddleware,
		),
	)
	// Timeout de 10 min para pruebas (evitar que corte solicitudes largas como generate speech con escenas)
	s.ReadTimeout = 30 * time.Minute
	s.WriteTimeout = 30 * time.Minute
	s.ReadHeaderTimeout = 30 * time.Minute
	s.IdleTimeout = 30 * time.Minute
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
	ioc.Register(NewRequestLoggerMiddleware)
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

type CORSMiddleware func(http.Handler) http.Handler

func NewCORSMiddleware() CORSMiddleware {
	// Configurar CORS con la librería rs/cors
	// En desarrollo, permitir todos los orígenes usando AllowOriginFunc
	c := cors.New(cors.Options{
		// Permitir cualquier origen en desarrollo
		AllowOriginFunc: func(origin string) bool {
			// Permitir todos los orígenes en desarrollo
			return true
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Idempotency-Key",
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           3600,
		// Habilitar debug para ver qué está pasando
		Debug: false,
	})

	return c.Handler
}
