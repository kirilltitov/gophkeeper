package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

const (
	flagOutput = "output"
)

//nolint:gocognit,maintidx // разбиение функции только усугубит её читабельность
func cmdGetSecret() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Gets secret with value",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagSecretName,
				Usage: "Secret name",
			},
			&cli.StringFlag{
				Name:    flagOutput,
				Aliases: []string{"o"},
				Usage:   "Outputs secret into provided file name (will create if not exists)",
			},
		},
		Before: setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			w := cmd.Root().Writer

			var err error

			name := cmd.String(flagSecretName)

			if name == "" {
				name = strings.TrimSpace(strings.Join(cmd.Args().Slice(), " "))
			}

			if name == "" {
				return fmt.Errorf(
					"you haven't provided secret name (as --%s flag or argument)",
					flagSecretName,
				)
			}

			if err := loadLocalSecrets(); err != nil {
				return err
			}

			if err := syncSecrets(ctx); err != nil {
				if isOffline(err) {
					fmt.Fprint(w, "Notice: Client is offline, using local secret (it might be outdated)\n\n")
				} else {
					return err
				}
			}

			name = strings.Trim(name, `"`)

			existingSecret, found := secretsByName[name]
			if !found {
				return fmt.Errorf("secret '%s' not found", name)
			}

			outputFileName := cmd.String(flagOutput)
			if existingSecret.Kind == api.KindBlob && outputFileName == "" {
				return fmt.Errorf(
					"secret '%s' is of type blob, and you haven't provided --output path for result",
					existingSecret.Name,
				)
			}

			var encryptionKeyBytes []byte
			if existingSecret.IsEncrypted {
				fmt.Fprint(w, noticeSecretIsEncrypted)
				encryptionKeyBytes, err = getEncryptionKeyBytes(cmd, true)
				if err != nil {
					return err
				}
			}

			var result []byte

			switch existingSecret.Kind {
			case api.KindBankCard:
				var value api.SecretBankCard
				if err := json.Unmarshal(existingSecret.Value, &value); err != nil {
					return errors.Wrap(err, "could not unmarshal secret bank card")
				}
				if existingSecret.IsEncrypted {
					var decryptedBytes []byte

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Name)
					if err != nil {
						return err
					}
					value.Name = string(decryptedBytes)

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Number)
					if err != nil {
						return err
					}
					value.Number = string(decryptedBytes)

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Date)
					if err != nil {
						return err
					}
					value.Date = string(decryptedBytes)

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.CVV)
					if err != nil {
						return err
					}
					value.CVV = string(decryptedBytes)
				}

				result = []byte(fmt.Sprintf(
					"Cardholder: %s\nNumber: %s\nExpiration date: %s\nCVV/CVC: %s\n",
					value.Name, value.Number, value.Date, value.CVV,
				))
			case api.KindCredentials:
				var value api.SecretCredentials
				if err := json.Unmarshal(existingSecret.Value, &value); err != nil {
					return errors.Wrap(err, "could not unmarshal secret credentials")
				}
				if existingSecret.IsEncrypted {
					var decryptedBytes []byte

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.URL)
					if err != nil {
						return err
					}
					value.URL = string(decryptedBytes)

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Login)
					if err != nil {
						return err
					}
					value.Login = string(decryptedBytes)

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Password)
					if err != nil {
						return err
					}
					value.Password = string(decryptedBytes)
				}

				result = []byte(fmt.Sprintf(
					"URL: %s\nLogin: %s\nPassword: %s\n",
					value.URL, value.Login, value.Password,
				))
			case api.KindNote:
				var value api.SecretNote
				if err := json.Unmarshal(existingSecret.Value, &value); err != nil {
					return errors.Wrap(err, "could not unmarshal secret note")
				}
				if existingSecret.IsEncrypted {
					var decryptedBytes []byte

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Body)
					if err != nil {
						return err
					}
					value.Body = string(decryptedBytes)
				}

				result = []byte(fmt.Sprintf(
					"%s\n",
					value.Body,
				))
			case api.KindBlob:
				var value api.SecretBlob
				if err := json.Unmarshal(existingSecret.Value, &value); err != nil {
					return errors.Wrap(err, "could not unmarshal secret blob")
				}
				if existingSecret.IsEncrypted {
					var decryptedBytes []byte

					decryptedBytes, err = decrypt(encryptionKeyBytes, value.Body)
					if err != nil {
						return err
					}
					result = decryptedBytes
				} else {
					result, err = base64.StdEncoding.DecodeString(value.Body)
					if err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("unexpected kind '%s'", existingSecret.Kind)
			}

			if existingSecret.Kind == api.KindBlob {
				if err := os.WriteFile(outputFileName, result, 0o660); err != nil {
					return errors.Wrap(err, "could not write secret blob to output file")
				}

				fmt.Fprintf(w, "Successfully written your secret to file %s\n", outputFileName)

				return nil
			}

			if outputFileName != "" {
				if err := os.WriteFile(outputFileName, result, 0o660); err != nil {
					return errors.Wrap(err, "could not write secret blob to output file")
				}

				fmt.Fprintf(w, "Successfully written your secret to file %s\n", outputFileName)

				return nil
			}

			fmt.Fprintf(w, "Name: %s\n", existingSecret.Name)

			if existingSecret.Description != "" {
				fmt.Fprintf(w, "Description: %s\n", existingSecret.Description)
			}

			if len(existingSecret.Tags) > 0 {
				fmt.Fprintf(w, "Tags: %s\n", strings.Join(existingSecret.Tags, ", "))
			}

			fmt.Fprintf(w, "\n%s", string(result))

			return nil
		},
	}
}
