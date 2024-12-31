package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"
)

func init() {
	ioc.Registry(New, configuration.NewConf)
	ioc.RegistryAtEnd(startAtEnd, New)
}

type Server struct {
	Manager *fuego.Server
	conf    configuration.Conf
}

func New(conf configuration.Conf) Server {
	s := fuego.NewServer(fuego.WithAddr(":" + conf.PORT))
	server := Server{
		Manager: s,
		conf:    conf,
	}
	server.healthCheck()
	fuego.Use(s, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			countryCode := r.Header.Get("country")
			country := countries.ByName(countryCode)
			if country == countries.Unknown {
				// Crear una instancia de BadRequestError
				err := fuego.BadRequestError{
					Err: errors.New("invalid country code"),
				}

				// Escribir la respuesta HTTP
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(err.StatusCode())
				fmt.Fprintf(w, `{"type": "%s", "title": "%s", "status": %d, "detail": "%s"}`,
					err.Type, err.Title, err.StatusCode(), err)
				return
			}

			// Continuar con el flujo si no hay error
			next.ServeHTTP(w, r)
		})
	})

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
	return server
}

func startAtEnd(e Server) error {
	return e.Manager.Run()
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
