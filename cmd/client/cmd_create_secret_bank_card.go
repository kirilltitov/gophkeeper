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
	flagCardHolder = "cardholder"
	flagCardNumber = "number"
	flagCardDate   = "date"
	flagCardCVV    = "cvv"
)

func cmdCreateSecretBankCard() *cli.Command {
	return &cli.Command{
		Name:        "create-bank-card",
		Description: "Creates secret bank card (name, number, exp date and CVV)",
		Usage:       "Creates secret bank card",
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

			encryptionKeyBytes, err := getEncryptionKeyBytes(cmd, false)
			if err != nil {
				return nil
			}
			isEncryptionEnabled := !cmd.Bool(flagNoEncrypt)

			var cardHolder = cmd.String(flagCardHolder)
			var cardNumber = cmd.String(flagCardNumber)
			var cardDate = cmd.String(flagCardDate)
			var cardCVV = cmd.String(flagCardCVV)

			if encryptionKeyBytes != nil {
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

			req := api.BaseCreateSecretRequest[api.SecretBankCard]{
				Name:        cmd.String(flagSecretName),
				Description: cmd.String(flagSecretDescription),
				IsEncrypted: isEncryptionEnabled,
				Value: api.SecretBankCard{
					Name:   cardHolder,
					Number: cardNumber,
					Date:   cardDate,
					CVV:    cardCVV,
				},
			}

			var resp api.CreatedSecretResponse

			code, err := SendRequest(c, ctx, "/api/secret/create/bank_card", http.MethodPost, req, &resp)
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

			fmt.Fprintf(w, "Successfully created secret bank card '%s' with id '%s'", req.Name, resp.ID.String())

			return nil
		},
	}
}
