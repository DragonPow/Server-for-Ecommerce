package migrate

import (
	"fmt"
	log "github.com/DragonPow/Server-for-Ecommerce/library/log"
	migrateV4 "github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"io/ioutil"
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
	logger := log.MustBuildLogR()

	return []*cli.Command{
		{
			Name:  "up",
			Usage: "Migrate up data",
			Action: func(c *cli.Context) error {
				m, err := migrateV4.New(sourceURL, databaseURL)
				if err != nil {
					logger.Error(err, "Error create migration")
				}

				logger.Info("migration up")
				if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
					logger.Error(err, "Migrate up fail")
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
					logger.Error(err, "Error create migration")
				}

				down, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					logger.Error(err, "rev should be a number")
				}

				logger.Info("migration down", zap.Int("down", -down))
				if err := m.Steps(-down); err != nil {
					logger.Error(err, "Migrate down fail")
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
					logger.Error(err, "Create migration up error")
				}
				if err := ioutil.WriteFile(down, []byte{}, 0600); err != nil {
					logger.Error(err, "Create migration down error")
				}
				return nil
			},
		},
	}
}
