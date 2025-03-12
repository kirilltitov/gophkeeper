package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func cmdEditSecretNote() *cli.Command {
	return &cli.Command{
		Name:        "edit-note",
		Description: "Edits secret note (plain text)",
		Usage:       "Edits secret note",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  flagSecretNoteFile,
				Usage: "Path to the file with secret note",
			},
			&cli.StringFlag{
				Name:  flagSecretNoteText,
				Usage: "Secret note text",
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

			fileName := cmd.String(flagSecretNoteFile)
			text := cmd.String(flagSecretNoteText)

			if fileName != "" && text != "" {
				return errors.New("both filename and text provided, choose one")
			}
			if fileName == "" && text == "" {
				return fmt.Errorf(
					"you haven't provided secret note (--%s) or the filename (--%s)",
					flagSecretNoteText,
					flagSecretNoteFile,
				)
			}

			var note string

			if fileName != "" {
				noteBytes, err := os.ReadFile(fileName)
				if err != nil {
					return errors.Wrap(err, "could not read secret note file")
				}
				note = string(noteBytes)
			}

			if text != "" {
				note = text
			}

			if existingSecret.IsEncrypted {
				fmt.Fprint(w, noticeSecretIsEncrypted)
				encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, true)
				if err != nil {
					return err
				}

				note, err = encrypt(encryptionKeyBytes, []byte(note))
				if err != nil {
					return err
				}
			}

			req := api.SecretNote{
				Body: note,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/edit/note/%s", existingSecret.ID),
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

			fmt.Fprintf(w, "Successfully edited secret note '%s'", existingSecret.Name)

			return nil
		},
	}
}
