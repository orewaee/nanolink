package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driving"
	"github.com/orewaee/nanolink/internal/datasource/disk"
	delivery "github.com/orewaee/nanolink/internal/delivery/redirect"
	"github.com/orewaee/nanolink/internal/gen"
	"github.com/orewaee/nanolink/internal/usecase"
	"github.com/urfave/cli/v3"
)

func main() {
	dir := "."
	linkApi := usecase.NewLinkService(&disk.YamlLinkRepo{Dir: dir + "/links"})

	cmd := &cli.Command{
		Name: "nanolink",
		Commands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "port",
						Value:    2222,
						Required: false,
					},
					&cli.StringFlag{
						Name:     "host",
						Value:    "127.0.0.1",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "tls",
						Value:    false,
						Required: false,
					},
					&cli.StringFlag{
						Name:     "cert-file",
						Value:    "",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "key-file",
						Value:    "",
						Required: false,
					},
				},
				Commands: []*cli.Command{
					{
						Name: "redirect",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							host := cmd.String("host")
							port := cmd.Int("port")

							certFile := cmd.String("cert-file")
							keyFile := cmd.String("key-file")
							tls := cmd.Bool("tls")
							if !tls {
								certFile = ""
								keyFile = ""
							}

							runRedirectDelivery(ctx, host, port, tls, certFile, keyFile, linkApi)
							return nil
						},
					},
					{
						Name: "dashboard",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							port := cmd.Int("port")
							log.Println("DASHBOARD PORT", port)
							return nil
						},
					},
				},
			},
			{
				Name: "add",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "id",
					},
				},
				ArgsUsage: "<location>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.String("id")
					if id == "" {
						provider := gen.NewAlphabeticIdProvider()
						generated, err := provider.GenerateId(ctx, 8)
						if err != nil {
							return err
						}

						id = generated
					}

					location := cmd.Args().First()
					if location == "" {
						return cli.Exit("location is required", 1)
					}

					link, err := linkApi.AddLink(ctx, id, location)
					if err != nil {
						return err
					}

					fmt.Printf("/%s -> %s\n", link.Id, link.Location)
					return nil
				},
			},
			{
				Name: "remove",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "id",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.StringArg("id")
					if id == "" {
						return cli.Exit("id is required", 1)
					}

					if err := linkApi.RemoveLinkById(ctx, id); err != nil {
						return err
					}

					fmt.Printf("%s removed\n", id)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runRedirectDelivery(
	ctx context.Context,
	host string,
	port int,
	tls bool,
	certFile string,
	keyFile string,
	linkApi driving.LinkApi,
) {
	opts := []domain.RedirectOpt{
		domain.WithHost(host),
		domain.WithPort(port),
	}

	if tls {
		opts = append(opts, domain.WithTLS(certFile, keyFile))
	}

	controller := delivery.NewHttpRedirectController(linkApi, opts...)
	go controller.Run()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	fmt.Println("shutting down redirect delivery...")
	controller.Shutdown(ctx)
}
