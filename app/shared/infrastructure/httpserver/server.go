package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/hellofresh/health-go/v5"
	"go.opentelemetry.io/otel/baggage"
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
	s := fuego.NewServer(
		fuego.WithAddr(":"+conf.PORT),
		fuego.WithGlobalMiddlewares(injectBaggageMiddleware))
	server := Server{
		Manager: s,
		conf:    conf,
	}
	server.healthCheck()
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

func injectBaggageMiddleware(next http.Handler) http.Handler {
	skipPaths := map[string]struct{}{
		"/login":         {},
		"/register":      {},
		"/health":        {},
		"/organizations": {},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, skip := skipPaths[r.URL.Path]; skip {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}
		orgHeader := r.Header.Get("organization")
		if orgHeader == "" {
			http.Error(w, "missing organization header", http.StatusBadRequest)
			return
		}

		parts := strings.SplitN(orgHeader, "-", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.Error(w, "invalid organization format, expected tenant-country (e.g., 1-CL)", http.StatusBadRequest)
			return
		}

		tenantID := parts[0]
		c := countries.ByName(strings.ToUpper(parts[1]))
		if c == countries.Unknown {
			http.Error(w, "invalid country code", http.StatusBadRequest)
			return
		}
		country := c.Alpha2()

		members := make([]baggage.Member, 0, 3)

		mTenant, _ := baggage.NewMember(sharedcontext.BaggageTenantID, tenantID)
		mCountry, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, country)
		members = append(members, mTenant, mCountry)

		if v := r.Header.Get("consumer"); v != "" {
			m, _ := baggage.NewMember(sharedcontext.BaggageConsumer, v)
			members = append(members, m)
		}

		if v := r.Header.Get("commerce"); v != "" {
			m, _ := baggage.NewMember(sharedcontext.BaggageCommerce, v)
			members = append(members, m)
		}

		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(r.Context(), bag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
