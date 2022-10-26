package cmd

import (
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	cfg *config.Config
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	app := cli.NewApp()
	app.Name = "service"
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Usage:  "Start grpc/http server",
			Action: serverAction,
		},
	}
	if app.Run(os.Args) != nil {
		panic(err)
	}
	return err
}

func serverAction(context *cli.Context) error {
	service, err := newService(cfg)
}

func newService(c *config.Config) (*service.Server, error) {

}
