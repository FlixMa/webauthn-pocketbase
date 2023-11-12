package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	Id                      string `db:"id" json:"id"`
	Username                string `db:"username" json:"username"`
	Name                    string `db:"name" json:"name"`
	WebAuthnId              string `db:"webauthn_id" json:"webauthn_id"`
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
func (u User) WebAuthnID() []byte {
	userHash := sha256.Sum256([]byte(u.Id))
	return userHash[:]
}

// WebAuthnName provides the name attribute of the user account during registration and is a human-palatable name for the user
// account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party SHOULD let the user
// choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dictdef-publickeycredentialuserentity)
func (u User) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName provides the name attribute of the user account during registration and is a human-palatable
// name for the user account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party
// SHOULD let the user choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://www.w3.org/TR/webauthn/#dom-publickeycredentialuserentity-displayname)
func (u User) WebAuthnDisplayName() string {
	return u.Name
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (user User) WebAuthnCredentials() []webauthn.Credential {
	/* user.WebAuthnCredentialsJSON: string
	{
		"ID": "7uYw/HdkaM/LRKUt10Cij109NVA=",
		"PublicKey": "pQECAyYgASFYIGWuzTIC8JPXKzGWWjUE7298wsJ8LSaoo7EfBTmEqLJiIlgg5GTk3v7PgDbv/ib3+CuV2q66+Ctl+WIZu3dF/I7tGmc=",
		"AttestationType": "none",
		"Transport": [
			"internal",
			"hybrid"
		],
		"Flags": {
			"UserPresent": true,
			"UserVerified": true,
			"BackupEligible": true,
			"BackupState": true
		},
		"Authenticator": {
			"AAGUID": "AAAAAAAAAAAAAAAAAAAAAA==",
			"SignCount": 0,
			"CloneWarning": false,
			"Attachment": "platform"
		}
		}
	}
	*/

	fmt.Printf("user.WebAuthnCredentialsJSON: %T %v\n", user.WebAuthnCredentialsJSON, user.WebAuthnCredentialsJSON)

	credential := webauthn.Credential{}
	err := json.Unmarshal([]byte(user.WebAuthnCredentialsJSON), &credential)
	if err != nil {
		fmt.Printf("unmarshal credential from db: err: %v\n", err)
		return []webauthn.Credential{}
	}
	fmt.Printf("credential from db: %v\n", credential)

	credentials := []webauthn.Credential{credential}
	return credentials
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (u User) WebAuthnIcon() string {
	return ""
}
