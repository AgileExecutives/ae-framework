import axios, { AxiosInstance } from 'axios'

export interface RuntimeOptions {
  baseURL?: string
  getToken?: () => string | null
}

let instance: AxiosInstance | null = null

export function createRuntime(options: RuntimeOptions = {}) {
  instance = axios.create({ baseURL: options.baseURL })
  instance.interceptors.request.use((config) => {
    const token = options.getToken?.()
    if (token) config.headers = { ...(config.headers ?? {}), Authorization: `Bearer ${token}` }
    return config
  })
  return instance
}

export function getInstance(): AxiosInstance {
  if (!instance) throw new Error('Runtime axios instance not initialized. Call createRuntime() first.')
  return instance
}

export function setAuthToken(token: string | null) {
  if (!instance) return
  if (token) instance.defaults.headers = { ...(instance.defaults.headers ?? {}), Authorization: `Bearer ${token}` }
  else if (instance.defaults.headers) delete (instance.defaults.headers as any).Authorization
}

export default {
  createRuntime,
  getInstance,
  setAuthToken
}
