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
				Name:     flagSecretLogin,
				Usage:    "Credentials login",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, false)
			isEncryptionEnabled := !cmd.Bool(flagNoEncrypt)

			password, err := readPassword(w, "Enter secret credentials password: ")
			if err != nil {
				return nil
			}

			var login = cmd.String(flagSecretLogin)

			if encryptionKeyBytes != nil {
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
				IsEncrypted: isEncryptionEnabled,
				Value: api.SecretCredentials{
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
				case http.StatusUnauthorized:
					return errors.New("unauthorized")
				case http.StatusConflict:
					return errors.New("secret with this name already exists")
				default:
					return errors.New(fmt.Sprintf("unexpected status code %d", code))
				}
			}

			fmt.Fprintf(w, "Succesfully created secret credentials '%s' with id '%s'", req.Name, resp.ID.String())

			return nil
		},
	}
}
