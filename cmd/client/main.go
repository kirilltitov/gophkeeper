package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"

	"github.com/kirilltitov/gophkeeper/pkg/auth"
)

const (
	appDir   = "com.kirilltitov.gophkeeper"
	authFile = "auth.jwt"
)

const (
	flagAddress        = "address"
	flagAuthCookieName = "auth-cookie-name"
	flagNoEncrypt      = "no-encrypt"
	flagSecretName     = "name"
	flagVerbose        = "verbose"
	flagVeryVerbose    = "very-verbose"
)

const noticeSecretIsEncrypted = "This secret is encrypted, so you'll have to enter encryption key\n\n"

var logger = logrus.New()

var isLoggedIn bool

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

var c *client
var secretsByName = make(map[string]*secret)
var secretsByID = make(map[uuid.UUID]*secret)

type secret struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"`
	IsEncrypted bool      `json:"is_encrypted"`
	Tags        []string  `json:"tags"`
}

func main() {
	cmd := cli.Command{
		Name:        "Gophkeeper Client",
		Usage:       "",
		UsageText:   "",
		ArgsUsage:   "",
		Description: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagAddress,
				Usage:   "Address (including protocol and port) of the service",
				Aliases: []string{"a"},
			},
			&cli.StringFlag{
				Name:  flagAuthCookieName,
				Usage: "Authentication cookie name",
				Value: auth.DefaultCookieName,
			},
			&cli.BoolFlag{
				Name:  flagNoEncrypt,
				Usage: "Force disable secret encryption",
			},
			&cli.BoolFlag{
				Name:    flagVerbose,
				Usage:   "Verbose mode",
				Aliases: []string{"v"},
			},
			&cli.BoolFlag{
				Name:    flagVeryVerbose,
				Usage:   "Very verbose mode",
				Aliases: []string{"vv"},
			},
		},
		Commands: []*cli.Command{
			cmdLogin(),
			cmdRegister(),
			cmdCreateSecretBankCard(),
			cmdCreateSecretCredentials(),
			cmdCreateSecretNote(),
			cmdCreateSecretBlob(),
			cmdRenameSecret(),
			cmdDeleteSecret(),
			cmdAddTag(),
			cmdDeleteTag(),
			cmdEditSecretBankCard(),
			cmdEditSecretCredentials(),
			cmdEditSecretNote(),
			cmdEditSecretBlob(),
			cmdGetSecrets(),
			cmdGetSecret(),
			cmdVersion(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func setup(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	w := cmd.Root().Writer

	logger.SetOutput(w)

	if cmd.Bool(flagVerbose) {
		logger.SetLevel(logrus.DebugLevel)
	}
	if cmd.Bool(flagVeryVerbose) {
		logger.SetLevel(logrus.TraceLevel)
	}

	jwtString, err := authenticate()
	if err != nil && !errors.Is(err, errAuthExpired) {
		return ctx, errors.Wrap(err, "could not authenticate user from local JWT file")
	}

	if jwtString != "" {
		isLoggedIn = true
	}

	address := cmd.String(flagAddress)
	if address == "" {
		return ctx, fmt.Errorf("you haven't provided --%s", flagAddress)
	}

	c = newClient(address, jwtString)

	return ctx, nil
}

func setupAndAuthorize(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	w := cmd.Root().Writer

	ctx, err := setup(ctx, cmd)
	if err != nil {
		return ctx, err
	}

	if !isLoggedIn {
		return ctx, errors.New("you are not authenticated")
	}

	var secretsList []secret

	code, err := SendRequest[[]secret](c, ctx, "/api/secret/list", http.MethodGet, nil, &secretsList)
	if err != nil {
		return ctx, errors.Wrap(err, "could not retrieve secrets list")
	}
	if code != http.StatusOK {
		return ctx, fmt.Errorf("unexpected status code during secrets list retrieval: %d", code)
	}

	for _, item := range secretsList {
		item := item
		secretsByName[item.Name] = &item
		secretsByID[item.ID] = &item
	}

	fmt.Fprintf(w, "Synchronized %d secrets from the server\n\n", len(secretsList))

	return ctx, nil
}

func readPassword(w io.Writer, prompt string) (string, error) {
	fmt.Fprint(w, prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(w, "Could not read from terminal\n")
		return "", nil
	}
	fmt.Fprint(w, "\n")

	return string(password), nil
}
