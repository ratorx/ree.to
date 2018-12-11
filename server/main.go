package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"flag"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
)

type config struct {
	httpPort  uint
	httpsPort uint
	certPath  string
	keyPath   string
}

var cfg config

func getPort(env string, def uint) uint {
	if value, success := os.LookupEnv(env); success {
		if port, err := strconv.Atoi(value); err == nil {
			return uint(port)
		}
	}

	return def
}

func getConfigValue(env string, def string) string {
	fmt.Println(strings.ToLower(env))
	if value, success := os.LookupEnv(env); success {
		return value
	}

	return def
}

func init() {
	flag.UintVar(&cfg.httpPort, "http.port", getPort("CV_HTTP_PORT", 8080), "HTTP port to listen on")
	flag.UintVar(&cfg.httpsPort, "https.port", getPort("CV_HTTPS_PORT", 8443), "HTTPS port to listen on")
	flag.StringVar(&cfg.certPath, "https.cert", getConfigValue("CV_CERT_PATH", ""), "Path to SSL fullchain certificate")
	flag.StringVar(&cfg.keyPath, "https.key", getConfigValue("CV_KEY_PATH", ""), "Path to SSL private key")

	flag.Parse()
}

func setupRoutes(router *httprouter.Router) {
	// Website
	// Workaround to allow other routes
	router.ServeFiles("/css/*filepath", http.Dir("public/css"))
	router.ServeFiles("/img/*filepath", http.Dir("public/img"))
	router.ServeFiles("/js/*filepath", http.Dir("public/js"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})

	// 404
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Warning("404 on path %s", r.URL.EscapedPath())
		http.ServeFile(w, r, "public/404.html")
	})

	// Utilities

	// Config init script
	router.GET("/init", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/init")
	})
}

func main() {
	router := httprouter.New()
	setupRoutes(router)

	if cfg.certPath != "" && cfg.keyPath != "" {
		glog.Infof("Listening for HTTPS on %v", cfg.httpsPort)
		go func() {
			glog.Error(http.ListenAndServeTLS(fmt.Sprintf(":%v", cfg.httpsPort), cfg.certPath, cfg.keyPath, router))
		}()
	}
	glog.Infof("Listening for HTTP on %v", cfg.httpPort)
	glog.Error(http.ListenAndServe(fmt.Sprintf(":%v", cfg.httpPort), router))
}
