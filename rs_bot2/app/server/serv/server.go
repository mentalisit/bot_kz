package servprof

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"
)

// Server представляет HTTP-сервер.
type Server struct {
	logger *log.Logger
}

// NewServer создает и запускает новый сервер.
func NewServer() *Server {
	s := &Server{logger: log.New(os.Stdout, "", 0)}
	go s.runServer()
	return s
}

func (s *Server) runServer() {
	// Создаем маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем маршруты pprof
	registerPprofRoutes(mux)

	// Запускаем сервер
	addr := ":80"
	s.logger.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		s.logger.Printf("Server failed: %v", err)
	}
}

// registerPprofRoutes регистрирует обработчики pprof на маршрутизаторе.
func registerPprofRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Регистрируем дополнительные обработчики pprof
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
