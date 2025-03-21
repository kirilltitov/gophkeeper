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

func cmdAddTag() *cli.Command {
	return &cli.Command{
		Name:        "tag",
		Description: "Adds tag to a secret (plain text)",
		Usage:       "Adds tag to a secret",
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

			tag := strings.TrimSpace(strings.Join(cmd.Args().Slice(), " "))
			if tag == "" {
				return errors.New("you haven't provided tag")
			}

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			name := cmd.String(flagSecretName)
			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			req := api.TagRequest{
				Tag: tag,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/tag/%s", existingSecret.ID),
				http.MethodPost,
				req,
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

			fmt.Fprintf(w, "Successfully added tag '%s' to secret '%s'", tag, existingSecret.Name)

			return nil
		},
	}
}
