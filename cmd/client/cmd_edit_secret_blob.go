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

func cmdEditSecretBlob() *cli.Command {
	return &cli.Command{
		Name:        "edit-blob",
		Description: "Edits secret blob (plain text)",
		Usage:       "Edits secret blob",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  flagSecretBlobFile,
				Usage: "Path to the file with secret blob",
			},
		},
		Before: setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			var err error

			if err := syncSecrets(ctx); err != nil {
				return err
			}

			name := cmd.String(flagSecretName)
			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			fileName := cmd.String(flagSecretBlobFile)
			blobBytes, err := os.ReadFile(fileName)
			if err != nil {
				return errors.Wrap(err, "could not read blob file")
			}

			var blob string
			if existingSecret.IsEncrypted {
				fmt.Fprint(w, noticeSecretIsEncrypted)
				encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, true)
				if err != nil {
					return err
				}

				blob, err = encrypt(encryptionKeyBytes, blobBytes)
				if err != nil {
					return err
				}
			} else {
				blob = base64.StdEncoding.EncodeToString(blobBytes)
			}

			req := api.SecretBlob{
				Body: blob,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/edit/blob/%s", existingSecret.ID),
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

			fmt.Fprintf(w, "Successfully edited secret blob '%s'", existingSecret.Name)

			return nil
		},
	}
}
