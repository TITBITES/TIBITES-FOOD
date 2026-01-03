-- Phase 5.0 initial schema
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);
CREATE INDEX IF NOT EXISTS idx_categories_active_order ON categories(is_active, sort_order, name);

CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    price_cents INT NOT NULL,
    currency TEXT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true
);
CREATE INDEX IF NOT EXISTS idx_menu_items_avail_name ON menu_items(is_available, name);
CREATE INDEX IF NOT EXISTS idx_menu_items_category ON menu_items(category_id);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID,
    status TEXT NOT NULL,
    total_cents INT NOT NULL,
    currency TEXT NOT NULL,
    delivery_address_id UUID,
    delivery_instructions TEXT,
    tip_cents INT NOT NULL DEFAULT 0,
    delivery_eta_minutes INT,
    delivered_at TIMESTAMPTZ,
    courier_name TEXT,
    delivery_tracking_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_orders_customer_created ON orders(customer_id, created_at);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL,
    name_snapshot TEXT NOT NULL DEFAULT '',
    unit_price_cents INT NOT NULL DEFAULT 0,
    quantity INT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

CREATE TABLE IF NOT EXISTS order_secrets (
    order_id UUID PRIMARY KEY REFERENCES orders(id) ON DELETE CASCADE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_intent_id TEXT NOT NULL,
    status TEXT NOT NULL,
    client_secret TEXT,
    authorization_url TEXT,
    amount_cents INT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
