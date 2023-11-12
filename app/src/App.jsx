import { Button, Center, Flex, Input, Text, VStack } from "@chakra-ui/react"
import { useFetcher } from "react-router-dom"
import { create as createCredential, parseCreationOptionsFromJSON, get as getCredential, parseRequestOptionsFromJSON } from "@github/webauthn-json/browser-ponyfill"

import { pocketbase } from "./pocketbase_singleton"

export async function registrationAction({ params, request }) {
  const data = Object.fromEntries(await request.formData())
  console.log("beginRegistrationAction data", data)
  const publicKeyCredentialCreationOptions = await pocketbase.send("/webauthn-begin-registration", {
    method: "POST",
    query: data
  })
  console.log("publicKeyCredentialCreationOptions", publicKeyCredentialCreationOptions)

  const credential = await createCredential(parseCreationOptionsFromJSON(publicKeyCredentialCreationOptions))
  console.log("finishRegistration: send credential", credential)

  const finalResult = await pocketbase.send("/webauthn-finish-registration", {
    method: "POST",
    //query: data,
    body: credential
  })
  console.log("beginRegistrationAction finalResult", finalResult)

  return finalResult
}

export async function loginAction({ params, request }) {
  const data = Object.fromEntries(await request.formData())
  console.log("loginAction data", data)

  const publicKeyCredentialRequestOptions = await pocketbase.send("/webauthn-begin-login", {
    method: "POST",
    query: data
  })
  console.log("publicKeyCredentialRequestOptions", publicKeyCredentialRequestOptions)
  
  const assertion = await getCredential(parseRequestOptionsFromJSON(publicKeyCredentialRequestOptions))
  console.log("finishLogin: send assertion", assertion)

  const finalResult = await pocketbase.send("/webauthn-finish-login", {
    method: "POST",
    //query: data,
    body: assertion
  })
  console.log("beginRegistrationAction finalResult", finalResult)

  return finalResult
}

export default function App() {
  const registrationFetcher = useFetcher()
  const loginFetcher = useFetcher()

  return (
    <Center p={5}>
      <VStack spacing={5}>
        <Text>Hello :-)</Text>
        <Flex gap={5}>
          <Flex direction="column">
            <registrationFetcher.Form method="post" action="/webauthn-registration">
              <Text>Register :-)</Text>
              <Input name="username" type="text" />
              <Button type="submit">Register</Button>
            </registrationFetcher.Form>
          </Flex>
          <Flex direction="column">
            <loginFetcher.Form method="post" action="/webauthn-login">
              <Text>Login :-)</Text>
              <Input name="username" type="text" />
              <Button type="submit">Login</Button>
            </loginFetcher.Form>
          </Flex>
        </Flex>
      </VStack>
    </Center>
  )
}
