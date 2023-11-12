package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id                      string `db:"id" json:"id"`
	Username                string `db:"username" json:"username"`
	Name                    string `db:"name" json:"name"`
	WebAuthnIdB64           string `db:"webauthn_id_b64" json:"webauthn_id_b64"`
	WebAuthnCredentialsJSON string `db:"webauthn_credentials" json:"webauthn_credentials"`
}

// WebAuthnID provides the user handle of the user account. A user handle is an opaque byte sequence with a maximum
// size of 64 bytes, and is not meant to be displayed to the user.
//
// To ensure secure operation, authentication and authorization decisions MUST be made on the basis of this id
// member, not the displayName nor name members. See Section 6.1 of [RFC8266].
//
// It's recommended this value is completely random and uses the entire 64 bytes.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dom-publickeycredentialuserentity-id)
func (user User) WebAuthnID() []byte {
	webAuthnId, err := base64.StdEncoding.DecodeString(user.WebAuthnIdB64)
	if err != nil {
		fmt.Printf("Could not base64 decode WebAuthnID from database err: %v (base64 id: %v)\n", err, user.WebAuthnIdB64)
		return []byte{}
	}
	return webAuthnId
}

// WebAuthnName provides the name attribute of the user account during registration and is a human-palatable name for the user
// account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party SHOULD let the user
// choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dictdef-publickeycredentialuserentity)
func (user User) WebAuthnName() string {
	return user.Username
}

// WebAuthnDisplayName provides the name attribute of the user account during registration and is a human-palatable
// name for the user account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party
// SHOULD let the user choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://www.w3.org/TR/webauthn/#dom-publickeycredentialuserentity-displayname)
func (user User) WebAuthnDisplayName() string {
	if user.Name == "" {
		return user.WebAuthnName()
	}
	return user.Name
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (user User) WebAuthnCredentials() []webauthn.Credential {
	// decode string from database into credential object
	credential := webauthn.Credential{}
	err := json.Unmarshal([]byte(user.WebAuthnCredentialsJSON), &credential)
	if err != nil {
		fmt.Printf("error while unmarshalling credential from db: %v\n", err)
		return []webauthn.Credential{}
	}

	// NOTE: database currently only stores a single credential per user.
	credentials := []webauthn.Credential{credential}
	return credentials
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (u User) WebAuthnIcon() string {
	return ""
}
