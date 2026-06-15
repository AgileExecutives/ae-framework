# shared-modules

Produktunabhängige Go-Module für alle AE-Produkte.

## Module

| Modul | Pfad |
|---|---|
| audit | `github.com/ae/shared-modules/audit` |
| booking | `github.com/ae/shared-modules/booking` |
| calendar | `github.com/ae/shared-modules/calendar` |
| documents | `github.com/ae/shared-modules/documents` |
| invoice | `github.com/ae/shared-modules/invoice` |
| invoice_number | `github.com/ae/shared-modules/invoice_number` |

## Abhängigkeit einbinden

```go
// go.mod eines Produkts
require (
    github.com/ae/shared-modules/booking v0.1.0
    github.com/ae/shared-modules/invoice v0.1.0
)
```

## Lokale Entwicklung (go.work)

```go
// go.work im Produkt-Repo (in .gitignore)
use (
    .
    ../shared-modules/booking
    ../shared-modules/invoice
)
```

## Versionierung

Tags folgen dem Schema `<modul>/vX.Y.Z`, z.B.:

```bash
git tag booking/v0.2.0
git push origin booking/v0.2.0
```
