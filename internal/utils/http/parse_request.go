package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func ParseRequest(w http.ResponseWriter, r *http.Request, target any) error {
	var buf bytes.Buffer
	defer r.Body.Close()

	if _, err := buf.ReadFrom(r.Body); err != nil {
		utils.Log.Infof("Could not get body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &target); err != nil {
		utils.Log.Infof("Could not parse request JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	return nil
}
