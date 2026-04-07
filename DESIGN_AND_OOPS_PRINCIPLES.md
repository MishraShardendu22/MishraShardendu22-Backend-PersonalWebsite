# Basic Design and OOP Principles Used in This Codebase

This document explains the core design principles already visible in this backend codebase and how they map to real files.

## 1) High-Level Architecture

The project follows a layered structure:

- Entry and setup: [main.go](main.go)
- Routing layer: [route](route)
- Request handling layer: [controller](controller)
- Cross-cutting concerns: [middleware](middleware), [util](util)
- Data models: [models](models)
- Database setup: [database/connect.database.go](database/connect.database.go)

This organization is a strong example of **Separation of Concerns**.

## 2) SOLID Principles (How They Appear Here)

## S - Single Responsibility Principle (SRP)

Each package has a focused purpose:

- [route](route): URL to handler mapping only.
- [controller](controller): request parsing, validation, business flow.
- [middleware/jwt.middleware.go](middleware/jwt.middleware.go): auth token validation.
- [util/apiResponse.util.go](util/apiResponse.util.go): uniform response shape.
- [database/connect.database.go](database/connect.database.go): DB connection setup only.

Result: easier maintenance and safer changes.

## O - Open/Closed Principle (OCP)

The code is open for extension without heavily modifying existing modules:

- New feature modules are added through new route + controller files.
- Route registration in [main.go](main.go) composes features by calling setup functions such as:
  - `SetupProjectRoutes`
  - `SetupExpRoutes`
  - `SetupCertificationRoutes`

This supports growth with limited impact on existing endpoints.

## L - Liskov Substitution Principle (LSP)

Go does not use classical OOP inheritance heavily, but substitution still appears through compatible handler types:

- Fiber handlers are passed around through common signatures (`fiber.Handler`, `func(*fiber.Ctx) error`).
- Middleware like `JWTMiddleware(secret)` can be inserted wherever a `fiber.Handler` is expected.

So while not inheritance-based, behavior-level substitutability is present.

## I - Interface Segregation Principle (ISP)

There are no large, forced interfaces in this codebase. Instead, it uses small focused functions and package-level APIs:

- `util.GenerateJWT`, `util.Tokenize`, `util.ResponseAPI` each solve one small problem.
- Route setup functions accept only what they need (`router`, `secret`, `adminPass`).

This keeps modules simple and avoids over-coupled contracts.

## D - Dependency Inversion Principle (DIP)

**Partially applied**:

- Good: high-level wiring injects runtime values (JWT secret/admin pass) from [main.go](main.go) into route/controller flow.
- Current gap: controllers directly depend on concrete DB access (`mgm.Coll(...)`) instead of repository interfaces.

Future improvement: define repository interfaces in domain/service layer, then inject implementations.

## 3) DRY (Don't Repeat Yourself)

DRY is clearly used in utility abstractions:

- Standard API response generation via [util/apiResponse.util.go](util/apiResponse.util.go)
- Shared token generation for search indexing in [util/tokenizer.util.go](util/tokenizer.util.go)
- Shared JWT generation in [util/jwt.util.go](util/jwt.util.go)
- Reusable rate limiter setup in [util/rate_limiter.util.go](util/rate_limiter.util.go)

Without these helpers, controllers would repeat formatting/security/rate-limit logic.

## 4) Other Important Design Principles Used

## Separation of Concerns

- Routing, business logic, infra setup, and utilities are clearly separated.
- This is one of the strongest design qualities in the project.

## Composition over Inheritance

- Middleware composition in Fiber (`app.Use(...)`, route-specific middleware).
- Struct embedding (`mgm.DefaultModel`) in [models/data.models.go](models/data.models.go).

## Defensive Programming

- Input checks and fail-fast validation in controllers.
- Auth checks in middleware.
- Global error handler and graceful shutdown in [main.go](main.go).

## Consistent API Contract

- Most successful/error responses use a common shape through `ResponseAPI`, improving frontend predictability.

## 5) Practical Summary

Principles used well today:

- SRP
- DRY
- Separation of Concerns
- Composition over Inheritance
- Partial OCP/ISP

Principle that can be improved next:

- DIP (introduce service/repository interfaces to reduce controller-to-database coupling)

## 6) Suggested Next Refactor (If You Want Stronger OOP/SOLID)

1. Add service layer (`service/`) between controllers and DB.
2. Add repository interfaces for each domain (project, experience, certificate, volunteer).
3. Inject repository/service instances into handlers.
4. Keep controller focused on HTTP concerns only.

This would move the architecture from "clean modular" to "strongly SOLID with testable dependency boundaries".
