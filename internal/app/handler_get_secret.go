package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// HandlerGetSecret retrieves secret with value.
//
// Example request:
//
// GET /api/secret/{ID}
//
// Example response:
//
//	{
//	  "success": true,
//	  "result": {
//	    "id":           "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	    "user_id":      "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	    "name":         "my secret card",
//	    "tags":         ["MIR"],
//	    "kind":         "bank_card",
//	    "is_encrypted": true,
//	    "value": {
//	      "id":     "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	      "name":   "KIRILL TITOV",
//	      "number": "1234 5678 9012 3456",
//	      "date":   "12/34/5678",
//	      "cvv":    "322"
//	    }
//	  },
//	  "error": null
//	}
//
//	{
//	  "success": true,
//	  "result": {
//	    "id":           "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	    "user_id":      "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	    "name":         "my secret credentials",
//	    "tags":         [],
//	    "kind":         "credentials",
//	    "is_encrypted": true,
//	    "value": {
//	      "id":       "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	      "login":    "frank.strino",
//	      "password": "secret password",
//	    }
//	  },
//	  "error": null
//	}
//
//	{
//	  "success": true,
//	  "result": {
//	    "id":           "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	    "user_id":      "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	    "name":         "my secret note",
//	    "tags":         ["notes", "secret notes"],
//	    "kind":         "note",
//	    "is_encrypted": true,
//	    "value": {
//	      "id":   "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	      "body": "secret body",
//	    }
//	  },
//	  "error": null
//	}
//
//	{
//	  "success": true,
//	  "result": {
//	    "id":           "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	    "user_id":      "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	    "name":         "my secret blob",
//	    "tags":         [],
//	    "kind":         "blob",
//	    "is_encrypted": true,
//	    "value": {
//	      "id":   "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	      "body": "0JAg0LXRidC1INGPINC/0LjRiNGDINC80YPQt9GL0LrRgyA6KSBodHRwczovL2NsY2sucnUvM0doZW5B",
//	    }
//	  },
//	  "error": null
//	}
//
// May response with codes 200, 401, 500.
func (a *Application) HandlerGetSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	secretID, err := getUUIDFromRequest(r, "ID")
	if err != nil {
		returnErrorWithCode(w, http.StatusBadRequest, err.Error())
		return
	}

	secret, err := a.Gophkeeper.GetSecretWithValueByID(ctx, *secretID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, gophkeeper.ErrNoAuth) {
			code = http.StatusUnauthorized
		}
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode(w, http.StatusOK, secret)
}
