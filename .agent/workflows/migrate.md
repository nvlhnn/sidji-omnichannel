---
description: How to create and run database migrations
---

## Creating Migrations

### 1. Check the last migration number
// turbo
```powershell
Get-ChildItem migrations\*.up.sql | Sort-Object Name | Select-Object -Last 1
```

### 2. Create migration files
Create both up and down files with the next sequential number:
- `migrations/NNN_description.up.sql` — the forward migration
- `migrations/NNN_description.down.sql` — the rollback migration

### 3. Restart the server
Migrations are auto-applied on server startup. Just restart:

Local:
```powershell
# The Go server will auto-run pending migrations on boot
go run cmd/server/main.go
```

Production (VPS):
```bash
docker compose -f docker-compose.prod.yml restart api
```

### Notes
- Migrations run automatically — no manual `psql` needed!
- The `schema_migrations` table tracks which migrations have been applied
- Each migration runs in a transaction (all-or-nothing)
- Always create both `.up.sql` and `.down.sql` files
- Number migrations sequentially (check existing ones first!)
