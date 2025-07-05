package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"transport-app/app/adapter/in/graphql/graph"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/sharedcontext"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hellofresh/health-go/v5"
	"github.com/vektah/gqlparser/v2/ast"
	"go.opentelemetry.io/otel/baggage"
)

func init() {
	ioc.Registry(New, configuration.NewConf)
	ioc.RegistryAtEnd(startAtEnd, New, graph.NewResolver)
}

type Server struct {
	Manager *fuego.Server
	conf    configuration.Conf
}

// Middleware para forzar headers JSON
func NewJSONMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

func New(conf configuration.Conf) Server {
	s := fuego.NewServer(
		fuego.WithAddr(":"+conf.PORT),
		fuego.WithGlobalMiddlewares(
			injectBaggageMiddleware,
			NewJSONMiddleware(),
		),
	)
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

func startAtEnd(e Server, resolver *graph.Resolver) error {
	fmt.Println(`
████████╗██████╗  █████╗ ███╗   ██╗███████╗██████╗  ██████╗ ██████╗ ████████╗     █████╗ ██████╗ ██████╗ 
╚══██╔══╝██╔══██╗██╔══██╗████╗  ██║██╔════╝██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝    ██╔══██╗██╔══██╗██╔══██╗
   ██║   ██████╔╝███████║██╔██╗ ██║███████╗██████╔╝██║   ██║██████╔╝   ██║       ███████║██████╔╝██████╔╝
   ██║   ██╔══██╗██╔══██║██║╚██╗██║╚════██║██╔═══╝ ██║   ██║██╔══██╗   ██║       ██╔══██║██╔═══╝ ██╔═══╝ 
   ██║   ██║  ██║██║  ██║██║ ╚████║███████║██║     ╚██████╔╝██║  ██║   ██║       ██║  ██║██║     ██║     
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝       ╚═╝  ╚═╝╚═╝     ╚═╝    
   `)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	fuego.PostStd(e.Manager, "/query", srv.ServeHTTP, option.Tags("graphql"))
	fuego.GetStd(e.Manager,
		"/",
		playground.Handler("GraphQL playground", "/query"),
		option.Summary("GraphQL Playground"),
		option.Tags("graphql"))

	return e.Manager.Run()
}

func (s Server) healthCheck() error {
	// Retryable client configurado
	client := retryablehttp.NewClient()
	client.RetryMax = 2
	client.RetryWaitMin = 300 * time.Millisecond
	client.RetryWaitMax = 1 * time.Second
	client.HTTPClient.Timeout = 3 * time.Second

	/*
			Puedes ejecutar este curl aca?
			$body = '{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}'
		Invoke-RestMethod -Uri 'http://localhost:3000/' -Method Post -Body $body -ContentType 'application/json'
	*/

	h, err := health.New(
		health.WithComponent(health.Component{
			Name:    s.conf.PROJECT_NAME,
			Version: s.conf.VERSION,
		}),
		health.WithSystemInfo(),
		func(h *health.Health) error {
			h.Register(health.Config{
				Name: "vroom",
				Check: func(ctx context.Context) error {
					body := `{"jobs":[{"id":1,"location":[-70.6483,-33.4372]}],"vehicles":[{"id":1,"start":[-70.6483,-33.4372]}]}`
					req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:3000/", strings.NewReader(body))
					if err != nil {
						return fmt.Errorf("failed to create request: %w", err)
					}
					req.Header.Set("Content-Type", "application/json")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return fmt.Errorf("failed to make request: %w", err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
					}
					return nil
				},
				Timeout:   3 * time.Second,
				SkipOnErr: false,
			})
			return nil
		})
	if err != nil {
		return err
	}

	fuego.GetStd(s.Manager, "/health", h.Handler().ServeHTTP, option.Summary("healthCheck"))
	return nil
}

func WrapPostStd(s Server, path string, f func(w http.ResponseWriter, r *http.Request)) {
	fuego.PostStd(s.Manager, path, f)
}

func injectBaggageMiddleware(next http.Handler) http.Handler {
	skipPaths := map[string]struct{}{
		"/":                      {},
		"/login":                 {},
		"/register":              {},
		"/health":                {},
		"/favicon.ico":           {},
		"/tenants":               {},
		"/.well-known/jwks.json": {},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, skip := skipPaths[r.URL.Path]; skip || strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/query" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				var body struct {
					OperationName string `json:"operationName"`
				}
				_ = json.Unmarshal(bodyBytes, &body)
				if body.OperationName == "IntrospectionQuery" {
					r.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) // restaura el body original
					next.ServeHTTP(w, r)
					return
				}
				// restaura igualmente para el flujo normal
				r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			}
		}

		orgHeader := r.Header.Get("tenant")
		if orgHeader == "" {
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}

		// Cortamos por el último guion para separar UUID del país
		lastDash := strings.LastIndex(orgHeader, "-")
		if lastDash == -1 || lastDash == len(orgHeader)-1 {
			http.Error(w, "invalid tenant format, expected uuid-country (e.g., a946ac90-...-CL)", http.StatusBadRequest)
			return
		}

		tenantID := orgHeader[:lastDash]
		countryCode := strings.ToUpper(orgHeader[lastDash+1:])

		if _, err := uuid.Parse(tenantID); err != nil {
			http.Error(w, "invalid tenant UUID", http.StatusBadRequest)
			return
		}

		c := countries.ByName(countryCode)
		if c == countries.Unknown {
			http.Error(w, "invalid country code", http.StatusBadRequest)
			return
		}
		country := c.Alpha2()

		members := make([]baggage.Member, 0, 4)

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

		if v := r.Header.Get("channel"); v != "" {
			m, _ := baggage.NewMember(sharedcontext.BaggageChannel, v)
			members = append(members, m)
		}

		bag, _ := baggage.New(members...)
		ctx := baggage.ContextWithBaggage(r.Context(), bag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
