package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
				},
				Commands: []*cli.Command{
					{
						Name: "redirect",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							host := cmd.String("host")
							port := cmd.Int("port")
							addr := fmt.Sprintf("%s:%d", host, port)
							runRedirectDelivery(ctx, addr, linkApi)
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runRedirectDelivery(ctx context.Context, addr string, linkApi driving.LinkApi) {
	controller := delivery.NewHttpRedirectController(addr, linkApi)

	fmt.Printf("running redirect delivery on %s\n", addr)
	go controller.Run()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	fmt.Println("shutting down redirect delivery...")
	controller.Shutdown(ctx)
}
