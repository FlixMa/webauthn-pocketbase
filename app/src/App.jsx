import { Button, Flex, HStack, Heading, Input, Text } from "@chakra-ui/react"
import { useFetcher, useLoaderData, useRevalidator } from "react-router-dom"
import { create as createCredential, parseCreationOptionsFromJSON, get as getCredential, parseRequestOptionsFromJSON } from "@github/webauthn-json/browser-ponyfill"

import { pocketbase } from "./pocketbase_singleton"
import { useEffect } from "react"

export async function registrationAction({ request }) {
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

export async function loginAction({ request }) {
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

  pocketbase.authStore.save(finalResult.token, finalResult.user)
  return finalResult
}

export async function userLoader() {
  return { user: pocketbase.authStore.model }
}

export default function App() {
  const { user } = useLoaderData()
  const { revalidate } = useRevalidator()
  const registrationFetcher = useFetcher()
  const loginFetcher = useFetcher()
  console.log(user)
  useEffect(() => {
    pocketbase.collection('users').subscribe("*", async (_) => {
      revalidate()
    })
    return () => pocketbase.collection('users').unsubscribe()
  }, [])
  

  return (
    <Flex w={"100vw"} h={"100vh"} p={5} direction={'column'} gap={10}>
      <Heading alignSelf={'center'}>WebAuthn + Pocketbase</Heading>
      <Flex direction={{base: 'column', sm: 'row'}} w={'100%'} gap={5} justify={'space-around'}>
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

      { user &&
      <Flex p={5} direction={'column'} gap={10}>
          <Flex direction={{base: 'column', md: 'row'}}><Text fontWeight={'bold'}>Username:</Text><Text>{user.username}</Text></Flex>
          <Flex direction={{base: 'column', md: 'row'}}><Text fontWeight={'bold'}>Name:</Text><Text>{user.name}</Text></Flex>
          <Flex direction={{base: 'column', md: 'row'}}><Text fontWeight={'bold'}>WebAuthn ID:</Text><Text wordBreak={'break-all'}>{user.webauthn_id_b64}</Text></Flex>
          <Flex direction={{base: 'column', md: 'row'}}><Text fontWeight={'bold'}>WebAuthn Credentials:</Text><Text wordBreak={'break-all'}>{JSON.stringify(user.webauthn_credentials)}</Text></Flex>
          <Button onClick={() => {
            pocketbase.authStore.clear()
            revalidate()
          }}>Logout</Button>
      </Flex>
      }
      
    </Flex>
  )
}
