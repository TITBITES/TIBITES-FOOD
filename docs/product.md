# Product Definition — FoodApp (Phase 2 — API Contract Locking)

Purpose: Lock product scope and API contracts without business logic or implementation. This is a single-business delivery app (not a marketplace).

Business: FoodApp (placeholder)
Timezone: Africa/Lagos (UTC+1)
Default currency: NGN
Service mode: Delivery only (Phase 2). Pickup may be added later. Dine-in excluded.

Key Decisions (Locked)
- Guest checkout: Allowed using an orderSecret returned at order creation.
- Cart: Client-side only (no server-side Cart APIs in Phase 2).
- Payments: Online only via provider (Paystack initially). Abstraction remains provider-agnostic (Stripe-ready).
- Admin authentication: OAuth2 only.
- OrderStatus vs Payment state: OrderStatus remains as-is; payment lifecycle lives in PaymentIntent.status.
- Delivery enforcement: No geofencing or operating-hours enforcement in backend.
- Receipts: Provider-generated receipts are sufficient.

1. Problem Statement
Provide a reliable, API-first system for a single food business to publish a menu and accept delivery orders. Customers (including guests) can browse menu, place orders, tip, pay securely online, and track delivery status. Admins manage menu and oversee orders.

2. Goals
- API-first and versioned when necessary; no breaking changes without v2.
- Clean Architecture across future implementation phases.
- Provider-agnostic payments via a PaymentIntent abstraction.

3. Non-Goals (Phase 2)
- No server-side Cart
- No dine-in, table booking, or rider management
- No backend operating-hours enforcement
- No system-generated invoices

4. Target Users
- Customer (guest or authenticated): browse menu, create order, view order with orderSecret, initiate payment, track order delivery status (read-only)
- Admin: manage MenuItem/Category; confirm/cancel orders; monitor payments via provider dashboard (OAuth2 only)

5. Core Use Cases (High level)
- Browse Menu by Category (read-only)
- Create Order (guest allowed) and receive orderSecret
- Create and retrieve PaymentIntent using Bearer token or orderSecret
- Track OrderStatus and read-only delivery status fields

6. Constraints
- Data: PostgreSQL primary, Redis for cache/queues
- Auth: JWT (access + refresh) for customers; OAuth2 for admins
- Security: bcrypt/argon2, rate limiting, CORS/CSRF protection

7. Operational Requirements
- Idempotent creation endpoints
- Structured logs and basic metrics
- Environment-based configuration

8. Success Metrics (unchanged from Phase 1)
- 99.9% API uptime target; latency and error-rate targets as defined previously

9. Risks & Assumptions
- Single business only; expansion to multi-vendor requires new versioned APIs
- Payment details handled by provider; no PCI scope for card data
- Delivery addresses accepted without geofencing; operational feasibility handled offline

10. Phase 2 Scope Lock
- Guest checkout via orderSecret
- Client-side cart only
- Online payments only (Paystack first; abstraction maintained)
- Admin OAuth2 only
- OrderStatus fixed set; payment state in PaymentIntent.status only
- Read-only delivery status exposed on Order
