package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func cmdDeleteTag() *cli.Command {
	return &cli.Command{
		Name:        "delete-tag",
		Description: "Deletes a tag from secret (plain text)",
		Usage:       "Deletes a tag from secret",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			tag := strings.TrimSpace(strings.Join(cmd.Args().Slice(), " "))
			if tag == "" {
				return errors.New("you haven't provided tag")
			}

			name := cmd.String(flagSecretName)
			existingSecret, found := secretsByName[name]
			if !found {
				return errors.New(fmt.Sprintf("secret '%s' not found", name))
			}

			req := api.TagRequest{
				Tag: tag,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/tag/%s", existingSecret.ID),
				http.MethodDelete,
				req,
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
					return errors.New(fmt.Sprintf("unexpected status code %d", code))
				}
			}

			fmt.Fprintf(w, "Succesfully deleted tag '%s' from secret '%s'", tag, existingSecret.Name)

			return nil
		},
	}
}
