package migrate

import (
	"fmt"
	migrateV4 "github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	// import posgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// import file
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func CliCommand(sourceURL string, databaseURL string) []*cli.Command {
	// Migration should always run on development mode
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return []*cli.Command{
		{
			Name:  "up",
			Usage: "Migrate up data",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					log.Fatal("Error create migration", zap.Error(err))
				}

				logger.Info("migration up")
				if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
					logger.Fatal(err.Error())
				}
				return err
			},
		},
		{
			Name:  "down",
			Usage: "step down migration by N(int)",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					logger.Fatal("Error create migration", zap.Error(err))
				}

				down, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					logger.Fatal("rev should be a number", zap.Error(err))
				}

				logger.Info("migration down", zap.Int("down", -down))
				if err := m.Steps(-down); err != nil {
					logger.Fatal(err.Error())
				}
				return err
			},
		},
		{
			Name: "create",
			Action: func(c *cli.Context) error {
				folder := strings.ReplaceAll(sourceURL, "file://", "")
				now := time.Now()
				versionTimeFormat := "20060102150405"
				ver := now.Format(versionTimeFormat)
				name := strings.Join(c.Args().Slice(), "-")

				up := fmt.Sprintf("%s/%s_%s.up.sql", folder, ver, name)
				down := fmt.Sprintf("%s/%s_%s.down.sql", folder, ver, name)

				logger.Info("create migration", zap.String("name", name))
				logger.Info("up script", zap.String("up", up))
				logger.Info("down script", zap.String("down", up))

				if err := ioutil.WriteFile(up, []byte{}, 0600); err != nil {
					logger.Fatal("Create migration up error", zap.Error(err))
				}
				if err := ioutil.WriteFile(down, []byte{}, 0600); err != nil {
					logger.Fatal("Create migration down error", zap.Error(err))
				}
				return nil
			},
		},
	}
}
