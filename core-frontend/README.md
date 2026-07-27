# Base App

Vue 3 + TypeScript + Vite base application with authentication and password management.

## Features

### Authentication
- ✅ User registration
- ✅ User login with JWT tokens
- ✅ Token persistence (localStorage)
- ✅ Protected routes
- ✅ Authentication middleware

### Password Management
- ✅ **Forgot Password** - Request password reset via email
- ✅ **Reset Password** - Reset password using token from email
- ✅ **Change Password** - Change password when authenticated
- ✅ Dynamic password requirements from server
- ✅ Real-time validation with Zod
- ✅ Mock email mode for development

### UI Components
- Built with Reka UI (headless components)
- Styled with Tailwind CSS
- Form validation with vee-validate + Zod
- Accessible and responsive design

## Documentation

- [Environment Setup](./ENVIRONMENT_SETUP.md) - **START HERE** - Configure environment variables
- [E2E Test Guide](./E2E_TEST_GUIDE.md) - Complete E2E testing documentation
- [Password Management Summary](./PASSWORD_MANAGEMENT_SUMMARY.md) - Implementation details
- [Registration Form Guide](./REGISTRATION_FORM_GUIDE.md) - Registration form documentation
- [Project Structure](#project-structure)

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd) 
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

### 1. Install Dependencies

```sh
yarn
```

### 2. Configure Environment Variables

Copy the example environment file and configure it:

```sh
cp .env.example .env
```

Edit `.env` to set your API URL and other settings. See [Environment Setup](./ENVIRONMENT_SETUP.md) for details.

**Required Variables:**
- `VITE_API_BASE_URL` - Backend API URL (default: `http://localhost:8080/api/v1`)
- `VITE_APP_URL` - Frontend URL (default: `http://localhost:3003`)

### Compile and Hot-Reload for Development

```sh
yarn dev
```

### Type-Check, Compile and Minify for Production

```sh
yarn build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
yarn test:unit
```

### Run End-to-End Tests with [Playwright](https://playwright.dev)

#### Quick Start (Recommended)
```sh
# Run all E2E tests with automated setup
./run-e2e-tests.sh
```

The script will automatically:
- Check if backend server is running
- Start dev server if needed
- Install Playwright browsers if needed
- Run all tests
- Generate reports

#### Manual Setup
```sh
# Install browsers for the first run
npx playwright install

# Ensure backend server is running with mock email
cd ../server-api
MOCK_EMAIL=true air

# In another terminal, start dev server
cd ../base-app
yarn dev

# Run the end-to-end tests
yarn test:e2e

# Run specific test suites
yarn test:e2e e2e/auth-flow.spec.ts
yarn test:e2e e2e/password-management.spec.ts

# Run tests in headed mode (see browser)
yarn test:e2e --headed

# Run tests in debug mode
yarn test:e2e --debug

# Run specific test by name
yarn test:e2e -g "forgot password"
```

#### Test Coverage
- ✅ Authentication flows (login, register, logout)
- ✅ Password management (forgot, reset, change)
- ✅ Form validation
- ✅ API integration
- ✅ Error handling

See [E2E_TEST_GUIDE.md](./E2E_TEST_GUIDE.md) for detailed testing documentation.

### Lint with [ESLint](https://eslint.org/)

```sh
yarn lint
```
