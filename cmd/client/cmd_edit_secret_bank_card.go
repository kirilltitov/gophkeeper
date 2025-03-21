package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func cmdEditSecretBankCard() *cli.Command {
	return &cli.Command{
		Name:        "edit-bank-card",
		Description: "Edits secret bank card (name, number, exp date and CVV)",
		Usage:       "Edits secret bank card",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSecretName,
				Usage:    "Secret name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagCardHolder,
				Usage:    "Cardholder name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagCardNumber,
				Usage:    "Card number",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagCardDate,
				Usage:    "Card expiration date",
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagCardCVV,
				Usage:    "CVV/CVC",
				Required: true,
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

			var cardHolder = cmd.String(flagCardHolder)
			var cardNumber = cmd.String(flagCardNumber)
			var cardDate = cmd.String(flagCardDate)
			var cardCVV = cmd.String(flagCardCVV)

			if existingSecret.IsEncrypted {
				fmt.Fprint(w, noticeSecretIsEncrypted)
				encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, true)
				if err != nil {
					return err
				}

				cardHolder, err = encrypt(encryptionKeyBytes, []byte(cardHolder))
				if err != nil {
					return err
				}

				cardNumber, err = encrypt(encryptionKeyBytes, []byte(cardNumber))
				if err != nil {
					return err
				}

				cardDate, err = encrypt(encryptionKeyBytes, []byte(cardDate))
				if err != nil {
					return err
				}

				cardCVV, err = encrypt(encryptionKeyBytes, []byte(cardCVV))
				if err != nil {
					return err
				}
			}

			req := api.SecretBankCard{
				Name:   cardHolder,
				Number: cardNumber,
				Date:   cardDate,
				CVV:    cardCVV,
			}

			code, err := SendRequest[any](
				c,
				ctx,
				fmt.Sprintf("/api/secret/edit/bank_card/%s", existingSecret.ID),
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

			fmt.Fprintf(w, "Successfully edited secret bank card '%s'", existingSecret.Name)

			return nil
		},
	}
}
