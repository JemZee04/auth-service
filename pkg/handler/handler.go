package handler

import (
	"context"
	"net/http"
	"regexp"

	"strings"
)

type Handler struct {
}

//func (h *Handler) InitRoutes() *http {
//	mux := http.NewServeMux()
//
//	mux.HandleFunc("/")
//
//	mux.Ge
//}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

//var routes = []route{
//	newRoute("GET", "/", signUp),
//	newRoute("POST", "/sign-in", refresh),
//	newRoute("GET", "/fec", home),
//}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	var routes = []route{
		newRoute("GET", "/", h.signUp),
		newRoute("POST", "/sign-in", h.refresh),
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
