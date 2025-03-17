package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"
)

func cmdDeleteSecret() *cli.Command {
	return &cli.Command{
		Name:  "delete-secret",
		Usage: "Deletes secret",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
		},
		Before: setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			name := cmd.String(flagSecretName)

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			type req struct {
				Name string `json:"name"`
			}

			code, err := SendRequest[req](
				c,
				ctx,
				fmt.Sprintf("/api/secret/%s", existingSecret.ID),
				http.MethodDelete,
				nil,
				nil,
			)
			if err != nil {
				return err
			}
			if code != http.StatusOK {
				return fmt.Errorf("unexpected status code %d", code)
			}

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			fmt.Fprintf(w, "Successfully delete secret '%s'", name)

			return nil
		},
	}
}
