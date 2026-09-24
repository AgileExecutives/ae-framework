Module lifecycle — AE Framework (serverbase)

This document explains how modules are registered, initialized, migrated, seeded, and started within the AE Framework (`serverbase`). It also describes per-tenant contract registration, event hooks, and recommended best practices.

Overview

- The application boot sequence (see `bootstrap.Application`) drives module lifecycle: registration → Initialize() → migrations → Start event bus/modules → seed an empty database → serve HTTP.
- Modules implement the `core.Module` interface (or use the adapter helpers) and register with the app via `app.RegisterModule(module)`.

Core concepts

- `core.Module` — the module abstraction. Key methods typically implemented:
  - `Name() string`, `Version() string`, `Dependencies() []string`
  - `Initialize(ctx core.ModuleContext) error` — called during app initialization (synchronous)
  - `Start(ctx context.Context) error` — start background work
  - `Stop(ctx context.Context) error` — stop background work
  - `Entities() []core.Entity` — used for AutoMigrate
  - `Routes() []core.RouteProvider` — register HTTP routes

- `core.ModuleContext` — the runtime context passed into modules. Contains:
  - `DB` — *gorm.DB* database handle
  - `Router` — *gin.Engine* or router group
  - `EventBus` — central event bus for pub/sub
  - `Config`, `Logger`, `Services` (service registry), `DocRegistry`, `TokenService`, etc.

- `core.ModuleRegistry` — manages initialization order via module dependencies, and invokes lifecycle methods.

Adapter helpers

- Many modules use `module.NewAdapterModule(...)` and helper options:
  - `module.WithEntities(...)` — declare entities for migrations
  - `module.WithInit(func(ctx core.ModuleContext) error)` — inserts initialization behavior
  - `module.WithServices(...)` — register services into the `Services` registry
  - `module.WithRoutes(...)` — provide `RouteProvider`s
  - `module.WithContractRegistration(func(ctx core.ModuleContext, tenantID uint) error)` — register per-tenant template contracts and seed templates for an individual tenant

Boot sequence (high level)

1. Application constructs `bootstrap.Application` and registers modules.
2. `Application.Initialize()`:
   - initialize core services and create `ModuleContext`
   - call `registry.InitializeAll(app.context)` — runs each module's Initialize (dependency order)
   - merge module swagger docs (modules can register docs during Initialize)
   - run migrations in dependency order: collects `Entities()` from modules and calls `DB.AutoMigrate(...)`
3. `Application.Start()`:
   - start EventBus
   - call `registry.StartAll(ctx)` — start modules (background services)
  - if no tenants exist, seed the initial tenant, organization, and admin
  - the tenant service publishes `tenant.created`; registered modules handle this event and invoke `startup.RegisterContractsForTenant` to create tenant-scoped contracts and templates
   - start HTTP server

Migrations and seeding

- Modules should expose their GORM models via `Entities()` so migrations run centrally in the correct order.
- Do not perform tenant-scoped seeding inside `Initialize()` unless you can guarantee tenants already exist. Prefer to:
  - keep `Initialize()` idempotent and limited to wiring services, AutoMigrate for module entities, and registering doc JSON
  - perform tenant-scoped seeding via `WithContractRegistration` or by reacting to the `tenant.created` event
- The top-level `ensureInitialTenantAndOrganization()` in `bootstrap/application.go` creates the initial tenant through `TenantService`. Because the event bus and modules are already started, the normal `tenant.created` event flow registers its templates/contracts before startup seeding continues.

Contract registration & per-tenant templates

- Two mechanisms:
  1. Module-owned registration: implement a module-level `RegisterContracts(core.ModuleContext, tenantID uint) error` (exposed by adapter `WithContractRegistration`) that inserts `templates` and `template_contracts` rows scoped to `tenantID`.
  2. Filesystem scanner: `startup.RegisterAllContracts` scans `modules/*/contracts` and `shared-modules/*/contracts` and uses a `ContractRegistrar` service to insert contract rows per tenant.
- Central orchestrator (`startup.RegisterContractsForTenant`) enumerates modules and invokes module-provided callbacks and falls back to filesystem scanning.
- Important: always set `tenant_id` on inserted rows. Avoid legacy seeding paths that insert tenantless (0) rows.

Event-driven tenant lifecycle

- `TenantService` publishes `tenant.created` after the tenant row is created; modules may subscribe to this event to perform tenant-specific initialization (create MinIO buckets, default organization settings, templates, template_contracts, etc.).
- Advantages: idempotent, runs after tenant record exists, and supports external tenant creation flows.

Service registry & wiring

- Modules can register services into the central `Services` registry using provider factories (example: `module.WithServices`). Other modules fetch services by name with `ctx.Services.Get("service_name")`.
- Prefer typed service interfaces and defensive type assertions. Log and avoid panics on missing services.

Best practices

- Keep `Initialize()` safe and idempotent: migrations, wiring, and registering docs only.
- Perform tenant-scoped seeding only after tenants exist: use `WithContractRegistration` callbacks, `startup.RegisterContractsForTenant`, or react to `tenant.created` events.
- Implement module-level `RegisterContracts(ctx core.ModuleContext, tenantID uint) error` when the module owns templates or contracts. Centralize per-tenant registration to avoid duplication.
- Avoid `panic()` inside startup paths — return errors or log and continue where appropriate; make registries tolerant and retryable.
- Write tests for per-tenant registration (use in-memory DB) and verify tenant_id is set on inserted rows.

Testing tips

- Unit-test `RegisterContractsForTenant` and module `RegisterContracts` implementations by creating a temporary `gorm` in-memory DB and invoking the callbacks.
- Ensure migrations run in tests before seeding. Clear relevant tables between test runs to avoid UNIQUE constraint conflicts.

Quick checklist when adding a new module

- [ ] Implement `core.Module` (or adapter) and declare `Entities()`
- [ ] Use `module.WithServices` to provide services and register service names
- [ ] Use `module.WithRoutes` to expose HTTP handlers
- [ ] If module owns templates/contracts implement `WithContractRegistration` for per-tenant registration
- [ ] Ensure `Initialize()` does not insert tenant-scoped rows unconditionally
- [ ] Add contract JSON to `modules/<name>/contracts` if using filesystem fallback

References

- bootstrap: `serverbase/pkg/bootstrap/application.go`
- startup contract orchestration: `serverbase/pkg/startup`
- templates module example: `serverbase/modules/templates/module.go`


