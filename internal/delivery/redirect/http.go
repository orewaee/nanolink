package delivery

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driving"
	"github.com/orewaee/nanolink/internal/delivery"
)

type HttpRedirectController struct {
	httpServer   *http.Server
	linkApi      driving.LinkApi
	redirectOpts domain.RedirectOpts
}

func NewHttpRedirectController(
	linkApi driving.LinkApi,
	opts ...domain.RedirectOpt,
) delivery.Controller {
	controller := &HttpRedirectController{
		linkApi:      linkApi,
		redirectOpts: domain.DefaultRedirectOpts,
	}

	for _, opt := range opts {
		if err := opt(&controller.redirectOpts); err != nil {
			panic(err)
		}
	}

	return controller
}

func (controller *HttpRedirectController) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

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
			if controller.redirectOpts.NotFoundTempl {
				templ, err := template.ParseFiles("templates/404.templ")
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					return
				}

				templ.Execute(writer, nil)
				return
			}

			writer.WriteHeader(http.StatusNotFound)
			return
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	opts := controller.redirectOpts

	host := opts.Host
	if opts.HostType == domain.IPv6 {
		host = fmt.Sprintf("[%s]", host)
	}

	addr := fmt.Sprintf("%s:%d", host, opts.Port)
	controller.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("running redirect delivery on %s\n", addr)

	if opts.TLS {
		err := controller.httpServer.ListenAndServeTLS(opts.CertFile, opts.KeyFile)
		if err != nil {
			return err
		}
	} else {
		err := controller.httpServer.ListenAndServe()
		if err != nil {
			return err
		}
	}

	return nil
}

func (controller *HttpRedirectController) Shutdown(ctx context.Context) error {
	return controller.httpServer.Shutdown(ctx)
}
