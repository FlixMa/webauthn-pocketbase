# WebAuthn + Pocketbase Example

This repository contains a proof-of-concept implementation of the user registration and authentication flow described in the webauthn standard.
Useful resources can be found here:

- A Guide to Web Authentication: [webauthn.guide](https://webauthn.guide)
- A demo of the WebAuthn specification: [webauthn.io](https://webauthn.io)
- Web Authentication API: [developer.mozilla.org](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API)

This implementation uses the [Webauthn/FIDO2 library](https://github.com/go-webauthn/webauthn) in golang which is complemented by [some helper functions](https://github.com/github/webauthn-json) on the javascript side until [browser support](https://developer.mozilla.org/en-US/docs/Web/API/PublicKeyCredential/parseCreationOptionsFromJSON_static#browser_compatibility) is widespread.


## Setup

1. Clone the repository
1. Spin up the pocketbase backend:
    ```bash
    cd backend
    go run . serve
    ```
1. Open up a web browser and complete the initial setup at `http://localhost:8090/_/`.
1. Go to the *users* collection and add to fields:
    1. text: `webauthn_id_b64`
    1. json: `webauthn_credentials`
1. (Optional) Go edit the collection and disable all auth methods (password, oauth).
1. In a seperate terminal install the prerequisites for the web frontend
    ```bash
    cd app
    npm install
    ```
1. Spin up a development server for the web app
    ```bash
    npm run dev
    ```
    or build the web app into the *pb_public* directory
    ```bash
    npm run build
    ```
1. Open up a web browser and point it to either `localhost:5173` (dev server) or `localhost:8090` (static) depending on the method you chose.
1. Try it out!


## Try it out

Once everything is setup, try to register a user by entering a username and clicking register.
You'll be prompted to create a some credentials for this webpage. Confirm using your method of choice (e.g. biometrics or physical key).
Now try to login by entering your username and and clicking login. Again, you'll be asked to identify yourself with the previously chosen authenticator.
If everything goes well, the authentication token should have been printed to the console.