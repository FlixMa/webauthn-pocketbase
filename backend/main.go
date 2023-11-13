package main

import (
	"fmt"

	"log"
	"net/http"
	"os"

	"crypto/rand"
	"encoding/base64"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"

	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	WEBAUTHN_COLLECTION_NAME       string = "users"
	WEBAUTHN_CREDENTIALS_FIELDNAME string = "webauthn_credentials"
	WEBAUTHN_ID_B64_FIELDNAME      string = "webauthn_id_b64"
)

func findUser(app *pocketbase.PocketBase, username string) (*User, error) {
	// Find user
	user := User{}
	err := app.Dao().DB().
		NewQuery(fmt.Sprintf(
			"SELECT id, username, name, %s, %s FROM %s WHERE username={:username}",
			WEBAUTHN_ID_B64_FIELDNAME, WEBAUTHN_CREDENTIALS_FIELDNAME, WEBAUTHN_COLLECTION_NAME)).
		Bind(dbx.Params{"username": username}).
		One(&user)
	if err != nil {
		return nil, err
	}

	err = user.ensureWebAuthnId(app)
	return &user, err
}

func createUser(app *pocketbase.PocketBase, username string) error {

	collection, err := app.Dao().FindCollectionByNameOrId(WEBAUTHN_COLLECTION_NAME)
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)

	// create long random password (NOTE: password auth is disabled anyway)
	randomBuffer := make([]byte, 32)
	rand.Read(randomBuffer)
	password := base64.StdEncoding.EncodeToString(randomBuffer)

	err = form.LoadData(map[string]any{
		"username":                     username,
		"password":                     password,
		"passwordConfirm":              password,
		WEBAUTHN_ID_B64_FIELDNAME:      "",
		WEBAUTHN_CREDENTIALS_FIELDNAME: "",
	})
	if err != nil {
		return err
	}

	err = form.Submit()
	return err
}

func findOrCreateUser(app *pocketbase.PocketBase, username string) (*User, error) {

	// try to find user
	user, err := findUser(app, username)

	// error? -> create user if not existent
	if err != nil {
		err = createUser(app, username)
		if err != nil {
			return nil, err
		}
		user, err = findUser(app, username)
		if err != nil {
			return nil, err
		}
	}

	err = user.ensureWebAuthnId(app)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (user *User) ensureWebAuthnId(app *pocketbase.PocketBase) error {
	authRecord, err := app.Dao().FindAuthRecordByUsername(WEBAUTHN_COLLECTION_NAME, user.Username)
	if err != nil {
		return apis.NewNotFoundError("User not found.", err)
	}

	// create webauthn id only if it doesnt exist yet
	if authRecord.GetString(WEBAUTHN_ID_B64_FIELDNAME) != "" {
		return nil
	}

	// create 64 bytes of random data
	randomBuffer := make([]byte, 64)
	rand.Read(randomBuffer)
	user.WebAuthnIdB64 = base64.StdEncoding.EncodeToString(randomBuffer)

	// store in database
	authRecord.Set(WEBAUTHN_ID_B64_FIELDNAME, user.WebAuthnIdB64)
	err = app.Dao().Save(authRecord)
	if err != nil {
		return NewInternalServerError("Could not store webauthn id to db.", err)
	}

	return nil
}

func (user User) addCredential(app *pocketbase.PocketBase, credential webauthn.Credential) error {
	authRecord, err := app.Dao().FindAuthRecordByUsername(WEBAUTHN_COLLECTION_NAME, user.Username)
	if err != nil {
		return apis.NewNotFoundError("User not found.", err)
	}
	authRecord.Set(WEBAUTHN_CREDENTIALS_FIELDNAME, credential)
	return app.Dao().Save(authRecord)
}

func (user User) sendAuthTokenResponse(app *pocketbase.PocketBase, c echo.Context) error {

	authRecord, err := app.Dao().FindAuthRecordByUsername(WEBAUTHN_COLLECTION_NAME, user.Username)
	if err != nil {
		return apis.NewNotFoundError("User not found.", err)
	}
	token, err := tokens.NewRecordAuthToken(app, authRecord)
	if err != nil {
		return NewInternalServerError("Failed to create auth token.", err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user":  authRecord,
	})
}

func NewInternalServerError(message string, data any) *apis.ApiError {
	return apis.NewApiError(500, message, data)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////        MAIN        //////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	// configure and initialize webauthn
	wconfig := &webauthn.Config{
		RPDisplayName: "Felix' PB Webauthn",                                       // Display Name for your site
		RPID:          "localhost",                                                // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:8090", "http://localhost:5173"}, // The origin URLs allowed for WebAuthn requests
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create a map holding the sessions used during registration and login flow
	webAuthnSessions := make(map[string]*webauthn.SessionData)

	// Configure the pocketbase server
	app := pocketbase.New()
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		// register routes for registration:
		// 	1. *begin* creates a challenge for the authenticator to sign
		// 	2. *finish* validates the generated credential and stores it in the user record
		e.Router.POST("/webauthn-begin-registration/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			// Find or create the new user
			// TODO: if the user exists, throw an error (cant register a credential for existing users, if not authenticated)
			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewNotFoundError("User not found.", err)
			}

			options, session, err := webAuthn.BeginRegistration(user)
			if err != nil {
				return NewInternalServerError("Could not start WebAuthn registration flow.", err)
			}

			// store the sessionData values
			webAuthnSessions[user.WebAuthnIdB64] = session

			// return the options generated
			// -> options.publicKey contain our registration options
			return c.JSON(http.StatusOK, options)
		})

		e.Router.POST("/webauthn-finish-registration/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			// Find or create the new user
			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewNotFoundError("User not found.", err)
			}
			session := webAuthnSessions[user.WebAuthnIdB64]
			delete(webAuthnSessions, user.WebAuthnIdB64)

			credential, err := webAuthn.FinishRegistration(user, *session, c.Request())
			if err != nil {
				// Handle Error and return.
				return apis.NewBadRequestError("Failed to verify login credentials.", err)
			}

			// If creation was successful, store the credential object
			err = user.addCredential(app, *credential)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]string{"status": "success"})
		})

		// register routes for authentication (similar to registration):
		// 	1. *begin* creates a challenge for the authenticator to sign
		// 	2. *finish* validates the signature and responds with a pocketbase auth token
		e.Router.POST("/webauthn-begin-login/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			user, err := findUser(app, string(username))
			if err != nil {
				return apis.NewNotFoundError("User not found.", err)
			}

			options, session, err := webAuthn.BeginLogin(user)
			if err != nil {
				// Handle Error and return.
				return err
			}

			// store the session values
			webAuthnSessions[user.WebAuthnIdB64] = session

			// return the options generated
			// options.publicKey contain our registration options
			return c.JSON(http.StatusOK, options)

		})

		e.Router.POST("/webauthn-finish-login/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewNotFoundError("User not found.", err)
			}

			// Get the session data stored from the function above and complete login
			session := webAuthnSessions[user.WebAuthnIdB64]
			_, err = webAuthn.FinishLogin(user, *session, c.Request())
			if err != nil {
				return apis.NewBadRequestError("Failed to verify login credentials.", err)
			}

			// If login was successful, send auth token for further communication
			return user.sendAuthTokenResponse(app, c)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
