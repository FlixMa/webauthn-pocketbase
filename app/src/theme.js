import { extendTheme } from "@chakra-ui/react"

export const config = {
  initialColorMode: 'dark',
  useSystemColorMode: false,
}

const theme = extendTheme({ config })
export default theme
