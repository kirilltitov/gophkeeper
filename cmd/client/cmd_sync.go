package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
)

const (
	secretsByNameFileName = "secrets_by_name"
	secretsByIDFileName   = "secrets_by_id"
)

func getSecretsByNameFileName() string {
	return fmt.Sprintf("%s/%s.json", getConfigDir(), secretsByNameFileName)
}

func getSecretsByIDFileName() string {
	return fmt.Sprintf("%s/%s.json", getConfigDir(), secretsByIDFileName)
}

func storeSecrets(file string, secrets []byte) error {
	if err := os.WriteFile(file, secrets, 0o660); err != nil {
		return err
	}

	return nil
}

func cmdSync() *cli.Command {
	return &cli.Command{
		Name:        "sync",
		Description: "Performs explicit sync of all user's secrets from server",
		Before:      setupAndAuthorize,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := syncSecrets(ctx); err != nil {
				return err
			}

			fmt.Fprintf(cmd.Root().Writer, "Synchronized %d secrets from the server\n\n", len(secretsByID))

			return nil
		},
	}
}

func syncSecrets(ctx context.Context) error {
	var err error
	var bytes []byte

	var secretsList []*secret

	code, err := SendRequest[[]*secret](c, ctx, "/api/secret/list", http.MethodGet, nil, &secretsList)
	if err != nil {
		return errors.Wrap(err, "could not retrieve secrets list")
	}
	if code != http.StatusOK {
		return fmt.Errorf("unexpected status code during secrets list retrieval: %d", code)
	}

	for _, item := range secretsList {
		secretsByName[item.Name] = item
		secretsByID[item.ID] = item
	}

	bytes, err = json.Marshal(secretsByName)
	if err != nil {
		return errors.Wrap(err, "could not marshal secrets to json")
	}
	if err := storeSecrets(getSecretsByNameFileName(), bytes); err != nil {
		return errors.Wrap(err, "could not save secrets to local file")
	}

	bytes, err = json.Marshal(secretsByID)
	if err != nil {
		return errors.Wrap(err, "could not marshal secrets to json")
	}
	if err := storeSecrets(getSecretsByIDFileName(), bytes); err != nil {
		return errors.Wrap(err, "could not save secrets to local file")
	}

	return nil
}

func loadLocalSecrets() error {
	localSecretsByName, err := loadLocalSecretsFile[map[string]*secret](getSecretsByNameFileName())
	if err != nil {
		return err
	}
	localSecretsByID, err := loadLocalSecretsFile[map[uuid.UUID]*secret](getSecretsByIDFileName())
	if err != nil {
		return err
	}

	secretsByName = *localSecretsByName
	secretsByID = *localSecretsByID

	return nil
}

func loadLocalSecretsFile[R any](fileName string) (*R, error) {
	var err error
	var bytes []byte

	bytes, err = os.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not read local secrets file '%s'", fileName))
	}

	var result R
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal local secrets file")
	}

	return &result, nil
}
