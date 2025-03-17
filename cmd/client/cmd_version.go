package main

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/version"
)

func cmdVersion() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Prints client version, release date and git ref",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			v := version.Version{
				BuildVersion: buildVersion,
				BuildDate:    buildDate,
				BuildCommit:  buildCommit,
			}
			v.Print(cmd.Root().Writer)

			return nil
		},
	}
}
