package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
		Before: checkAuth,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			name := cmd.String(flagSecretName)

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
				switch code {
				case http.StatusUnauthorized:
					return errors.New("unauthorized")
				default:
					return fmt.Errorf("unexpected status code %d", code)
				}
			}

			fmt.Fprintf(w, "Successfully delete secret '%s'", name)

			return nil
		},
	}
}
