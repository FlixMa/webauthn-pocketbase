package main

import (
	"fmt"

	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	webAuthn         *webauthn.WebAuthn
	err              error
	webAuthnSessions map[string]*webauthn.SessionData
)

func findUser(app *pocketbase.PocketBase, username string) (User, error) {
	// Find user
	user := User{}
	err := app.Dao().DB().
		NewQuery("SELECT id, username, name, webauthn_credentials FROM users WHERE username={:username}").
		Bind(dbx.Params{"username": username}).
		One(&user)

	return user, err
}

func findOrCreateUser(app *pocketbase.PocketBase, username string) (User, error) {

	user, err := findUser(app, username)
	if err != nil {
		// TODO: create user if not existent
	}

	return user, err
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

		e.Router.POST("/webauthn-begin-registration", func(c echo.Context) error {
			username := c.FormValue("username")

			// Find or create the new user
			user, err := findOrCreateUser(app, username)
			if err != nil {
				return err
			}

			options, session, err := webAuthn.BeginRegistration(user)
			if err != nil {
				return err
			}

			// store the sessionData values
			webAuthnSessions[user.WebAuthnId] = session

			// return the options generated
			// -> options.publicKey contain our registration options
			return c.JSON(http.StatusOK, options)
		})

		e.Router.POST("/webauthn-finish-registration", func(c echo.Context) error {
			//username := c.FormValue("username")
			//fmt.Println(c.FormValues())
			//fmt.Println(c.QueryParams())

			// Find or create the new user
			user, err := findOrCreateUser(app, "felix")
			if err != nil {
				return err
			}
			session := webAuthnSessions[user.WebAuthnId]

			credential, err := webAuthn.FinishRegistration(user, *session, c.Request())
			if err != nil {
				// Handle Error and return.
				return err
			}

			// If creation was successful, store the credential object
			err = user.addCredential(app, *credential)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]string{"status": "success"})
		})

		e.Router.POST("/webauthn-begin-login", func(c echo.Context) error {
			username := c.QueryParam("username")
			user, err := findUser(app, username)
			if err != nil {
				return err
			}
			fmt.Printf("user.WebAuthnCredentialsJSON: %T %v\n", user.WebAuthnCredentialsJSON, user.WebAuthnCredentialsJSON)

			options, session, err := webAuthn.BeginLogin(user)
			if err != nil {
				// Handle Error and return.
				return err
			}

			// store the session values
			webAuthnSessions[user.WebAuthnId] = session

			// return the options generated
			// options.publicKey contain our registration options
			return c.JSON(http.StatusOK, options)

		})

		e.Router.POST("/webauthn-finish-login", func(c echo.Context) error {
			//username := c.QueryParam("username")
			user, err := findOrCreateUser(app, "felix")
			if err != nil {
				return err
			}

			// Get the session data stored from the function above
			session := webAuthnSessions[user.WebAuthnId]

			credential, err := webAuthn.FinishLogin(user, *session, c.Request())
			if err != nil {
				// Handle Error and return.

				return err
			}
			fmt.Println("Finish-Login-Credential:")
			fmt.Println(credential)

			// If login was successful, handle next steps
			return user.sendAuthTokenResponse(app, c)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
