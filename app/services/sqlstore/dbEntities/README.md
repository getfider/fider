# dbEntities Package

This package contains database entity mapping structs that convert PostgreSQL query results to domain models.

## Purpose

The `dbEntities` package provides a centralized location for database-to-model conversion logic, separating data mapping concerns from the business logic in the `postgres` package.

## Exported Types

### User (and UserProvider)

The `User` type is **exported** so it can be used across packages:

```go
import "github.com/getfider/fider/app/services/sqlstore/dbEntities"

var user dbEntities.User
err := trx.Get(&user, "SELECT id, name, email... FROM users WHERE id = $1", userID)
if err != nil {
    return err
}

// Convert to entity.User
entityUser := user.ToModel(ctx)
```

## Unexported Types

All other types (comment, post, tag, tenant, etc.) are **unexported** (lowercase) and only used internally by the `postgres` package. They can be exported later if needed by simply capitalizing the type name and its `toModel` method.

## Testing

The package includes unit tests for the exported User type. Run them with:

```bash
godotenv -f .test.env go test ./app/services/sqlstore/dbEntities/...
```
