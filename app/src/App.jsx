import { Button, Flex, Heading, Input } from "@chakra-ui/react"
import { useFetcher } from "react-router-dom"
import { create as createCredential, parseCreationOptionsFromJSON, get as getCredential, parseRequestOptionsFromJSON } from "@github/webauthn-json/browser-ponyfill"

import { pocketbase } from "./pocketbase_singleton"

export async function registrationAction({ params, request }) {
  const data = Object.fromEntries(await request.formData())
  console.log("beginRegistrationAction data", data)
  const publicKeyCredentialCreationOptions = await pocketbase.send(`/webauthn-begin-registration/${btoa(data.username)}`, {
    method: "POST"
  })
  console.log("publicKeyCredentialCreationOptions", publicKeyCredentialCreationOptions)

  const credential = await createCredential(parseCreationOptionsFromJSON(publicKeyCredentialCreationOptions))
  console.log("finishRegistration: send credential", credential.toJSON())

  const finalResult = await pocketbase.send(`/webauthn-finish-registration/${btoa(data.username)}`, {
    method: "POST",
    body: credential
  })
  console.log("beginRegistrationAction finalResult", finalResult)

  return finalResult
}

export async function loginAction({ params, request }) {
  const data = Object.fromEntries(await request.formData())
  console.log("loginAction data", data)

  const publicKeyCredentialRequestOptions = await pocketbase.send(`/webauthn-begin-login/${btoa(data.username)}`, {
    method: "POST"
  })
  console.log("publicKeyCredentialRequestOptions", publicKeyCredentialRequestOptions)
  
  const assertion = await getCredential(parseRequestOptionsFromJSON(publicKeyCredentialRequestOptions))
  console.log("finishLogin: send assertion", assertion)

  const finalResult = await pocketbase.send(`/webauthn-finish-login/${btoa(data.username)}`, {
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
    <Flex w={"100vw"} h={"100vh"} p={5} direction={'column'} gap={10}>
      <Heading alignSelf={'center'}>WebAuthn + Pocketbase</Heading>
      <Flex w={'100%'} gap={5} justify={'space-around'}>
          
          <registrationFetcher.Form method="post" action="/webauthn-registration">
            <Flex direction="column" p={10} gap={5} bg={'gray.700'} borderRadius={'md'}>
              <Heading size={'sm'}>Registration:</Heading>
              <Input name="username" type="text" placeholder="Username" />
              <Button type="submit">Register</Button>
            </Flex>
          </registrationFetcher.Form>
        
          <loginFetcher.Form method="post" action="/webauthn-login">
            <Flex direction="column" p={10} gap={5} bg={'gray.700'} borderRadius={'md'}>
              <Heading size={'sm'}>Login:</Heading>
              <Input name="username" type="text" placeholder="Username" />
              <Button type="submit">Login</Button>
            </Flex>
          </loginFetcher.Form>
        
        </Flex>
    </Flex>
  )
}
