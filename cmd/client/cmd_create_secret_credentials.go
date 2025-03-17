package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

const (
	flagSecretURL   = "url"
	flagSecretLogin = "login"
)

func cmdCreateSecretCredentials() *cli.Command {
	return &cli.Command{
		Name:        "create-credentials",
		Description: "Creates secret credentials (login and password)",
		Usage:       "Creates secret credentials",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  flagSecretDescription,
				Usage: "Secret description",
			},
			&cli.StringFlag{
				Name:     flagSecretURL,
				Usage:    "Credentials URL",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagSecretLogin,
				Usage:    "Credentials login",
				Required: true,
			},
		},
		Before: setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, false)
			if err != nil {
				return nil
			}
			isEncryptionEnabled := !cmd.Bool(flagNoEncrypt)

			password, err := readPassword(w, "Enter secret credentials password: ")
			if err != nil {
				return nil
			}

			var login = cmd.String(flagSecretLogin)
			var URL = cmd.String(flagSecretURL)

			if encryptionKeyBytes != nil {
				URL, err = encrypt(encryptionKeyBytes, []byte(URL))
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

			req := api.BaseCreateSecretRequest[api.SecretCredentials]{
				Name:        cmd.String(flagSecretName),
				Description: cmd.String(flagSecretDescription),
				IsEncrypted: isEncryptionEnabled,
				Value: api.SecretCredentials{
					URL:      URL,
					Login:    login,
					Password: password,
				},
			}

			var resp api.CreatedSecretResponse

			code, err := SendRequest(c, ctx, "/api/secret/create/credentials", http.MethodPost, req, &resp)
			if err != nil {
				return err
			}
			if code != http.StatusCreated {
				switch code {
				case http.StatusConflict:
					return errors.New("secret with this name already exists")
				default:
					return fmt.Errorf("unexpected status code %d", code)
				}
			}

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			fmt.Fprintf(w, "Successfully created secret credentials '%s' with id '%s'", req.Name, resp.ID.String())

			return nil
		},
	}
}
