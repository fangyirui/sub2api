CREATE TABLE IF NOT EXISTS sms_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(100) NOT NULL,
    phone_number VARCHAR(50) NOT NULL DEFAULT '',
    hero_sms_id VARCHAR(100) NOT NULL DEFAULT '',
    sms_content TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'created',
    pending_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS sms_orders_order_no_key ON sms_orders (order_no);
CREATE INDEX IF NOT EXISTS idx_sms_orders_status ON sms_orders (status);
