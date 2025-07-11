package delivery

import (
	"net/http"

	"github.com/orewaee/nanolink/internal/core/driving"
)

type RestController struct {
	linkApi driving.LinkApi
}

func (controller *RestController) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/v1/", controller.muxV1())

	return nil
}
