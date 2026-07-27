// Types for the auth store
export interface User {
  id: number
  username: string
  email: string
  first_name?: string
  last_name?: string
  role: string
  tenant_id: number
  active: boolean
  created_at: string
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface RegisterCredentials {
  username: string
  email: string
  password: string
  first_name: string
  last_name: string
  role?: string
  tenant_id?: number
  accept_terms: boolean
  newsletter_opt_in: boolean
  company_name?: string
}

export interface ChangePasswordCredentials {
  current_password: string
  new_password: string
}