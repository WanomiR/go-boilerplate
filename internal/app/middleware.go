package app

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var totalRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "goboilerplate",
	Name:      "total_requests_count",
	Help:      "Total number of requests received on all API's",
}, []string{"type"})

func init() {
	prometheus.MustRegister(totalRequestsCount)
}

func (a *App) zapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		a.logger.Info("new request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("addr", r.RemoteAddr),
		)
	})
}

func (a *App) rateLimiter(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.RWMutex
		clients = make(map[string]*client)
	)

	// clear inactive clients
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 1*time.Hour {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get client's ip address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mu.RLock()
		_, found := clients[ip]
		mu.RUnlock()

		mu.Lock()
		if !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 100)}
		}
		clients[ip].lastSeen = time.Now()
		mu.Unlock()

		if !clients[ip].limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) requestsCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		totalRequestsCount.With(prometheus.Labels{"type": "requestst"}).Inc()
	})
}
