package serve

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"
)

// redirect
// Middleware that handles requests to GoLoad server
// GoLoad routes continue on.
// Non-GET requests are reverse proxied to target server
// Non-HTML GET requests are sent as-is
// HTML GET requests are modified with goload.js and sent
func redirect(routes []string, devPort, targetPort int) func(next http.Handler) http.Handler {
	targetHost := fmt.Sprintf("http://localhost:%d", targetPort)
	htmlMod := fmt.Sprintf("<script src=\"http://localhost:%v/goload.js\"></script></body>", devPort)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Checks if request is any of GoLoad's internal routes
			requestPath := r.URL.String()
			if slices.Contains(routes, requestPath) {
				next.ServeHTTP(w, r)
				return
			}

			// Proxy any non-GET requests
			if r.Method != http.MethodGet {
				proxy(targetHost, w, r)
				return
			}

			// First request the GET content
			getTarget := fmt.Sprintf("%s%s", targetHost, r.URL.Path)
			req, err := http.NewRequestWithContext(r.Context(), "GET", getTarget, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			for k, vv := range resp.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}

			// If GET does not return HTML return response as-is
			if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
				log.Println("returning", r.URL.String(), "unmodified")
				_, _ = io.Copy(w, resp.Body)
				return
			}

			// Modify HTML to include GoLoad and return
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			modifiedHTML := bytes.Replace(body, []byte("</body>"), []byte(htmlMod), 1)
			_, _ = w.Write(modifiedHTML)
			log.Println("returning", r.URL.String(), "modified")
		})
	}
}

// proxy
// Reverse proxy
func proxy(targetHost string, w http.ResponseWriter, r *http.Request) {
	log.Println("proxying", r.URL.String())

	targetUrl, _ := url.Parse(targetHost)

	p := httputil.NewSingleHostReverseProxy(targetUrl)
	r.URL.Host = targetUrl.Host
	r.URL.Scheme = targetUrl.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetUrl.Host

	p.ServeHTTP(w, r)
}
