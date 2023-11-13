import React from 'react'
import ReactDOM from 'react-dom/client'
import { ChakraProvider, ColorModeScript, extendTheme } from '@chakra-ui/react'
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import App, { registrationAction, loginAction, userLoader } from './App.jsx'
import theme, { config as themeConfig } from './theme.js'

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    loader: userLoader,
    children: [
      {
        path: 'webauthn-registration',
        action: registrationAction
      },
      {
        path: 'webauthn-login',
        action: loginAction
      },
    ]
  }
])

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ColorModeScript initialColorMode={themeConfig.initialColorMode} />
    <ChakraProvider theme={theme}>
      <RouterProvider router={router} />
    </ChakraProvider>
  </React.StrictMode>,
)
