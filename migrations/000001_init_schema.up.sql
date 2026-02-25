-- Create businesses table
CREATE TABLE IF NOT EXISTS businesses (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'business',
    hmac_key VARCHAR(255) UNIQUE NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    created_at_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_businesses_email ON businesses(email);

-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at BIGINT NOT NULL,
    created_at_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);

-- Create services table
CREATE TABLE IF NOT EXISTS services (
    id VARCHAR(36) PRIMARY KEY,
    business_id VARCHAR(36) NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    total_slots INT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    created_at_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_services_business_id ON services(business_id);

-- Create slots table
CREATE TABLE IF NOT EXISTS slots (
    id VARCHAR(36) PRIMARY KEY,
    service_id VARCHAR(36) NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    slot_number INT NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('free', 'occupied')),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    created_at_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(service_id, slot_number)
);

CREATE INDEX IF NOT EXISTS idx_slots_service_id ON slots(service_id);
CREATE INDEX IF NOT EXISTS idx_slots_service_status ON slots(service_id, status);
CREATE INDEX IF NOT EXISTS idx_slots_status ON slots(status);

-- Create tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id VARCHAR(36) PRIMARY KEY,
    service_id VARCHAR(36) NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    slot_id VARCHAR(36) REFERENCES slots(id) ON DELETE SET NULL,
    slot_number INT NOT NULL,
    customer_id VARCHAR(36) REFERENCES customers(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'released')),
    hmac_digest VARCHAR(255),
    issued_at BIGINT NOT NULL,
    released_at BIGINT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    created_at_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tickets_service_id ON tickets(service_id);
CREATE INDEX IF NOT EXISTS idx_tickets_customer_id ON tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_hmac ON tickets(hmac_digest);
CREATE INDEX IF NOT EXISTS idx_tickets_slot_id ON tickets(slot_id);
