# Domain Model — FoodApp (Phase 1 — Foundation)

Use locked terms consistently. Do not introduce synonyms. No business logic here.

## Locked Terms
- User
- Customer
- Admin
- MenuItem
- Category
- Order
- OrderItem
- PaymentIntent
- OrderStatus: PLACED, CONFIRMED, PREPARING, OUT_FOR_DELIVERY, COMPLETED, CANCELLED

## Global Conventions
- Identifiers: UUID v4
- Timestamps: createdAt, updatedAt (UTC, RFC3339)
- Currency: NGN (default)
- Timezone: Africa/Lagos (UTC+1)

## Entities (structure only)

### User
- id: UUID
- email: string (unique)
- passwordHash: string (bcrypt/argon2)
- role: enum [Customer, Admin]
- createdAt: timestamp
- updatedAt: timestamp

### Customer
- id: UUID (references User.id)
- name: string
- phone: string (E.164)
- defaultAddressId: UUID (optional)
- createdAt: timestamp
- updatedAt: timestamp

### Admin
- id: UUID (references User.id)
- name: string
- createdAt: timestamp
- updatedAt: timestamp

### Category
- id: UUID
- name: string
- description: string (optional)
- sortOrder: int
- isActive: boolean
- createdAt: timestamp
- updatedAt: timestamp

### MenuItem
- id: UUID
- name: string
- description: string
- categoryId: UUID
- priceCents: int (includes tax)
- currency: string (ISO 4217; default NGN)
- isAvailable: boolean
- createdAt: timestamp
- updatedAt: timestamp

### Order
- id: UUID
- customerId: UUID
- status: OrderStatus
- totalCents: int
- currency: string (ISO 4217; default NGN)
- deliveryAddressId: UUID (optional)
- deliveryInstructions: string (optional)
- tipCents: int (>= 0)
- createdAt: timestamp
- updatedAt: timestamp

### OrderItem
- id: UUID
- orderId: UUID
- menuItemId: UUID
- nameSnapshot: string (MenuItem name at time of order)
- unitPriceCents: int
- quantity: int
- createdAt: timestamp
- updatedAt: timestamp

### PaymentIntent
- id: UUID
- orderId: UUID (unique)
- provider: string (abstracted; e.g., Paystack)
- providerIntentId: string
- status: string (provider-mapped)
- clientSecret: string
- amountCents: int
- currency: string (ISO 4217; default NGN)
- createdAt: timestamp
- updatedAt: timestamp

## Value Objects
- Money: amountCents (int), currency (ISO 4217)
- Address: line1, line2, city, state, postalCode, country
- PhoneNumber: e164 string

## Aggregates
- Order aggregate: Order + [OrderItem...] + PaymentIntent reference
- Menu aggregate: Category + [MenuItem...]

## Invariants (descriptive only)
- Order.totalCents = sum(OrderItem.quantity * unitPriceCents) + tipCents
- All OrderItems share same currency as Order.currency
- PaymentIntent.amountCents = Order.totalCents at intent creation (adjustments require new flow)
- One PaymentIntent per Order at a time
- Order Customer must exist and be role Customer

## OrderStatus Transitions
- PLACED → CONFIRMED → PREPARING → OUT_FOR_DELIVERY → COMPLETED
- PLACED → CANCELLED
- CONFIRMED → CANCELLED (admin only, reason required)
- PREPARING → CANCELLED (exception only, reason required)
- OUT_FOR_DELIVERY → CANCELLED (exception only, reason required)
- COMPLETED and CANCELLED are terminal

## Events (future)
- OrderPlaced, OrderConfirmed, OrderCancelled, OrderCompleted
- PaymentIntentCreated, PaymentSucceeded, PaymentFailed

## Security & Auditing (domain requirements)
- Status changes attributable to a User (Admin) or system process
- Cancellations/refunds require reason and audit trail entry
