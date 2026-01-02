# Domain Model — FoodApp (Phase 2 — API Contract Locking)

Use locked terms consistently. Do not introduce synonyms. No business logic here.

Locked Terms
- User, Customer, Admin, MenuItem, Category, Order, OrderItem, PaymentIntent
- OrderStatus: PLACED, CONFIRMED, PREPARING, OUT_FOR_DELIVERY, COMPLETED, CANCELLED

Global Conventions
- Identifiers: UUID v4
- Timestamps: createdAt, updatedAt (UTC, RFC3339)
- Currency: NGN (default)
- Timezone: Africa/Lagos (UTC+1)

Entities (structure only)
- User
  - id: UUID
  - email: string (unique)
  - passwordHash: string (bcrypt/argon2)
  - role: enum [Customer, Admin]
  - createdAt, updatedAt

- Customer
  - id: UUID (references User.id)
  - name: string
  - phone: string (E.164)
  - defaultAddressId: UUID (optional)
  - createdAt, updatedAt

- Admin
  - id: UUID (references User.id)
  - name: string
  - createdAt, updatedAt

- Category
  - id: UUID
  - name: string
  - description: string (optional)
  - sortOrder: int
  - isActive: boolean
  - createdAt, updatedAt

- MenuItem
  - id: UUID
  - name: string
  - description: string
  - categoryId: UUID
  - priceCents: int (includes tax)
  - currency: string (ISO 4217; default NGN)
  - isAvailable: boolean
  - createdAt, updatedAt

- Order
  - id: UUID
  - customerId: UUID (optional for guest)
  - status: OrderStatus
  - totalCents: int
  - currency: string (ISO 4217; default NGN)
  - deliveryAddressId: UUID (optional)
  - deliveryInstructions: string (optional)
  - tipCents: int (>= 0)
  - deliveryEtaMinutes: int (read-only; optional)
  - deliveredAt: timestamp (read-only; optional)
  - courierName: string (read-only; optional)
  - deliveryTrackingUrl: string (read-only; optional)
  - createdAt, updatedAt

- OrderItem
  - id: UUID
  - orderId: UUID
  - menuItemId: UUID
  - nameSnapshot: string (MenuItem name at time of order)
  - unitPriceCents: int
  - quantity: int
  - createdAt, updatedAt

- PaymentIntent
  - id: UUID
  - orderId: UUID (unique)
  - provider: string (abstracted; e.g., Paystack)
  - providerIntentId: string
  - status: string (provider-mapped)
  - clientSecret: string (optional; provider-dependent)
  - authorizationUrl: string (optional; provider-dependent)
  - amountCents: int
  - currency: string (ISO 4217; default NGN)
  - createdAt, updatedAt

Value Objects
- Money: amountCents (int), currency (ISO 4217)
- Address: line1, line2, city, state, postalCode, country
- PhoneNumber: e164 string
- OrderSecret: opaque string used to authorize guest access to a specific Order (not stored in Order response except on creation)

Aggregates
- Order aggregate: Order + [OrderItem...] + PaymentIntent reference
- Menu aggregate: Category + [MenuItem...]

Invariants (descriptive only)
- Order.totalCents = sum(OrderItem.quantity * unitPriceCents) + tipCents
- All OrderItems share same currency as Order.currency
- PaymentIntent.amountCents = Order.totalCents at intent creation (adjustments require new flow)
- One PaymentIntent per Order at a time

OrderStatus Transitions
- PLACED → CONFIRMED → PREPARING → OUT_FOR_DELIVERY → COMPLETED
- PLACED → CANCELLED
- CONFIRMED → CANCELLED (admin only, reason required)
- PREPARING → CANCELLED (exception only, reason required)
- OUT_FOR_DELIVERY → CANCELLED (exception only, reason required)
- COMPLETED and CANCELLED are terminal

Events (future)
- OrderPlaced, OrderConfirmed, OrderCancelled, OrderCompleted
- PaymentIntentCreated, PaymentSucceeded, PaymentFailed

Security & Auditing (domain requirements)
- Status changes attributable to a User (Admin) or system process
- Cancellations/refunds require reason and audit trail entry
