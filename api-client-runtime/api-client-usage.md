## API Client Usage (template)

This snippet explains how to wire a generated API client together with the shared runtime when `ae-cli` scaffolds a project.

1. Install dependencies in the generated app:

```bash
pnpm add @agile-exec/api-client-runtime @agile-exec/api-client-<service>
```

2. Initialize the runtime in your app entry (e.g., `src/main.ts`):

```ts
import { createRuntime } from '@agile-exec/api-client-runtime'
import { ApiFactory } from '@agile-exec/api-client-<service>'

createRuntime({ baseURL: import.meta.env.VITE_API_BASE_URL, getToken: () => localStorage.getItem('auth_token') })

// Example: create a typed API client instance (depends on generator)
const api = ApiFactory({ /* options */ })
export default api
```

3. Use in stores/components:

```ts
import api from '@/api'
await api.auth.login({ email, password })
```

Place this file in the `ae-cli/templates` folder so `ae-cli` can include it in scaffolded projects or print it as a usage hint after generation.
