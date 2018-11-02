package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lestrrat-go/server-starter/listener"
	"github.com/muranoya/echo/log"
	"go.uber.org/zap"
)

// Env is HTTP environment
type Env struct {
}

// Handler is a custom handler struct.
type Handler struct {
	*Env
	HandleFunc func(e *Env, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.RequestLog().Info("Request",
		zap.String("ua", r.UserAgent()),
		zap.String("ip", r.RemoteAddr),
		zap.String("ref", r.Referer()))

	if err := h.HandleFunc(h.Env, w, r); err != nil {
		log.SystemLog().Warn("Request handler failed", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func helloHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "world")
	return nil
}

func main() {
	log.SystemLog().Info("Startup", zap.Int("pid", os.Getpid()))
	listeners, err := listener.ListenAll()
	if err != nil {
		os.Exit(1)
	}

	env := &Env{}
	handler := http.NewServeMux()
	handler.Handle("/hello", Handler{Env: env, HandleFunc: helloHandler})
	for _, l := range listeners {
		go func(l net.Listener) {
			if err := http.Serve(l, handler); err != nil {
				log.SystemLog().Error("http.Serve failed", zap.Error(err))
			}
		}(l)
	}

	loop := true
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM)
	for loop {
		s := <-sigCh
		switch s {
		case syscall.SIGTERM:
			log.SystemLog().Info("Receive SIGTERM")
			for _, l := range listeners {
				l.Close()
			}
			time.Sleep(5 * time.Second)
			loop = false
		default:
			time.Sleep(time.Second)
		}
	}
	log.SystemLog().Info("Good bye")
}
