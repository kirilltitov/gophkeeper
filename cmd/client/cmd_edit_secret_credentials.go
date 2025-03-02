package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func cmdEditSecretCredentials() *cli.Command {
	return &cli.Command{
		Name:        "edit-credentials",
		Description: "Edits secret credentials (login, password)",
		Usage:       "Edits secret credentials",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagSecretLogin,
				Usage:    "Credentials login",
				Required: true,
			},
		},
		Before: checkAuth,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			var err error

			name := cmd.String(flagSecretName)
			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			var login = cmd.String(flagLogin)

			password, err := readPassword(w, "Enter secret credentials password: ")
			if err != nil {
				return nil
			}

			if existingSecret.IsEncrypted {
				fmt.Fprintf(w, "This secret is encrypted, so you'll have to enter encryption key\n")
				encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, true)
				if err != nil {
					return err
				}

				login, err = encrypt(encryptionKeyBytes, []byte(login))
				if err != nil {
					return err
				}

				password, err = encrypt(encryptionKeyBytes, []byte(password))
				if err != nil {
					return err
				}
			}

			req := api.SecretCredentials{
				Login:    login,
				Password: password,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/edit/credentials/%s", existingSecret.ID),
				http.MethodPost,
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
					return fmt.Errorf("unexpected status code %d", code)
				}
			}

			fmt.Fprintf(w, "Successfully edited secret credentials '%s'", existingSecret.Name)

			return nil
		},
	}
}
