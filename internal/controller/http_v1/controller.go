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

// HelloWorld
// @Summary dummy controller
// @Description Returns "Hello, World!" on a GET request
// @Tags hello-world
// @Success 200 {string} string "Hello, World!"
// @Failure 405 {string} string "Method not allowed"
// @Router /hello [get]
func (c *HttpController) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
