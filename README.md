# NapiWay - New API Generation Way

NapiWay is a specification-driven API code generator.

It generates:

* Go HTTP server helpers and models
* Go Client SDK
* TypeScript Client SDK

All code is generated from a custom YAML specification file.

# Specification Rules

## 1. File Format

* The specification file **must be written in YAML**.
* At least one generation target must be defined:

  * `goServer`
  * `goSdk`
  * `tsSdk`
* The `spec` section is required.

## 2. Naming Convention (Strict)

All logical names **must be PascalCase**, except transport-level names.

PascalCase format:

```
^[A-Z][a-zA-Z0-9]*$
```

### Must Be PascalCase

* `spec.apiName`
* Endpoint names (keys inside `endpoints`)
* `auth[].id`
* `auth[].name`
* `Param.name`
* All body field names inside:

  * `requestBody.properties`
  * `responses[*].body.properties`
  * nested `object.properties`

### Must NOT Be PascalCase

Transport-level names must match actual HTTP usage and are not required to follow PascalCase:

* `transportName` (headers, query params, path params)
* Path parameter placeholders inside `path`

Example:

```
path: /users/{userId}
```

Here:

* `UserId` → Param.name (PascalCase)
* `userId` → transportName and path placeholder

## 3. Path Parameters

Path parameters must be declared using:

```
{paramName}
```

Example:

```
/users/{userId}
```

Rules:

* The value inside `{}` must match the `transportName` of the corresponding `pathParam`.
* `Param.name` must be PascalCase.
* `transportName` must match the path placeholder exactly.

This is required because:

* Go server parsing relies on matching the placeholder
* TypeScript SDK replaces placeholders using `transportName`

## 4. Authentication

Authentication methods are defined globally inside `spec.auth`.

* Only `type: header` is supported.
* Endpoint authentication references `auth[].id`.
* Referencing an undefined auth ID will fail validation.
* If `auth` is omitted in an endpoint, the endpoint requires no authentication.

Endpoint auth supports:

* `all`
* `any`
* If both are defined, both conditions must be satisfied.

## 5. Field Rules

Supported field types:

* `string`
* `number`
* `boolean`
* `object`
* `array`

### Object Type

If `type: object` is used, `properties` must be defined.

### Array Type

If `type: array` is used, `items` must be defined.

### nonEmpty

* Only valid for `string` and `array`.
* If `nonEmpty: true`, then `required: true` must also be set.
* Using `nonEmpty` on other types fails validation.

## 6. Request and Response Rules

* HTTP methods allowed:

  * `GET`
  * `POST`
  * `PUT`
  * `DELETE`

* HTTP status codes must be between `100` and `599`.

* If `contentType` is omitted, it defaults to:

  ```
  application/json
  ```

* If `requestBody` is defined, `properties` is required.

## 7. Generated Code Policy

All generated code is fully managed by NapiWay.

Manual modification of generated files is discouraged, as regeneration will overwrite changes.

If customization is required, extend the generated code externally rather than modifying generated files directly.

# Generated Files

## Go

### Server Generation

Generates models and helper functions for API endpoints.

Files generated:

* `helpers.go`
  Contains per-endpoint `DecodeAndValidate` request functions.

* `models.go`
  Contains request and per-status response models.

### Go Client SDK Generation

Generates a standalone Go SDK.

Files generated:

* `client.go`
  Client struct and endpoint methods.

* `models.go`
  Request and response models.

* `go.mod`
  Go module definition.

* `.gitignore`

* `README.md`

## TypeScript

### Client SDK Generation

Generates a standalone TypeScript SDK.

Files generated:

* `client.ts`
  Client class and endpoint methods.

* `models.ts`
  Request and response models.

* `package.json`

* `tsconfig.json`

* `.gitignore`

* `README.md`
