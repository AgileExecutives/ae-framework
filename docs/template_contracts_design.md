**Modular Template Contracts & Registration Design**

Ziel
- Beschreibe einen integrierten, redundanzarmen Ansatz, bei dem die Verantwortung für Template‑Contracts bei den jeweiligen Modulen liegt (sowohl in `serverbase/modules` als auch in `shared-modules`).

Kernaussage
- Jedes Modul ist Eigentümer seiner Template‑Contracts. Module implementieren optional eine Methode `RegisterContracts(registrar *templates.ContractRegistrar, tenantID uint) error` am Modul‑Objekt. Der zentrale Bootstrap ruft diese Methode einmal pro Tenant auf. Keine zentrale Datei pflegt Listen von Modulkontrakten manuell — Module bleiben Quelle der Wahrheit.

Design‑Übersicht
- Zentrale Bootstrap‑Funktion: `RegisterAllContracts(db *gorm.DB)` in [serverbase/pkg/startup/contracts.go](serverbase/pkg/startup/contracts.go)
- Module sind verantwortlich für:
  - Optional: Implementieren der Methode `RegisterContracts(registrar, tenantID)` auf dem Modul‑Objekt (Variante 2, empfohlen) oder alternativ Bereitstellung einer freien Funktion `Register<ModuleName>Contracts`.
  - Hosten ihrer `contracts/*.json` Dateien innerhalb ihres Modulpfads.
- Registrar: `ContractRegistrar` in [serverbase/modules/templates/services/service.go](serverbase/modules/templates/services/service.go) bietet `RegisterContractFromFile` / `RegisterContractFromBytes`.

Verantwortlichkeiten (Single Source of Truth)
- Module:
  - Pflegen die `contracts/*.json` Dateien.
  - Implementieren optional `RegisterContracts` auf dem Modul‑Objekt, die `ContractRegistrar` verwendet.
  - Falls ein Modul keine Contracts hat, liefert es eine leere/no‑op Implementierung oder keine Implementierung (wird übersprungen).
- Central Bootstrap (`RegisterAllContracts`):
  - Iteriert Tenants und ruft für jeden Tenant die `RegisterContracts`‑Methode der Module (falls implementiert) auf.
  - Optionaler Fallback: scannt `shared-modules/*/contracts` für Module, die noch keine modulare Methode implementiert haben.

Warum dieses Muster?
- Minimale Redundanz: Contract‑Listen leben neben dem Modulcode und werden nicht doppelt in zentralen Maps gepflegt.
- Klare Ownership: Module halten ihre Contracts nahe am Code, der sie nutzt.
- Testbarkeit: Modulspezifische Registration kann isoliert getestet werden.

Konkrete Änderungen / Empfehlungen
1) Ensure `templates` module does not duplicate module‑owned contracts
  - Entferne oder markiere `ensureStandardTemplates` in [serverbase/modules/templates/module.go](serverbase/modules/templates/module.go) als Test/Harness‑Only.
  - Wenn weiterhin Default‑Templates für Tests benötigt werden, verschiebe deren Erzeugung in `test`‑Fixtures oder in eine `--test-seed` Option.

2) Einheitliche Module‑API (Variante 2)
  - Konvention (empfohlen): Module implementieren optional die Methode
    `RegisterContracts(registrar *templates.ContractRegistrar, tenantID uint) error` auf dem Modul‑Objekt.
  - Alternativ: Module können weiterhin `Register<ModuleName>Contracts(registrar, tenantID) error` als freie Funktion anbieten.
  - Beispiele: `modules/user/services/contract_registration.go`, `shared-modules/booking/services/contract_registration.go`.

3) Bootstrap Anpassungen
  - `RegisterAllContracts` bleibt verantwortlich für Tenant‑Iteration und orchestriert Aufrufe an Modul‑Registrierer.
  - Refactor `RegisterAllContracts` so it iterates module registry and, via a type assertion, invokes `RegisterContracts` for modules that implement it.

4) Fallback‑Scan
  - Behalte für `shared-modules` optional den Dateisystemscan als Fallback, aber logge deutlich, dass modulare `RegisterContracts`‑Methoden bevorzugt sind.

5) Seeder / Templates Erzeugung
  - Seeder, der `templates` Einträge erstellt, darf NICHT automatisch Dateien in `statics/` referenzieren. Stattdessen:
    - Wenn ein Modul liefert eine `contracts/*.json` und optional eine `templates/*.html`, dann soll der Modul‑Registrierer die Contract (via `ContractRegistrar`) aufnehmen und optional die Template‑Inhalte per `StorageService.StoreTemplate` in den document‑Store hochladen, und dann `templates`‑Metadaten erstellen.
    - Seeder im zentralen Bootstrap darf lediglich die Existenz prüfen und ggf. Module‑Registrierer für Default‑Templates anstoßen — die Erstellungslogik bleibt im Modul.

6) Default vs Managed Templates
  - Optional Feld `is_managed` (DB‑Migration später). Wenn `true`, erlaubt das System, Managed Defaults durch ein Reapply zu aktualisieren. Ohne Migration: Seeder nur initial erzeugen, niemals überschreiben.

Implementierung
 - Erweitere das Modul‑Interface um eine optionale Methode `RegisterContracts(registrar *templates.ContractRegistrar, tenantID uint) error` (keine Breaking‑Change, da optional). `RegisterAllContracts` führt pro Modul eine Typ‑Assertion aus und ruft die Methode für Module auf, die sie implementieren.
   - Vorteile: explizite, typisierte API; Module bleiben Quelle der Wahrheit; einfache Testbarkeit.

Priorisierte TODO‑Liste
 - **High:** Ensure module ownership and remove duplication
   - **Task:** Remove contract seeding from `ensureStandardTemplates` (or mark it test‑only). File: [serverbase/modules/templates/module.go](serverbase/modules/templates/module.go).
   - **Task:** Ensure every module that has contract files implements optional `RegisterContracts` on its module object. Files to update: modules without helpers (scan report will list these).
 - **High:** Implement per‑tenant invocation (Variante 2)
   - **Task:** Keep `RegisterAllContracts(db)` as orchestrator but refactor to iterate modules and call `RegisterContracts` via type assertion for modules that implement it. Remove hardcoded lists after migration.
 - **Medium:** Storage handling for template files
   - **Task:** Define and use `StorageService.StoreTemplate` in module registrars when uploading default HTML. File: [serverbase/pkg/services/storage_service.go](serverbase/pkg/services/storage_service.go).
 - **Medium:** Add `is_managed` flag (deferred migration)
   - **Task:** Prepare migration plan and DB migration script (apply when you accept schema changes).
 - **Medium:** Tests
   - **Task:** Add unit tests for `RegisterAllContracts` and per‑module `RegisterContracts` implementations. Use `server-test` harness patterns.
 - **Low:** Developer ergonomics
   - **Task:** Add CLI or admin endpoint to re‑apply module defaults for a tenant (respect `is_managed`).

Developer Workflows / How to add a new module contract
1. Add `contracts/<name>-contract.json` to your module folder.
2. Implement optional `RegisterContracts(contractRegistrar, tenantID) error` on the module object, which calls `contractRegistrar.RegisterContractFromFile(tenantID, "<module>", path)` for each contract.
3. (Optional) If your module provides default template HTML files, implement uploading inside the registrar using `StorageService` and create appropriate `templates` rows.

Minimal Code Snippets (patterns)
1) Module registrar (example `modules/user/services/contract_registration.go`)
 - See existing implementation in [modules/user/services/contract_registration.go](serverbase/modules/user/services/contract_registration.go).

2) ContractRegistrar usage (example)
 - Use `templateServices.NewContractRegistrar(db)` and call `RegisterContractFromFile(tenantID, module, path)`.

Notes & Risks
 - Backwards compatibility: Some tests and the `templates` module seeders currently expect central seeds. Adjust tests or keep a test‑only path.
 - Schema changes (e.g., `is_managed`) require migrations and coordination — defer unless you want to recreate DB.

Nächste Schritte (wenn du zustimmst)
   1. Definiere das optionale Interface `ContractRegistrarCapable` (z. B. in `pkg/core` oder `pkg/startup`).
   2. Passe `RegisterAllContracts` an: während der Modul‑Iteration wird per Typ‑Assertion geprüft und ggf. `RegisterContracts` aufgerufen.
   3. Aktualisiere Module mit Contracts (z. B. `modules/user`, `shared-modules/booking`, `shared-modules/invoice`) und implementiere `RegisterContracts` dort — die Implementierung ruft intern `services.RegisterXContracts` oder `contractRegistrar.RegisterContractFromFile`.
   4. Entferne die zentralen, harten Listen in `RegisterAllContracts` (oder halte sie vorübergehend als Fallback, bis alle Module migriert sind).
   5. Aktualisiere Tests / server‑test Harness, um die neue Methode zu berücksichtigen.

