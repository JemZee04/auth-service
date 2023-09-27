package handler

import (
	"auth-service/pkg/service"
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	var routes = []route{
		newRoute("POST", "/sign-up", h.signUp),
		newRoute("POST", "/sign-in", h.signIn),
		newRoute("POST", "/refresh", h.refresh),
	}
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}
