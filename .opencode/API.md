# API Reference

Base URL: `/api/v1`

All responses are wrapped in `JSONResponse`:

```json
{
  "items": "<T | null>",
  "isSuccess": true,
  "message": "string"
}
```

On error, `items` is `null`, `isSuccess` is `false`, and `message` includes the error detail.

**Note on field naming:** All DB models embed `gorm.Model`, which serializes fields as PascalCase: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`. These fields have NO explicit `json` tags in Go, so they use Go's default JSON serialization (exported field name as-is). All other fields (e.g., `username`, `email`, `domainId`, `parentId`) have explicit lowercase `json` tags.

## Authentication

| Header | Description | Required |
|---|---|---|
| `Authorization: Bearer <token>` | JWT access token | Most endpoints |
| `X-Refresh-Token` | JWT refresh token | Refresh only |

## Multi-Tenancy Headers

Most authenticated endpoints accept these headers for RBAC context:

| Header | Type | Description |
|---|---|---|
| `X-App-Id` | `integer` | App identifier |
| `X-Domain-Id` | `integer` | Domain identifier |

## Pagination

List endpoints (`GET` on collections) support query parameters:

| Param | Type | Description |
|---|---|---|
| `page` | `integer` | Page number |
| `size` | `integer` | Page size |
| `sort` | `string` | Sort expression |
| `filters` | `string` | Filter expression |

---

## Users

### Register

```
POST /v1/users
```

No auth required.

**Body:** `UserInput`

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:** `201` — `{ "items": <newUserId: int>, "isSuccess": true, "message": "User created" }`

### Login

```
POST /v1/users/login
```

No auth required.

**Body:** `LoginInput`

```json
{
  "username": "string",
  "password": "string",
  "isRemember": false
}
```

**Response:** `200` — `{ "items": LoginResponse, "isSuccess": true }`

`LoginResponse`:

```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "user": { "ID": 1, "username": "string", "email": "string", "CreatedAt": "string", "UpdatedAt": "string" }
}
```

### Refresh Token

```
POST /v1/users/refresh
```

**Header:** `X-Refresh-Token: <token>`

**Response:** `200` — `{ "items": LoginResponse }`

### List Users

```
GET /v1/users
```

Auth + `users#read`. Headers: `X-App-Id`, `X-Domain-Id`. Pagination query params.

**Response:** `200` — `{ "items": <paginated User[]>, "isSuccess": true }`

### Get User by ID

```
GET /v1/users/{id}
```

Auth + `users#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": User, "isSuccess": true }`

`User`:

```json
{
  "ID": 1,
  "username": "string",
  "email": "string",
  "CreatedAt": "string",
  "UpdatedAt": "string",
  "DeletedAt": null
}
```

### Update User

```
PUT /v1/users/{id}
```

Auth + `users#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `User`

```json
{
  "username": "string",
  "email": "string"
}
```

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "User updated" }`

### Delete User

```
DELETE /v1/users/{id}
```

Auth + `users#delete`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "User deleted" }`

### Get User Roles

```
GET /v1/users/{id}/roles
```

Auth + `users#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": string[], "isSuccess": true }` (role names)

### Assign Roles to User

```
PATCH /v1/users/{id}/roles
```

Auth + `users#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `uint[]` (role IDs)

```json
[1, 2, 3]
```

**Response:** `200` — `{ "items": string[], "isSuccess": true }` (role names)

### Get User Organizations

```
GET /v1/users/{id}/organizations
```

Auth required. Path param: user ID.

**Response:** `200` — `{ "items": Domain[], "isSuccess": true }` (list of domains/organizations the user belongs to)

### Get User Permissions

```
GET /v1/users/{id}/permissions
```

Auth + `users#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": string[][], "isSuccess": true }` (permission tuples: `[sub, app, dom, res, act]`)

---

## Roles

### List Roles

```
GET /v1/roles
```

Auth + `roles#read`. Headers: `X-App-Id`, `X-Domain-Id`. Pagination query params.

**Response:** `200` — `{ "items": <paginated Role[]>, "isSuccess": true }`

### Create Role

```
POST /v1/roles
```

Auth + `roles#create`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `RoleInput`

```json
{
  "name": "string",
  "description": "string",
  "domainId": 1
}
```

**Response:** `201` — `{ "items": <newRoleId: int>, "isSuccess": true }`

### Get Role by ID

```
GET /v1/roles/{id}
```

Auth + `roles#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": Role, "isSuccess": true }`

`Role`:

```json
{
  "ID": 1,
  "name": "string",
  "description": "string",
  "domainId": 1,
  "domain": { "ID": 1, "name": "string", "description": "string", "isOrganization": true, "parentId": null, "parent": null, "CreatedAt": "string", "UpdatedAt": "string" },
  "CreatedAt": "string",
  "UpdatedAt": "string"
}
```

### Update Role

```
PUT /v1/roles/{id}
```

Auth + `roles#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `Role`

```json
{
  "name": "string",
  "description": "string",
  "domainId": 1
}
```

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Role updated" }`

### Delete Role

```
DELETE /v1/roles/{id}
```

Auth + `roles#delete`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Role deleted" }`

### Get Role Permissions

```
GET /v1/roles/{id}/permissions
```

Auth + `roles#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": string[][], "isSuccess": true }` (permission tuples)

### Assign Permissions to Role

```
PATCH /v1/roles/{id}/permissions
```

Auth + `roles#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `uint[]` (permission IDs)

```json
[1, 2, 3]
```

**Response:** `200` — `{ "items": string[], "isSuccess": true }` (permission names)

---

## Permissions

### List Permissions

```
GET /v1/permissions
```

Auth + `permissions#read`. Headers: `X-App-Id`, `X-Domain-Id`. Pagination query params.

**Response:** `200` — `{ "items": <paginated Permission[]>, "isSuccess": true }`

### Create Permission

```
POST /v1/permissions
```

Auth + `permissions#create`. Header: `X-App-Id`.

**Body:** `PermissionInput`

```json
{
  "name": "string",
  "description": "string",
  "appId": 1
}
```

**Response:** `201` — `{ "items": <newPermissionId: int>, "isSuccess": true }`

### Get Permission by ID

```
GET /v1/permissions/{id}
```

Auth + `permissions#read`. Header: `X-App-Id`.

**Response:** `200` — `{ "items": Permission, "isSuccess": true }`

`Permission`:

```json
{
  "ID": 1,
  "name": "string",
  "description": "string",
  "appId": 1,
  "app": { "ID": 1, "name": "string", "description": "string", "CreatedAt": "string", "UpdatedAt": "string" },
  "CreatedAt": "string",
  "UpdatedAt": "string"
}
```

### Update Permission

```
PUT /v1/permissions/{id}
```

Auth + `permissions#update`. Header: `X-App-Id`.

**Body:** `Permission`

```json
{
  "name": "string",
  "description": "string",
  "appId": 1
}
```

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Permission updated" }`

### Delete Permission

```
DELETE /v1/permissions/{id}
```

Auth + `permissions#delete`. Header: `X-App-Id`.

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Permission deleted" }`

---

## Domains

### List Domains

```
GET /v1/domains
```

Auth + `domains#read`. Headers: `X-App-Id`, `X-Domain-Id`. Pagination query params.

**Response:** `200` — `{ "items": <paginated Domain[]>, "isSuccess": true }`

### Create Domain

```
POST /v1/domains
```

Auth + `domains#create`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `DomainInput`

```json
{
  "name": "string",
  "description": "string",
  "isOrganization": true,
  "parentId": null
}
```

**Response:** `201` — `{ "items": <newDomainId: int>, "isSuccess": true }`

### Get Domain by ID

```
GET /v1/domains/{id}
```

Auth + `domains#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": Domain, "isSuccess": true }`

`Domain`:

```json
{
  "ID": 1,
  "name": "string",
  "description": "string",
  "isOrganization": true,
  "parentId": null,
  "parent": { /* recursive Domain or null */ },
  "CreatedAt": "string",
  "UpdatedAt": "string"
}
```

### Update Domain

```
PUT /v1/domains/{id}
```

Auth + `domains#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `Domain`

```json
{
  "name": "string",
  "description": "string",
  "isOrganization": true,
  "parentId": null
}
```

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Domain updated" }`

### Delete Domain

```
DELETE /v1/domains/{id}
```

Auth + `domains#delete`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "Domain deleted" }`

---

## Apps

### List Apps

```
GET /v1/apps
```

Auth + `apps#read`. Headers: `X-App-Id`, `X-Domain-Id`. Pagination query params.

**Response:** `200` — `{ "items": <paginated App[]>, "isSuccess": true }`

### Create App

```
POST /v1/apps
```

Auth + `apps#create`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `AppInput`

```json
{
  "name": "string",
  "description": "string"
}
```

**Response:** `201` — `{ "items": <newAppId: int>, "isSuccess": true }`

### Get App by ID

```
GET /v1/apps/{id}
```

Auth + `apps#read`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": App, "isSuccess": true }`

`App`:

```json
{
  "ID": 1,
  "name": "string",
  "description": "string",
  "CreatedAt": "string",
  "UpdatedAt": "string"
}
```

### Update App

```
PUT /v1/apps/{id}
```

Auth + `apps#update`. Headers: `X-App-Id`, `X-Domain-Id`.

**Body:** `App`

```json
{
  "name": "string",
  "description": "string"
}
```

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "App updated" }`

### Delete App

```
DELETE /v1/apps/{id}
```

Auth + `apps#delete`. Headers: `X-App-Id`, `X-Domain-Id`.

**Response:** `200` — `{ "items": null, "isSuccess": true, "message": "App deleted" }`

---

## Access Controls

### Evaluate Access

```
POST /v1/access-controls/eval
```

No auth.

**Body:** `AccessControlEval`

```json
{
  "sub": "string",
  "app": "string",
  "dom": "string",
  "obj": "string",
  "act": "string"
}
```

**Response:** `200` — `{ "items": <bool>, "isSuccess": true }`

### Get Permissions for User

```
GET /v1/access-controls/permissions/{app}/{domain}/{user}
```

No auth. Path params: app name, domain name, username.

**Response:** `200` — `{ "items": string[][], "isSuccess": true }` (permission tuples)

### Get All Roles

```
GET /v1/access-controls/roles/{domain}
```

No auth. Path param: domain name.

**Response:** `200` — `{ "items": string[][], "isSuccess": true }` (role tuples)

---

## Token

### Check Token Permission

```
POST /v1/token
```

Auth required.

**Body:** `TokenEndpointRequest`

```json
{
  "appId": 1,
  "domainId": 1,
  "permission": "resource#action"
}
```

**Response:** `200` — `{ "items": <bool>, "isSuccess": true }`

---

## Error Responses

All errors follow the same `JSONResponse` shape:

```json
{
  "items": null,
  "isSuccess": false,
  "message": "error description: detail"
}
```

| Status | Meaning |
|---|---|
| `400` | Bad request / invalid JSON / invalid header |
| `401` | Missing or invalid JWT |
| `403` | Authenticated but lacks permission |
| `404` | Resource not found |
| `500` | Internal server error |
