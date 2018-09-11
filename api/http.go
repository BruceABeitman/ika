package api

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
)

// BuildHTTPServer builds the HTTP server serving the Ika
// API
func BuildHTTPServer(service core.Service) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/proxy", getProxy(service))
	mux.Handle("/proxy/meta", updateProxyMeta(service))
	mux.Handle("/proxy/refresh", proxyRefresh(service))
	return &http.Server{Addr: ":4242", Handler: BasicAuth(mux)}
}

// BasicAuth ...
func BasicAuth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(user)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(pass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="proxyAuth"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func getProxy(service core.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queueID := core.QueueID{
			Channel: r.URL.Query().Get("channel"),
			Domain:  r.URL.Query().Get("domain"),
		}
		proxy, err := service.GetProxy(queueID)
		if err != nil {
			internalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(proxy); err != nil {
			internalServerError(w, err)
		}
	})
}

func updateProxyMeta(service core.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queueID := core.QueueID{
			Channel: r.URL.Query().Get("channel"),
			Domain:  r.URL.Query().Get("domain"),
		}
		proxyMeta := core.ProxyMeta{
			Addr:  r.URL.Query().Get("addr"),
			Error: r.URL.Query().Get("error"),
		}

		if err := service.UpdateProxyMeta(queueID, proxyMeta); err != nil {
			internalServerError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func proxyRefresh(service core.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := service.RefreshProxies("proxy:master")
		if err != nil {
			internalServerError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
