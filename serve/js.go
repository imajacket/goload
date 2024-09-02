package serve

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// jsHandler
// Returns GoLoad's javascript file.
func jsHandler(port int) func(http.ResponseWriter, *http.Request) {
	goload, err := os.ReadFile("goload.js")
	if err != nil {
		panic(err)
	}

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	out, err := m.String("text/javascript", string(goload))
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.New("goload").Parse(out))

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		Port string
	}{
		Port: strconv.Itoa(port),
	})
	if err != nil {
		panic(err)
	}

	minifiedJs := buf.Bytes()

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write(minifiedJs)
	}
}
