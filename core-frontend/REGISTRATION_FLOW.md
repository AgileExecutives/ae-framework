# Registration Flows — SAAS vs Single-tenant

This document captures the registration UX flows, API contract mapping, and implementation plan for handling both SAAS (tenant creation + admin user) and single-tenant registration in `core-frontend`.

## Goals
- Support a single registration route that adapts to context (tenant slug, registration token, or host configuration).
- For SAAS: allow creating an organization (tenant) and initial admin user in one flow.
- For single-tenant: simple user creation for an existing tenant.
- Clear mapping to `@ae/api-client` endpoints and fallback behaviors.

## Flow Variants

1) SAAS - Open registration
- Entry: `/register` or `/:tenant?/register` when no registration token required.
- Steps:
  - Step A: Collect organization info (company name, domain/slug) and admin user info.
  - Step B: Call `api.createOrganization` (or equivalent) to create tenant.
  - Step C: Call `api.registerClient` / `api.createUser` to create admin user tied to tenant.
  - Step D: Show verification / success page with `success.registration` messages.

2) SAAS - Token-based (invite or paid flow)
- Entry: `/register?token=abc` or `/register/:token`
- Steps:
  - Validate token: `api.getRegistrationSettingsWithToken(token)` or `api.getRegistrationSettings(token)`
  - Pre-fill tenant info if provided by token
  - Create user via `api.registerClient` or `api.clients.registerWithToken`

3) Single-tenant
- Entry: `/:tenant/register` or embedded app registration
- Steps:
  - Collect user info (name, email, password)
  - Call `api.createUser` (or `clients.register`) scoped to known tenant
  - Redirect to login or dashboard depending on auto-login behavior

## API Client Mapping (audit)
The generated API client should expose endpoints for these operations. Suggested method names (adapt to generated names):

- getRegistrationSettings(token?) -> `getRegistrationSettings` or `registration.getSettings`
- createOrganization(payload) -> `createOrganization` or `organizations.create`
- registerClient(payload) -> `registerClient` or `clients.register`
- registerWithToken(token, payload) -> `clients.registerWithToken` or `registration.registerWithToken`
- createUser(payload) -> `users.create` or `clients.createUser`

If any of these are missing in `@ae/api-client`, add them to the audit TODO and implement a local wrapper in `src/lib/api-registration.ts` that translates to the available endpoints.

## UI/Router Plan
- Add route: `/register` and `/register/:token?` (token optional)
- View component: `src/views/RegisterView.vue` (responsible for flow orchestration)
- Child components:
  - `RegisterOrgForm.vue` — collect tenant/org info
  - `RegisterUserForm.vue` — collect user info
  - `RegistrationSuccess.vue`

## UX Decisions
- If token present and validated, skip org creation step and only collect user details.
- After successful registration: either auto-login and redirect to `/dashboard`, or show success and link to login depending on API behavior.

## Acceptance Criteria
- Single route supports SAAS & single-tenant flows.
- Works with token and without token.
- Integration tests mock API client methods to assert correct sequence of calls.

## Implementation TODO
- [ ] Create `RegisterView.vue` with flow orchestration.
- [ ] Create `RegisterOrgForm.vue`, `RegisterUserForm.vue`, `RegistrationSuccess.vue`.
- [ ] Add route(s) in router index and shared routes.
- [ ] Audit `@ae/api-client` for missing endpoints and add wrapper `src/lib/api-registration.ts` if needed.
- [ ] Add unit tests for form validation and integration tests for end-to-end flow mocking API client.

## Notes
- Prefer defensive parsing of API responses (some endpoints return nested shapes). Use existing patterns from `stores/auth.ts` for tolerant parsing.
- Consider UX for errors returned from backend (map API messages to `messages` i18n keys).
