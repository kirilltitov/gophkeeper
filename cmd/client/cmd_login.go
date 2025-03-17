package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

const flagLogin = "login"

func cmdLogin() *cli.Command {
	return &cli.Command{
		Name:        "login",
		Description: "Performs login into the service using given login and password",
		Usage:       "Performs login into the service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: flagLogin,
			},
		},
		Before: setup,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var login string

			login = cmd.String(flagLogin)

			if login == "" {
				login = strings.TrimSpace(strings.Join(cmd.Args().Slice(), " "))
			}

			if login == "" {
				return fmt.Errorf("you haven't provided login (--%s or argument)", flagLogin)
			}

			claims, err := getAuthClaims()
			if err != nil {
				return err
			}
			if claims.Login == login && claims.Valid() == nil {
				fmt.Fprintf(
					cmd.Root().Writer,
					"Note: You don't have to login as there is an active session for user '%s' until %s\n",
					claims.Login,
					claims.ExpiresAt.String(),
				)
			}

			password, err := readPassword(cmd.Root().Writer, "Enter password: ")
			if err != nil {
				return err
			}

			c := newClient(cmd.String(flagAddress), "")

			type req struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			resp, err := c.SendRawRequest(
				ctx,
				"/api/login",
				http.MethodPost,
				req{
					Login:    login,
					Password: password,
				},
			)
			if err != nil {
				if errors.Is(err, errUnauthorized) {
					return errors.New("incorrect login or password")
				}
				return errors.Wrap(err, "could not login")
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code %d", resp.StatusCode)
			}

			jwtCookie := findAuthCookie(resp.Cookies(), cmd.String(flagAuthCookieName))
			if jwtCookie == nil {
				return errors.New("no auth cookie in response")
			}

			if err := storeJWT(jwtCookie.Value); err != nil {
				return errors.Wrap(err, "could not save JWT locally")
			}

			fmt.Fprintf(cmd.Root().Writer, "Successfully logged in\n")

			return cmdSync().Run(ctx, cmd.Args().Slice())
		},
	}
}
