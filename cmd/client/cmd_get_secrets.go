package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func cmdGetSecrets() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "Secrets list",
		Aliases: []string{"secrets"},
		Before:  setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			fmt.Fprintf(w, "[ID] [Kind] Name Details\n\n")

			for _, item := range secretsByName {
				var details []string
				if item.IsEncrypted {
					details = append(details, "🔑")
				}
				if len(item.Tags) > 0 {
					details = append(details, fmt.Sprintf("🏷: %s", strings.Join(item.Tags, ", ")))
				}
				fmt.Fprintf(
					w,
					`[%s] [%-11s] "%s" %s%s`,
					item.ID, item.Kind, item.Name, strings.Join(details, " "), "\n",
				)
			}

			return nil
		},
	}
}
