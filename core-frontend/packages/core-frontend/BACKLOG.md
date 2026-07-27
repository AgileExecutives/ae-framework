# Core Frontend Backlog

This backlog contains user stories, todos, and testcases for the `core-frontend` module. Each user story includes acceptance criteria, implementation tasks, and concrete test cases (unit / integration / manual) to validate the work.

**Note:** prioritize API / auth fixes and the personalized styling feature for maintainers who embed `core-frontend` into host applications.

**User Story 1: Prevent browser CORS by using dev proxy for API calls**
- Description: Ensure front-end uses relative API base (`/api/v1`) in development so Vite proxy avoids CORS errors.
- Acceptance Criteria:
  - Dev builds use `/api/v1` by default when `VITE_API_BASE_URL` is not set to an absolute origin.
  - No browser CORS errors when calling backend endpoints during `npm run dev` with proxy configured.
- Todos:
  - [x] Audit repo for absolute origins (`http://localhost:8080` or similar) and replace with relative `/api/v1` where appropriate.
  - [x] Ensure `packages/core-frontend/src/config/api-config.ts` and `packages/unburdy-app/src/main.ts` default to `/api/v1` in dev.
  - [x] Update examples and docs to show usage with dev proxy.
- Testcases:
  - Unit: config value selection when `VITE_API_BASE_URL` set/unset.
  - Integration: Run `npm run dev` and confirm `/api/v1/password-security` requests pass through proxy (manual/automated E2E).

**User Story 2: Center Dashboard and float UI toggles**
- Description: Make dashboard content centered like the registration page and add floating locale/theme toggles top-right.
- Acceptance Criteria:
  - Dashboard layout matches `AuthLayout` centering patterns.
  - Locale and Theme toggles are positioned top-right and accessible on all auth pages.
- Todos:
  - [x] Update `DashboardView.vue` to use shared centering wrapper.
  - [x] Add `LocaleSwitcher` and `ThemeToggle` to top-right slot (or global header component).
  - [x] Add responsive behavior and keyboard access.
- Testcases:
  - Visual: Manual QA on desktop and mobile widths.
  - Unit: Snapshot tests for `DashboardView` component.

**User Story 3: Fill missing i18n keys**
- Description: Populate translation keys `success.registration.title`, `success.registration.message`, `success.goToLogin`, `forgot.successTitle` for all supported locales.
- Acceptance Criteria:
  - Keys are present in `src/locales/en` and `src/locales/de`.
  - No warnings in production build related to missing message keys.
- Todos:
  - [x] Add missing keys to `src/locales/en.ts` and `src/locales/de.ts`.
- Testcases:
  - [x] Add a lint/check step to fail on missing keys during CI (optional). Implemented via `scripts/check-i18n.cjs` and npm scripts `lint:i18n` / `lint:i18n:fix`.

**User Story 4: Fix login redirect after authentication**
- Description: Ensure router middleware and auth store coordinate so users are redirected to dashboard after login.
- Acceptance Criteria:
  - After successful `login()` (token received), user is redirected to intended route (or `/dashboard`).
  - Auth initialization does not clear a freshly set token due to race conditions.
- Todos:
  - [x] Harden `login()` to accept various token shapes and set `initialized = true` after success.
  - [x] Adjust router middleware to avoid revalidation that clears auth immediately after login.
  - [x] Add logs to trace auth initialization in dev.
- Testcases:
  - Unit: Simulate `login()` with different API shapes (token, access_token, nested data) and assert `auth_token` stored.
  - Integration/E2E: Full login flow and redirect to `/dashboard`.

**User Story 5: Define registration flows for SAAS vs Single-tenant**
- Description: Capture requirements and implement adaptive registration UI and front-end contract with back-end for both SAAS (tenant creation + admin user) and single-tenant flows.
- Acceptance Criteria:
  - A single registration route supports both modes depending on tenant context or registration token.
  - For SAAS, steps to create tenant and initial user are clearly defined and use the API client endpoints.
- Todos:
  - [ ] Draft flow docs and wireframe for `/:tenant?/register` and token based registration.
  - [x] Draft flow docs and wireframe for `/:tenant?/register` and token based registration. See `REGISTRATION_FLOW.md`.
  - [ ] Map API client endpoints: `POST /organizations`, `POST /clients/registration/:token`, etc.
  - [ ] Implement front-end forms that adapt (tenant slug, registration token).
  - [ ] Provide backend contract docs for missing endpoints.
- Testcases:
  - Unit: Validation unit tests for registration forms.
  - Integration: Simulate SAAS registration creating org + user (mocked API).

**User Story 6: Audit `api-client` for tenant/registration endpoints**
- Description: Verify the generated API client contains endpoints needed for tenant creation, registration tokens, and organization settings.
- Acceptance Criteria:
  - Confirmed availability of `getRegistrationSettings`, `createOrganization`, `getRegistrationSettingsWithToken`, `registerClient`, and other relevant endpoints.
  - Document any missing or mismatched endpoints and propose fixes.
- Todos:
  - [ ] Produce a mapping doc between front-end needs and `api-client` methods.
  - [ ] Add helper wrappers if the client uses different naming conventions (e.g., `organizations` vs `tenants`).
- Testcases:
  - Unit: Tests that call client methods with mocked `request()` to ensure correct path/verb.

**User Story 7: Personalized styling per-host application (feature request)**
- Description: Allow each application that consumes `core-frontend` as a module to provide per-app styling (colors, fonts, logo) without modifying the core package.
- Acceptance Criteria:
  - Host application can provide either:
    - a global CSS file at `/app-styles.css`, or
    - an inline CSS string via `window.__CORE_FRONTEND_INLINE_STYLES__`, or
    - a URL via `window.__CORE_FRONTEND_STYLE_URL__`.
  - Core automatically loads and applies provided styles at startup before mount.
  - Theming uses CSS variables to make overrides minimal (`--brand-primary`, `--brand-bg`, `--brand-font`).
- Todos:
  - [x] Add runtime style loader plugin in `src/plugins/appStyle.ts`.
  - [ ] Expose documentation and examples for host apps to provide `app-styles.css` or set `window.__CORE_FRONTEND_INLINE_STYLES__`.
  - [ ] Create a small set of CSS variables in `src/style.css` for overriding.
  - [ ] Add tests for plugin (DOM insertion) and e2e test to verify overrides applied.
- Testcases:
  - Unit: Test that `applyAppStyles()` injects a `<link>` when `__CORE_FRONTEND_STYLE_URL__` is set and injects `<style>` when `__CORE_FRONTEND_INLINE_STYLES__` is set.
  - Integration: Host app provides `/app-styles.css` with `--brand-primary: #112233` and assert a DOM element uses the overridden color.
  - Manual: Drop `app-styles.css` into host public folder and verify branding updates on load.

**Maintenance / Chore Stories**
- Replace remaining absolute backend URLs in docs, examples, and tests.
- Add CI checks for i18n key presence and for linting absolute API origins.
  - [x] Add i18n checker script `scripts/check-i18n.cjs` that compares locales to `en`.
  - [x] Add npm scripts `lint:i18n` and `lint:i18n:fix` (`--fix` auto-inserts missing keys using English placeholders).

---

If you want, I can:
- implement the runtime style loader and the minimal CSS variable set now (patch `src/plugins/appStyle.ts`, `src/style.css`, and wire into `src/main.ts`), and
- open a PR or create follow-up tasks for the registration UI.
