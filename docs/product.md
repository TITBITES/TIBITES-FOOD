# Product Definition — FoodApp (Phase 1 — Foundation)

Purpose: Establish clear scope and constraints before any feature implementation. No business logic here.

Business: FoodApp (placeholder)
Timezone: Africa/Lagos (UTC+1)
Default currency: NGN
Service mode: Delivery only (Phase 1). Pickup may be added later. Dine-in excluded.

## 1. Problem Statement
Provide a reliable, API-first system for a single food business to publish a menu and accept delivery orders. Customers can browse menu, place orders, tip, pay securely, and track order status. Admins manage menu and oversee order lifecycle.

## 2. Goals
- Deliver API-driven platform for a single business (not a marketplace)
- Strict separation between backend services and frontend (no business logic in frontend)
- Robust order lifecycle with auditable state transitions and provider-agnostic payments

## 3. Non-Goals (Phase 1)
- No multi-vendor marketplace
- No table booking or dine-in
- No direct DB access from handlers
- No backend enforcement of operating hours (frontend may enforce)
- No system-generated invoices (provider receipts suffice)

## 4. Target Users
- Customer: browses menu, places order, sets tipCents (>= 0), pays, tracks OrderStatus
- Admin: manages MenuItem and Category, confirms/cancels orders, monitors payments via provider dashboard

## 5. Core Use Cases (High level only)
- View Menu by Category (read-only API)
- Place Order with OrderItems and delivery address
- Initiate PaymentIntent and complete payment via provider (Paystack initially; abstraction preserved)
- Track OrderStatus: PLACED → CONFIRMED → PREPARING → OUT_FOR_DELIVERY → COMPLETED; or CANCELLED

## 6. Constraints
- API-first; contracts defined in OpenAPI before implementation; no breaking changes without versioning
- Clean Architecture; strict separation of concerns
- Data: PostgreSQL (primary), Redis (cache/queues)
- Auth: JWT (access + refresh), OAuth2 for admin
- Security: bcrypt/argon2, rate limiting, CORS/CSRF protection

## 7. Operational Requirements
- Idempotent creation endpoints
- Observability: structured logs; basic metrics (latency, error rate)
- Config via environment variables (12-factor)

## 8. Success Metrics (initial targets)
- 99.9% API uptime (Phase 2 target)
- p95 read latency < 300ms; p95 order create < 800ms (excluding provider latency)
- < 0.5% failed charges due to integration errors

## 9. Risks & Assumptions
- Single brand only; future marketplace would require new versioned APIs
- Payment provider abstraction; Paystack initially; Stripe-ready for future
- Delivery addresses accepted without geofencing/validation; operations enforce feasibility

## 10. Privacy & Compliance (baseline)
- Store minimal PII: name, phone, address for delivery
- PCI: payment details never touch our servers; handled by provider
- Data retention: Orders retained 24 months; PII removal on request

## 11. Phase 1 Scope Lock
- Delivery only
- Prices include tax (in priceCents)
- Tipping supported as free numeric entry (tipCents >= 0)
- No backend hour enforcement
- Provider-generated receipts only
