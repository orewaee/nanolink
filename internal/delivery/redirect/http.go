package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driving"
	"github.com/orewaee/nanolink/internal/delivery"
)

type HttpRedirectController struct {
	port       int
	httpServer *http.Server
	linkApi    driving.LinkApi
}

func NewHttpRedirectController(port int, linkApi driving.LinkApi) delivery.Controller {
	return &HttpRedirectController{
		port:    port,
		linkApi: linkApi,
	}
}

func (controller *HttpRedirectController) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")
		link, err := controller.linkApi.GetLinkById(request.Context(), id)
		if err == nil {
			writer.Header().Add("Location", link.Location)
			writer.WriteHeader(http.StatusPermanentRedirect)
			return
		}

		switch {
		case errors.Is(err, domain.ErrLinkNotFound):
			writer.WriteHeader(http.StatusNotFound)
			return
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	controller.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", controller.port),
		Handler: mux,
	}

	return controller.httpServer.ListenAndServe()
}

func (controller *HttpRedirectController) Shutdown(ctx context.Context) error {
	return controller.httpServer.Shutdown(ctx)
}
