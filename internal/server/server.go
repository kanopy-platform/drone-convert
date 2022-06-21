package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/kanopy-platform/drone-convert/internal/plugin"
	"github.com/kanopy-platform/go-http-middleware/logging"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	router *http.ServeMux
}

func New(secret string) http.Handler {
	s := &Server{router: http.NewServeMux()}
	logrusMiddleware := logging.NewLogrus()

	plugins := []converter.Plugin{
		plugin.New(),
	}

	s.router.Handle("/", s.handlePlugins(secret, plugins...))
	s.router.HandleFunc("/healthz", s.handleHealthz())

	return logrusMiddleware.Middleware(s.router)
}

// handlePlugins accepts multiple conversion plugins and executes
// their handlers in the passed in order
func (s *Server) handlePlugins(secret string, plugins ...converter.Plugin) http.HandlerFunc {
	logger := log.StandardLogger()
	handlers := []http.Handler{}

	for _, plugin := range plugins {
		handlers = append(handlers, converter.Handler(plugin, secret, logger))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler.ServeHTTP(w, r)
		}
	}
}

func (s *Server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"status": "ok",
		}

		bytes, err := json.Marshal(status)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(bytes))
	}
}
