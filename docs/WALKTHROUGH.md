# Project Walkthrough

## Dynamic Error Code

We added support for dynamic error codes in `AppError`.
- Use `.WithCode("CODE")` on `AppError`.
- Field `error_code` appears in JSON response.

## Database Migrations

We added commands to the `Makefile` to easily run migrations using `golang-migrate`.
These commands automatically load your `.env` file to configure the database connection.

### Commands

| Command | Description |
| :--- | :--- |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback the last migration group |
| `make migrate-force version=N` | Force the database version to N (use if migration state is dirty) |

### Prerequisites
- `make` tool installed (e.g., via MinGW, Chocolatey, or WSL).
- `migrate` tool installed (`golang-migrate`).
- `.env` file with `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`.

### Troubleshooting
If `make` is not recognized on Windows:
1. Install `make` via Chocolatey: `choco install make`
2. Or use the full command manually:
   ```powershell
   migrate -database "mysql://user:pass@tcp(host:port)/dbname" -path migrations up
   ```
