package httpv1

import (
	"net/http"

	"github.com/wanomir/rr"
	"go.uber.org/zap"
)

type HttpController struct {
	rr     *rr.ReadResponder
	logger *zap.Logger
}

func NewHttpController(rr *rr.ReadResponder, logger *zap.Logger) *HttpController {
	return &HttpController{
		rr:     rr,
		logger: logger,
	}
}

func (c *HttpController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
