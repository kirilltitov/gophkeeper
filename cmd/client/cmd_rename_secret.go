package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

const (
	flagSecretNewName = "new-name"
)

func cmdRenameSecret() *cli.Command {
	return &cli.Command{
		Name:  "rename-secret",
		Usage: "Renames secret",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagSecretNewName,
				Usage:    "New secret name",
				Required: true,
			},
		},
		Before: checkAuth,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			oldName := cmd.String(flagSecretName)
			newName := cmd.String(flagSecretNewName)

			existingSecret, found := secretsByName[oldName]
			if !found {
				return fmt.Errorf("secret '%s' not found", oldName)
			}

			type req struct {
				Name string `json:"name"`
			}

			code, err := SendRequest[req](
				c,
				ctx,
				fmt.Sprintf("/api/secret/%s/rename", existingSecret.ID),
				http.MethodPost,
				req{Name: newName},
				nil,
			)
			if err != nil {
				return err
			}
			if code != http.StatusOK {
				switch code {
				case http.StatusUnauthorized:
					return errors.New("unauthorized")
				case http.StatusConflict:
					return errors.New("secret with this name already exists")
				default:
					return fmt.Errorf("unexpected status code %d", code)
				}
			}

			fmt.Fprintf(w, "Succesfuly renamed secret from '%s' to '%s'", oldName, newName)

			return nil
		},
	}
}
