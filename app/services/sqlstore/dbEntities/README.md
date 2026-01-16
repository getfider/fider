# dbEntities Package

This package contains database entity mapping structs that convert PostgreSQL query results to domain models.

## Purpose

The `dbEntities` package provides a centralized location for database-to-model conversion logic, separating data mapping concerns from the business logic in the `postgres` package.

## Exported Types

### User (and UserProvider)

The `User` type is **exported** for use in commercial/premium features:

```go
import "github.com/getfider/fider/app/services/sqlstore/dbEntities"

// Use dbEntities.User in commercial code
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

## Architecture

**Current State:**
- `postgres` package: Contains all the db* struct definitions and uses them directly
- `dbEntities` package: Contains the same struct definitions with User exported
- Both coexist without conflicts

**Usage Pattern:**
- Regular Fider code continues using `postgres` package as before
- Commercial/premium code can import and use `dbEntities.User` directly
- No changes needed to existing `postgres` code

## Converting Additional Types

To export another type for commercial use:

1. Open the corresponding file (e.g., `comment.go`, `post.go`)
2. Capitalize the struct name: `type comment` → `type Comment`
3. Capitalize the method: `func (c *comment) toModel` → `func (c *Comment) ToModel`
4. Update any internal references if needed

## Testing

The package includes unit tests for the exported User type. Run them with:

```bash
godotenv -f .test.env go test ./app/services/sqlstore/dbEntities/...
```