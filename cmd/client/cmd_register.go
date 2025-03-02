package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

func cmdRegister() *cli.Command {
	return &cli.Command{
		Name:        "register",
		Description: "Performs signup into the service using given login and password",
		Usage:       "Performs signup into the service",
		Aliases:     []string{"signup"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagLogin,
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Root().Writer, "About to register user with login '%s'\n", cmd.String(flagLogin))

			password1, err := readPassword(cmd.Root().Writer, "Enter password: ")
			if err != nil {
				return err
			}

			password2, err := readPassword(cmd.Root().Writer, "Repeat password: ")
			if err != nil {
				return err
			}

			if password1 != password2 {
				return errors.New("entered passwords don't match")
			}

			c := newClient(cmd.String(flagAddress), "")

			type req struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			resp, err := c.SendRawRequest(
				ctx,
				"/api/register",
				http.MethodPost,
				req{
					Login:    cmd.String(flagLogin),
					Password: password1,
				},
			)
			if err != nil {
				return errors.Wrap(err, "could not register")
			}

			if resp.StatusCode != 200 {
				switch resp.StatusCode {
				case http.StatusConflict:
					return errors.New("user with given login already exists")
				case http.StatusBadRequest:
					return errors.New("empty login or password")
				default:
					return fmt.Errorf("unexpected status code %d", resp.StatusCode)
				}
			}

			jwtCookie := findAuthCookie(resp.Cookies(), cmd.String(flagAuthCookieName))
			if jwtCookie == nil {
				return errors.New("no auth cookie in response")
			}

			if err := storeJWT(jwtCookie.Value); err != nil {
				return fmt.Errorf("could not save JWT locally: %s", err.Error())
			}

			fmt.Fprintf(cmd.Root().Writer, "Successfuly registered\n")

			return nil
		},
	}
}
