package plugin_cond_middleware

import (
  "context"
  "net/http"
  "strings"
)

type Config struct {
  SourceRange []string `json:"ipRange"`
  PathPrefix  string   `json:"pathPrefix"`
}

func CreateConfig() *Config {
  return &Config{}
}

type CondMiddleware struct {
  next        http.Handler
  sourceRange []string
  pathPrefix  string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
  return &CondMiddleware{
    next:        next,
    sourceRange: config.SourceRange,
    pathPrefix:  config.PathPrefix,
  }, nil
}

func (m *CondMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  if strings.HasPrefix(req.URL.Path, m.pathPrefix) {
    clientIP := req.RemoteAddr
    allowed := false
    for _, ip := range m.sourceRange {
      if strings.HasPrefix(clientIP, ip) {
        allowed = true
        break
      }
    }
    if !allowed {
      http.Error(rw, "Forbidden", http.StatusForbidden)
      return
    }
  }
  m.next.ServeHTTP(rw, req)
}