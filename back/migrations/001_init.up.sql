CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      VARCHAR(50)  UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(150) NOT NULL,
    email         VARCHAR(150) UNIQUE NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('admin','inventory','viewer')),
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE departments (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(150) UNIQUE NOT NULL,
    location   VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE equipment (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_number      VARCHAR(50)  UNIQUE NOT NULL,
    name                  VARCHAR(255) NOT NULL,
    description           TEXT,
    category_id           UUID NOT NULL REFERENCES categories(id),
    serial_number         VARCHAR(100),
    model                 VARCHAR(150),
    manufacturer          VARCHAR(150),
    purchase_date         DATE,
    purchase_price        NUMERIC(12,2),
    warranty_expiry       DATE,
    status                VARCHAR(20) NOT NULL DEFAULT 'in_storage'
                              CHECK (status IN ('in_use','in_storage','in_repair','written_off','reserved')),
    department_id         UUID NOT NULL REFERENCES departments(id),
    responsible_person_id UUID REFERENCES users(id),
    notes                 TEXT,
    is_archived           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE equipment_photos (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    file_path    VARCHAR(500) NOT NULL,
    uploaded_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE movements (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    equipment_id     UUID NOT NULL REFERENCES equipment(id),
    from_department_id UUID NOT NULL REFERENCES departments(id),
    to_department_id   UUID NOT NULL REFERENCES departments(id),
    moved_by         UUID NOT NULL REFERENCES users(id),
    moved_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reason           TEXT
);

CREATE TABLE inventory_sessions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_id UUID NOT NULL REFERENCES departments(id),
    status        VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                      CHECK (status IN ('in_progress','completed')),
    created_by    UUID NOT NULL REFERENCES users(id),
    started_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at   TIMESTAMPTZ
);

CREATE TABLE inventory_items (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id      UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
    equipment_id    UUID NOT NULL REFERENCES equipment(id),
    expected_status VARCHAR(20) NOT NULL,
    actual_status   VARCHAR(20) NOT NULL CHECK (actual_status IN ('found','not_found','damaged')),
    comment         TEXT,
    checked_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id),
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID,
    details     JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_equipment_category      ON equipment(category_id);
CREATE INDEX idx_equipment_department    ON equipment(department_id);
CREATE INDEX idx_equipment_status        ON equipment(status);
CREATE INDEX idx_equipment_inv_number    ON equipment(inventory_number);
CREATE INDEX idx_equipment_is_archived   ON equipment(is_archived);
CREATE INDEX idx_movements_equipment     ON movements(equipment_id);
CREATE INDEX idx_movements_moved_at      ON movements(moved_at);
CREATE INDEX idx_inv_items_session       ON inventory_items(session_id);
CREATE INDEX idx_audit_logs_user         ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity       ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created      ON audit_logs(created_at);
