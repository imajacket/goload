package serve

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imajacket/goload/example"
)

// Serve
// Serves GoLoad
func Serve(devPort, targetPort int) {
	h := newWsHelper()
	r := chi.NewRouter()

	r.Use(redirect([]string{
		"/update",
		"/ws",
		"/example",
		"/goload.js",
	}, devPort, targetPort))

	r.Get("/update", h.reloadHandler)
	r.Get("/ws", h.wsHandler)
	r.Get("/goload.js", jsHandler(devPort))

	example.InitDevRoutes(r)

	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", devPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
