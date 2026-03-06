# OAuth Role-Based Access Control

This feature allows you to restrict OAuth login to users who have specific roles provided by your OAuth provider.

## Configuration

### Environment Variable

Add the following environment variable to your `.env` file:

```bash
OAUTH_ALLOWED_ROLES=ROLE_ADMIN,ROLE_TEACHER
```

- **Format**: Comma-separated list of role names
- **Case-sensitive**: Role names are matched exactly as they appear in the OAuth response
- **Optional**: If not set or empty, all users are allowed to login (default behavior)
- **Per-provider enforcement**: The role check is only applied to OAuth providers that have a **JSON User Roles Path** configured. Providers without a roles path (e.g. built-in Google or GitHub OAuth) are always allowed through, even when `OAUTH_ALLOWED_ROLES` is set. This prevents accidentally blocking users from providers that don't return roles.

### Examples

#### Allow only admins:
```bash
OAUTH_ALLOWED_ROLES=ROLE_ADMIN
```

#### Allow admins and teachers:
```bash
OAUTH_ALLOWED_ROLES=ROLE_ADMIN,ROLE_TEACHER
```

#### Allow all users (default):
```bash
# Don't set OAUTH_ALLOWED_ROLES or leave it empty
OAUTH_ALLOWED_ROLES=
```

## OAuth Provider Configuration

For custom OAuth providers, you need to configure the JSON path to extract roles from the OAuth response.

### Admin Panel Configuration

1. Go to **Site Settings** → **Authentication** → **OAuth Providers**
2. Edit or create a custom OAuth provider
3. Configure the **JSON User Roles Path** field

### JSON User Roles Path Examples

The roles path supports various OAuth response formats:

#### Array of strings:
```json
{
  "id": "12345",
  "name": "John Doe",
  "email": "john@example.com",
  "roles": ["ROLE_ADMIN", "ROLE_TEACHER"]
}
```
**JSON User Roles Path**: `roles`

#### Array of objects (extract specific field):
```json
{
  "id": "12345",
  "roles": [
    { "id": "ROLE_ADMIN", "displayName": "Administrator" },
    { "id": "ROLE_TEACHER", "displayName": "Teacher" }
  ]
}
```
**JSON User Roles Path**: `roles[].id`

This extracts the `id` field from each object in the `roles` array.

#### Nested array:
```json
{
  "id": "12345",
  "user": {
    "profile": {
      "roles": ["ROLE_ADMIN"]
    }
  }
}
```
**JSON User Roles Path**: `user.profile.roles`

#### Nested array of objects:
```json
{
  "id": "12345",
  "user": {
    "groups": [
      { "name": "ROLE_ADMIN" },
      { "name": "ROLE_USER" }
    ]
  }
}
```
**JSON User Roles Path**: `user.groups[].name`

#### Comma-separated string:
```json
{
  "id": "12345",
  "roles": "ROLE_ADMIN,ROLE_TEACHER"
}
```
**JSON User Roles Path**: `roles`

#### Single role as string:
```json
{
  "id": "12345",
  "role": "ROLE_ADMIN"
}
```
**JSON User Roles Path**: `role`
