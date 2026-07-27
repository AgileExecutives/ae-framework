# @agile-exec/api-client-runtime

Shared runtime utilities for generated API clients used by AE frontend apps.

Exports:
- `createRuntime({ baseURL, getToken })` — creates and configures an axios instance used by generated clients.
- `getInstance()` — returns the initialized axios instance.
- `setAuthToken(token)` — set or clear default Authorization header.

Usage example

```ts
import { createRuntime, setAuthToken } from '@agile-exec/api-client-runtime'

createRuntime({ baseURL: process.env.API_BASE_URL, getToken: () => localStorage.getItem('auth_token') })
setAuthToken('my-token')
```
