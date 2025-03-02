package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

const (
	flagSecretBlobFile = "file"
)

func cmdCreateSecretBlob() *cli.Command {
	return &cli.Command{
		Name:        "create-blob",
		Description: "Creates secret blob (any binary content)",
		Usage:       "Creates secret blob",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagSecretBlobFile,
				Usage:    "Path to the file with secret blob",
				Required: true,
			},
		},
		Before: checkAuth,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			fileName := cmd.String(flagSecretBlobFile)
			blobBytes, err := os.ReadFile(fileName)
			if err != nil {
				return errors.Wrap(err, "could not read blob file")
			}

			encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, false)
			if err != nil {
				return nil
			}
			isEncryptionEnabled := !cmd.Bool(flagNoEncrypt)

			var blob string
			if encryptionKeyBytes != nil {
				blob, err = encrypt(encryptionKeyBytes, blobBytes)
				if err != nil {
					return err
				}
			} else {
				blob = base64.StdEncoding.EncodeToString(blobBytes)
			}

			req := api.BaseCreateSecretRequest[api.SecretBlob]{
				Name:        cmd.String(flagSecretName),
				IsEncrypted: isEncryptionEnabled,
				Value: api.SecretBlob{
					Body: blob,
				},
			}

			var resp api.CreatedSecretResponse

			code, err := SendRequest(c, ctx, "/api/secret/create/blob", http.MethodPost, req, &resp)
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
					return fmt.Errorf("unexpected status code %d", code)
				}
			}

			fmt.Fprintf(w, "Successfully created secret blob '%s' with id '%s'", req.Name, resp.ID.String())

			return nil
		},
	}
}
