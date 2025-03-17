package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// HandlerGetSecrets retrieves all secrets for current user.
//
// Example request:
//
// GET /api/secret/list
//
// Example response:
//
//	{
//	  "success": true,
//	  "result": [
//	    {
//	      "id": "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	      "user_id": "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	      "name": "foo",
//	      "description": "my secret description",
//	      "tags": ["bar","baz"],
//	      "kind": "note",
//	      "is_encrypted": false,
//	      "value": {
//	        "id": "1ee1416c-d537-6ae0-b6c7-0f48c8929427",
//	        "body": "foo body"
//	      }
//	    },
//	    {
//	      "id": "1ee1416c-d537-6ae0-b6c7-0f48c8929428",
//	      "user_id": "1ee06239-36d2-6142-b86b-55c4f2f680df",
//	      "name": "foo",
//	      "description": "my secret description",
//	      "tags": [],
//	      "kind": "credentials",
//	      "is_encrypted": false,
//	      "value": {
//	        "id": "1ee1416c-d537-6ae0-b6c7-0f48c8929428",
//	        "url": "https://passport.yandex.ru/",
//	        "login": "teonoman",
//	        "password": "megapass"
//	      }
//	    }
//	  ],
//	  "error": null
//	}
//
// May response with codes 200, 401, 500.
func (a *Application) HandlerGetSecrets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	secrets, err := a.Gophkeeper.GetSecrets(ctx)
	if err != nil {
		code := http.StatusInternalServerError
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode(w, http.StatusOK, &secrets)
}
