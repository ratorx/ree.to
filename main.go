package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type config struct {
	httpPort  uint16
	httpsPort uint16
	certPath  string
	keyPath   string
}

var cfg config

func getPort(env string, def uint16) uint16 {
	if value, success := os.LookupEnv(env); success {
		if port, err := strconv.Atoi(value); err == nil {
			return uint16(port)
		}
	}

	return def
}

func getConfigValue(env string, def string) string {
	if value, success := os.LookupEnv(env); success {
		return value
	}

	return def
}

func init() {
	cfg.httpPort = getPort("CV_HTTP_PORT", 8080)
	cfg.httpsPort = getPort("CV_HTTPS_PORT", 8443)
	cfg.certPath = getConfigValue("CV_CERT_PATH", "")
	cfg.keyPath = getConfigValue("CV_KEY_PATH", "")
}

func main() {
	router := httprouter.New()

	// File Server
	// Workaround to allow other paths
	router.ServeFiles("/css/*filepath", http.Dir("public/css"))
	router.ServeFiles("/img/*filepath", http.Dir("public/img"))
	router.ServeFiles("/js/*filepath", http.Dir("public/js"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})

	// Config init
	router.GET("/init", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/init")
	})

	// Handle 404
	router.NotFound = func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/404.html")
	}

	if cfg.certPath != "" && cfg.keyPath != "" {
		go func() {
			fmt.Println(http.ListenAndServeTLS(fmt.Sprintf(":%v", cfg.httpsPort), cfg.certPath, cfg.keyPath, router))
		}()
	}
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%v", cfg.httpPort), router))
}
