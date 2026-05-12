# OAuth Role-Based Access Control

This feature allows you to restrict OAuth login to users who have specific roles provided by your OAuth provider.

Role restrictions are configured **per custom OAuth provider** in the Admin UI.
This means different providers can have completely independent role requirements.

## Configuration

### Admin Panel

1. Go to **Site Settings** → **Authentication** → **OAuth Providers**
2. Edit or create a custom OAuth provider
3. Configure two fields together:
   - **JSON User Roles Path** — the path in the provider's profile JSON where roles are found
   - **Allowed Roles** — a comma-separated list of roles that are permitted to sign in

Both fields must be set for the role check to run. If either is empty, all users are allowed through for that provider.

### Allowed Roles field format

- Comma-separated list of role names, e.g. `ROLE_ADMIN,ROLE_TEACHER`
- Case-sensitive: role names are matched exactly as they appear in the OAuth response
- Leave empty to allow all users (default)

### Examples

#### Allow only admins:
```
ROLE_ADMIN
```

#### Allow admins and teachers:
```
ROLE_ADMIN,ROLE_TEACHER
```

#### Allow all users (default):
Leave the **Allowed Roles** field empty.

## JSON User Roles Path

The roles path tells Fider where to find roles in the OAuth profile response.

### Examples

#### Single role as string:
```json
{
  "id": "12345",
  "role": "ROLE_ADMIN"
}
```
**JSON User Roles Path**: `role`

#### Comma-separated string:
```json
{
  "id": "12345",
  "roles": "ROLE_ADMIN,ROLE_TEACHER"
}
```
**JSON User Roles Path**: `roles`

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
