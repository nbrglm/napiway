---
# API Generator Specification Guide

This document explains how to write a valid YAML specification file for the generator.

The specification defines:

* API metadata
* Authentication methods
* Endpoints
* Request/response models
* Code generation targets (Go server / Go SDK / TS SDK)

---

# 1. File Structure Overview

A valid spec file has this top-level structure:

```yaml
goServer:   # optional
goSdk:      # optional
tsSdk:      # optional

spec:       # REQUIRED
```

At least one of `goServer`, `goSdk`, or `tsSdk` must be defined.

---

# 2. Code Generation Configuration

These sections define what should be generated.

## 2.1 Go Server

```yaml
goServer:
  outputDir: ./path/to/output      # REQUIRED
  packageName: packagename         # REQUIRED
```

* `outputDir`: Directory where generated files will be written.
* `packageName`: Go package name used in generated files.

---

## 2.2 Go SDK

```yaml
goSdk:
  outputDir: ./path/to/output      # REQUIRED
  moduleName: github.com/org/sdk   # REQUIRED
```

* `outputDir`: Output directory.
* `moduleName`: Go module path for go.mod.

---

## 2.3 TypeScript SDK

```yaml
tsSdk:
  outputDir: ./path/to/output      # REQUIRED
  packageName: "@org/sdk"          # REQUIRED

  description: "SDK description"   # optional
  author: "Author Name"            # optional
  license: "MIT"                   # optional
  repository: "https://..."        # optional
  website: "https://..."           # optional
  keywords:
    - api
    - sdk
```

Only `outputDir` and `packageName` are required.

---

# 3. The `spec` Section (Required)

This defines the actual API.

```yaml
spec:
  apiName: "My API"          # REQUIRED
  version: "1.0.0"           # REQUIRED
  description: "API desc"    # REQUIRED

  website: "https://..."     # optional
  contact: "support@..."     # optional
  license: "MIT"             # optional

  auth: []                   # optional
  endpoints: {}              # REQUIRED
```

---

# 4. Authentication

Authentication methods are defined globally inside `spec.auth`.

```yaml
auth:
  - id: apiKey                  # REQUIRED (unique identifier)
    name: API Key               # REQUIRED (display name)
    transportName: X-API-Key    # REQUIRED (actual HTTP header name)
    type: header                # REQUIRED (currently only 'header')
    description: "Optional"
    format: "Bearer {token}"    # Optional documentation field
```

### Important Rules

* `id` is referenced by endpoints.
* Only `type: header` is supported.
* If an endpoint references an unknown `id`, validation fails.

---

# 5. Endpoints

Endpoints are defined as a map:

```yaml
endpoints:
  GetUser:
    ...
  CreateUser:
    ...
```

The key (e.g., `GetUser`) is a unique identifier used for code generation.

---

## 5.1 Basic Endpoint Structure

```yaml
method: GET | POST | PUT | DELETE   # REQUIRED
path: /users/{id}                   # REQUIRED

contentType: application/json       # optional
description: "..."                  # optional

pathParams: []                      # optional
headers: []                         # optional
queryParams: []                     # optional

auth: {}                            # optional
requestBody: {}                     # optional
responses: {}                       # optional
```

---

# 6. Parameters

Used for path params, headers, and query params.

```yaml
- name: id                  # REQUIRED (internal name)
  transportName: id         # REQUIRED (HTTP name)
  type: string              # REQUIRED (string | number | boolean)
  description: "Optional"
  required: true
```

### Required Behavior

* For `string`, `required: true` means non-empty.
* For `number` and `boolean`, `required: true` means must be present.

---

# 7. Endpoint Authentication Rules

```yaml
auth:
  all:
    - apiKey
  any:
    - sessionToken
```

* `all`: every listed method must be satisfied.
* `any`: at least one must be satisfied.
* If both exist, both conditions apply.
* If omitted, endpoint requires no authentication.

---

# 8. Request Body

```yaml
requestBody:
  description: "Optional"
  properties:
    fieldName:
      type: string
```

`properties` is required if `requestBody` is defined.

---

# 9. Responses

Responses are defined per HTTP status code:

```yaml
responses:
  200:
    description: "Success"
    contentType: application/json     # optional
    headers: []                       # optional
    body: {}                          # optional
```

Status code must be between 100 and 599.

---

# 10. Fields (Body Properties)

Fields define request and response models.

```yaml
fieldName:
  type: string | number | boolean | object | array
  description: "Optional"
  required: true
  nonEmpty: true
```

---

## 10.1 Field Rules

### Primitive Types

For:

* `string`
* `number`
* `boolean`

You cannot define `properties` or `items`.

---

### Object Type

If:

```yaml
type: object
```

You MUST define:

```yaml
properties:
  subField:
    type: string
```

---

### Array Type

If:

```yaml
type: array
```

You MUST define:

```yaml
items:
  type: string
```

---

### nonEmpty Rule

* Only valid for:

  * `string`
  * `array`
* If `nonEmpty: true`, then `required: true` must also be set.
* Using `nonEmpty` on `number`, `boolean`, or `object` will fail validation.

---

# 11. Default Behaviors

* If `contentType` is not specified → defaults to `application/json`.
* If response `contentType` is not specified → defaults to `application/json`.
* If `MaxBodyBytes` is not specified → default limit is applied (256 KB).

---

# 12. Common Validation Failures

Your spec will fail if:

* `apiName` or `version` is missing.
* No generation target is defined.
* `object` type has no `properties`.
* `array` type has no `items`.
* `nonEmpty` used without `required`.
* Endpoint references unknown auth ID.
* Invalid HTTP method.
* Invalid HTTP status code.

---

# 13. Naming Conventions (Strict)

All logical names in the specification **must use PascalCase**, except transport-level names.

## Must Be PascalCase

The following fields must follow PascalCase:

* `spec.apiName`
* Endpoint names (keys inside `endpoints`)
* `auth[].id`
* `auth[].name`
* `Param.name`
* All field names inside:

  * `requestBody.properties`
  * `response.body.properties`
  * nested `object.properties`

### Valid Examples

```
CreateUser
UserProfile
ApiKey
SessionToken
UserId
EmailAddress
IsActive
CreatedAt
```

---

## Must NOT Be PascalCase

Transport-level names must reflect actual HTTP usage and may follow standard HTTP conventions:

* `transportName`

  * Header names (e.g. `X-API-Key`)
  * Query parameter names (e.g. `user_id`)
  * Path parameter names (e.g. `id`)

These must match the real HTTP contract and are not required to be PascalCase.

---

## Why This Rule Exists

PascalCase ensures:

* Clean Go struct generation
* Predictable TypeScript model naming
* No case-mismatch issues across SDKs
* Deterministic code generation

Breaking this rule will lead to:

* Inconsistent struct names
* Ugly auto-corrections in generators
* Possible conflicts in strongly typed languages

---

For more info, see the [spec.go](./spec/spec.go)