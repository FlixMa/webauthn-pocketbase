import PocketBase from 'pocketbase'

const isDevelopmentEnvironment = (!process.env.NODE_ENV || process.env.NODE_ENV === 'development')
const baseURL = isDevelopmentEnvironment ? 'http://localhost:8090' : window.location.origin

export const pocketbase = new PocketBase(baseURL)