package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

var dsn string

func RunMigratorCLI(migrations *migrate.Migrations) error {
	app := &cli.App{
		Name: "migrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dsn",
				Destination: &dsn,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(ctx *cli.Context) error {
					return getMigrator(migrations).Init(ctx.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(ctx *cli.Context) error {
					if err := getMigrator(migrations).Init(ctx.Context); err != nil {
						return err
					}
					group, err := getMigrator(migrations).Migrate(ctx.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no new migrations to run (database is up to date)\n")
						return nil
					}
					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(ctx *cli.Context) error {
					group, err := getMigrator(migrations).Rollback(ctx.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}
					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")
					mf, err := getMigrator(migrations).CreateGoMigration(ctx.Context, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create SQL migration",
				Action: func(ctx *cli.Context) error {
					name := strings.Join(ctx.Args().Slice(), "_")
					mf, err := getMigrator(migrations).CreateSQLMigrations(ctx.Context, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf[0].Name, mf[0].Path)
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(ctx *cli.Context) error {
					mgs, err := getMigrator(migrations).MigrationsWithStatus(ctx.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", mgs)
					fmt.Printf("unapplied migrations: %s\n", mgs.Unapplied())
					fmt.Printf("last migration group: %s\n", mgs.LastGroup())
					return nil
				},
			},
		},
	}

	return app.Run(os.Args)
}

func getDb(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	return bun.NewDB(sqldb, pgdialect.New())
}

func getMigrator(migrations *migrate.Migrations) *migrate.Migrator {
	db := getDb(dsn)

	return migrate.NewMigrator(db, migrations)
}
