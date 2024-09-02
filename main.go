package main

import (
	"flag"

	"github.com/imajacket/goload/serve"
)

const DEFAULT_PORT = 3000
const DEFAULT_TARGET_PORT = 8888

func main() {
	devPort := flag.Int("dev_port", DEFAULT_PORT, "Port to run dev server on. Defaults to 3000.")
	targetPort := flag.Int("target_port", DEFAULT_TARGET_PORT, "Port that target server is running on.")
	flag.Parse()

	serve.Serve(*devPort, *targetPort)
}
