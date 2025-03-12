package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"
)

func cmdChangeSecretDescription() *cli.Command {
	return &cli.Command{
		Name:  "change-description",
		Usage: "Changes secret description",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagSecretDescription,
				Usage:    "New secret description",
				Required: true,
			},
		},
		Before: setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			name := cmd.String(flagSecretName)
			newDescription := cmd.String(flagSecretDescription)

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			type req struct {
				Description string `json:"description"`
			}

			code, err := SendRequest[req](
				c,
				ctx,
				fmt.Sprintf("/api/secret/%s/change_description", existingSecret.ID),
				http.MethodPost,
				req{Description: newDescription},
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

			fmt.Fprint(w, "Successfully changed secret description\n")

			return nil
		},
	}
}
