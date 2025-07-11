package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/orewaee/nanolink/internal/datasource/disk"
	delivery "github.com/orewaee/nanolink/internal/delivery/redirect"
	"github.com/orewaee/nanolink/internal/usecase"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "nanolink",
		Commands: []*cli.Command{
			{
				Name: "run",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "port",
						Value: 2222,
						Usage: "nanolink instance port",
					},
				},
				Commands: []*cli.Command{
					{
						Name: "redirect",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							port := cmd.Int("port")
							log.Println("REDIRECT PORT", port)
							runRedirectDelivery(port)
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runRedirectDelivery(port int) {
	linkApi := usecase.NewLinkService(&disk.YamlLinkRepo{Dir: "links"})
	controller := delivery.NewHttpRedirectController(port, linkApi)

	fmt.Printf("running redirect delivery on :%d\n", port)
	go controller.Run()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	fmt.Println("shutting down redirect delivery...")
	controller.Shutdown(context.Background())
}
