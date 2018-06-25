package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"time"

	"github.com/v-zhidu/orb/logging"
	context "golang.org/x/net/context"
)

type contextKey string

const (
	//HeaderKey context key
	HeaderKey = contextKey("headers")
)

type ApiHandler interface {
	Serve(context.Context, *http.Request) (interface{}, int)
}

type ApiHandlerFunc func(context.Context, *http.Request) (interface{}, int)

func (f ApiHandlerFunc) Serve(ctx context.Context, req *http.Request) (interface{}, int) {
	return f(ctx, req)
}

// ServeHTTP implement http.handler interface
func (f ApiHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), HeaderKey, newHTTPHeaders())

	rsp, _ := f.Serve(ctx, r)

	header, _ := ctx.Value(HeaderKey).(*httpHeaders)
	if header != nil {
		headers := header.getHeader()
		for headerkey, headervalue := range headers {
			rw.Header().Add(headerkey, headervalue)
		}
	}

	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(rsp); err != nil {
		http.Error(rw, "Internal Error", http.StatusInternalServerError)
	}
}

// ----------------------------------------------------------------------------
// HTTP Server
// ----------------------------------------------------------------------------

type HTTPServer struct {
	mux    *http.ServeMux
	host   string
	port   int
	prefix string
}

func NewHTTPServer(host string, port int, prefix string) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		host:   host,
		port:   port,
		prefix: prefix,
	}
}

func (s *HTTPServer) RegisterApiHandler(url string, handler ApiHandler) {
	if len(url) == 0 {
		logging.Errorln("register url is invalid")
	}

	logging.Debug("mapping handler", logging.Fields{
		"prefix":  s.prefix,
		"url":     url,
		"handler": reflect.TypeOf(handler),
	})
	s.mux.Handle(fmt.Sprintf("%s%s", s.prefix, url), loggingHandler(handler))
}

func (s *HTTPServer) Run() {
	server := &http.Server{
		Handler:     s.mux,
		Addr:        fmt.Sprintf("%s:%s", s.host, strconv.Itoa(s.port)),
		ReadTimeout: 60 * time.Second,
	}

	//gracefully shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logging.Infoln("http server Shutdown")
		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logging.Fatalln("http server Shutdown error")
		}
		close(idleConnsClosed)
	}()

	logging.Info("http server started and served", logging.Fields{
		"host": s.host,
		"port": s.port,
	})
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logging.Fatal("HTTP server ListenAndServe", logging.Fields{
			"error": err.Error(),
		})
	}

	<-idleConnsClosed
}

// ----------------------------------------------------------------------------
// logging
// ---------------------------------------------------------------------------

func loggingHandler(next ApiHandler) ApiHandlerFunc {
	f := func(ctx context.Context, req *http.Request) (interface{}, int) {
		start := time.Now()
		logging.Info("request", logging.Fields{
			"method": req.Method,
			"url":    req.RequestURI,
			"ip":     req.RemoteAddr,
			"header": req.Header,
		})
		res, code := next.Serve(ctx, req)
		logging.Info("response", logging.Fields{
			"response": res,
			"code":     code,
			"url":      req.RequestURI,
			"duration": time.Since(start),
		})

		return res, code
	}

	return ApiHandlerFunc(f)
}

// ----------------------------------------------------------------------------
// HTTP headers
// ---------------------------------------------------------------------------

type httpHeaders struct {
	headers map[string]string
}

func newHTTPHeaders() *httpHeaders {
	return &httpHeaders{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func SetHeader(ctx context.Context, key string, value string) {
	if ctx.Value(HeaderKey) != nil {
		header, _ := ctx.Value(HeaderKey).(*httpHeaders)
		header.setHeader(key, value)
	}
}

func (t *httpHeaders) getHeader() map[string]string {
	if t.headers == nil {
		tmp := make(map[string]string)
		return tmp
	}
	return t.headers
}

func (t *httpHeaders) setHeader(key string, value string) {
	if t.headers == nil {
		t.headers = make(map[string]string)
	}
	t.headers[key] = value
}
