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

var (
	webAuthn         *webauthn.WebAuthn
	err              error
	webAuthnSessions map[string]*webauthn.SessionData
)

func findUser(app *pocketbase.PocketBase, username string) (*User, error) {
	// Find user
	user := User{}
	err := app.Dao().DB().
		NewQuery("SELECT id, username, name, webauthn_id_b64, webauthn_credentials FROM users WHERE username={:username}").
		Bind(dbx.Params{"username": username}).
		One(&user)
	if err != nil {
		return nil, err
	}

	err = user.ensureWebAuthnId(app)
	return &user, err
}

func createUser(app *pocketbase.PocketBase, username string) error {

	collection, err := app.Dao().FindCollectionByNameOrId("users")
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
		"username":             username,
		"password":             password,
		"passwordConfirm":      password,
		"webauthn_id_b64":      "",
		"webauthn_credentials": "",
	})
	if err != nil {
		return err
	}

	err = form.Submit()
	return err
}

func findOrCreateUser(app *pocketbase.PocketBase, username string) (*User, error) {

	user, err := findUser(app, username)

	// create user if not existent
	if err != nil {
		err = createUser(app, username)
		if err != nil {
			return user, err
		}
		user, err = findUser(app, username)
		if err != nil {
			return user, err
		}
	}

	err = user.ensureWebAuthnId(app)
	return user, err
}

func (user *User) ensureWebAuthnId(app *pocketbase.PocketBase) error {
	// create webauthn id only if it doesnt exist yet
	if user.WebAuthnIdB64 != "" {
		return nil
	}

	authRecord, err := app.Dao().FindAuthRecordByUsername("users", user.Username)
	if err != nil {
		return apis.NewBadRequestError("user not found.", err)
	}

	// create 64 bytes of random data
	randomBuffer := make([]byte, 64)
	rand.Read(randomBuffer)
	user.WebAuthnIdB64 = base64.StdEncoding.EncodeToString(randomBuffer)

	// store in database
	authRecord.Set("webauthn_id_b64", user.WebAuthnIdB64)
	err = app.Dao().Save(authRecord)
	if err != nil {
		return apis.NewBadRequestError("could not store webauthn id to db.", err)
	}

	return nil
}

func (user User) addCredential(app *pocketbase.PocketBase, credential webauthn.Credential) error {
	authRecord, err := app.Dao().FindAuthRecordByUsername("users", user.Username)
	if err != nil {
		return apis.NewBadRequestError("user not found.", err)
	}
	authRecord.Set("webauthn_credentials", credential)
	return app.Dao().Save(authRecord)
}

func (user User) sendAuthTokenResponse(app *pocketbase.PocketBase, c echo.Context) error {

	authRecord, err := app.Dao().FindAuthRecordByUsername("users", user.Username)
	if err != nil {
		return apis.NewBadRequestError("User not found.", err)
	}
	token, err := tokens.NewRecordAuthToken(app, authRecord)
	if err != nil {
		return apis.NewBadRequestError("Failed to create auth token.", err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user":  authRecord,
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////        MAIN        //////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	wconfig := &webauthn.Config{
		RPDisplayName: "Felix' PB Webauthn",                                       // Display Name for your site
		RPID:          "localhost",                                                // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:8090", "http://localhost:5173"}, // The origin URLs allowed for WebAuthn requests
	}

	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
	}
	webAuthnSessions = make(map[string]*webauthn.SessionData)

	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		e.Router.POST("/webauthn-begin-registration/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			data := apis.RequestInfo(c).Data
			fmt.Printf("/webauthn-begin-registration data: %v\n", data)

			// Find or create the new user
			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewBadRequestError("User not found.", err)
			}

			options, session, err := webAuthn.BeginRegistration(user)
			if err != nil {
				return err
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

			data := apis.RequestInfo(c).Data
			fmt.Printf("/webauthn-finish-registration data: %v\n", data)

			// Find or create the new user
			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewBadRequestError("User not found.", err)
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

		e.Router.POST("/webauthn-begin-login/:userb64", func(c echo.Context) error {
			usernameb64 := c.PathParam("userb64")
			username, err := base64.StdEncoding.DecodeString(usernameb64)
			if err != nil {
				return apis.NewBadRequestError("Could not decode user from path.", err)
			}

			data := apis.RequestInfo(c).Data
			fmt.Printf("/webauthn-begin-login data: %v\n", data)

			user, err := findUser(app, string(username))
			if err != nil {
				return apis.NewBadRequestError("User not found.", err)
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

			data := apis.RequestInfo(c).Data
			fmt.Printf("/webauthn-finish-login data: %v\n", data)

			user, err := findOrCreateUser(app, string(username))
			if err != nil {
				return apis.NewBadRequestError("User not found.", err)
			}

			// Get the session data stored from the function above
			session := webAuthnSessions[user.WebAuthnIdB64]

			_, err = webAuthn.FinishLogin(user, *session, c.Request())
			if err != nil {
				// Handle Error and return.

				return apis.NewBadRequestError("Failed to verify login credentials.", err)
			}

			// If login was successful, handle next steps
			return user.sendAuthTokenResponse(app, c)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
