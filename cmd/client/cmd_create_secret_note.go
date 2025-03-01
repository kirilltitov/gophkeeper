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

const (
	flagSecretNoteFile = "file"
	flagSecretNoteText = "text"
)

func cmdCreateSecretNote() *cli.Command {
	return &cli.Command{
		Name:        "create-note",
		Description: "Creates secret note (any plain text content)",
		Usage:       "Creates secret note",
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			fileName := cmd.String(flagSecretNoteFile)
			text := cmd.String(flagSecretNoteText)

			if fileName != "" && text != "" {
				return errors.New("both filename and text provided, choose one")
			}
			if fileName == "" && text == "" {
				return errors.New(fmt.Sprintf(
					"you haven't provided secret note (--%s) or the filename (--%s)",
					flagSecretNoteText,
					flagSecretNoteFile,
				))
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

			encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, false)
			isEncryptionEnabled := !cmd.Bool(flagNoEncrypt)

			if encryptionKeyBytes != nil {
				note, err = encrypt(encryptionKeyBytes, []byte(note))
				if err != nil {
					return err
				}
			}

			req := api.BaseCreateSecretRequest[api.SecretNote]{
				Name:        cmd.String(flagSecretName),
				IsEncrypted: isEncryptionEnabled,
				Value: api.SecretNote{
					Body: note,
				},
			}

			var resp api.CreatedSecretResponse

			code, err := SendRequest(c, ctx, "/api/secret/create/note", http.MethodPost, req, &resp)
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

			fmt.Fprintf(w, "Succesfully created secret note '%s' with id '%s'", req.Name, resp.ID.String())

			return nil
		},
	}
}
