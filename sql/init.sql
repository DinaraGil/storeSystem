DROP TABLE IF EXISTS worker_role CASCADE;
DROP TABLE IF EXISTS counterparty_role CASCADE;
DROP TABLE IF EXISTS worker CASCADE;
DROP TABLE IF EXISTS counterparty CASCADE;
DROP TABLE IF EXISTS delivery CASCADE;
DROP TABLE IF EXISTS shipment CASCADE;
DROP TABLE IF EXISTS delivery_list CASCADE;
DROP TABLE IF EXISTS shipment_list CASCADE;
DROP TABLE IF EXISTS item CASCADE;
DROP TABLE IF EXISTS event CASCADE;
DROP TABLE IF EXISTS state CASCADE;
DROP TABLE IF EXISTS quantity_mistake CASCADE;
DROP TABLE IF EXISTS scan_mistake CASCADE;

-- =========================
-- COUNTERPARTY ROLE
-- =========================

CREATE TABLE counterparty_role
(
    counterparty_role_id SERIAL PRIMARY KEY,
    role_name            TEXT NOT NULL
);

-- =========================
-- WORKER ROLE
-- =========================

CREATE TABLE worker_role
(
    worker_role_id SERIAL PRIMARY KEY,
    role_name      TEXT NOT NULL
);

-- =========================
-- COUNTERPARTY
-- =========================

CREATE TABLE counterparty
(
    counterparty_id     SERIAL PRIMARY KEY,
    full_name           TEXT        NOT NULL,
    legal_form          TEXT        NOT NULL,
    inn                 VARCHAR(12),
    kpp                 VARCHAR(9),
    ogrn                VARCHAR(13),
    legal_address       TEXT        NOT NULL,
    bank_name           TEXT        NOT NULL,
    bik                 VARCHAR(9)  NOT NULL,
    bank_account_number VARCHAR(20) NOT NULL,
    contact_person      TEXT        NOT NULL,
    phone               VARCHAR(11),
    role_id             INT REFERENCES counterparty_role (counterparty_role_id)
);

-- =========================
-- WORKER
-- =========================

CREATE TABLE worker
(
    worker_id     SERIAL PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    role_id       INT REFERENCES worker_role (worker_role_id)
);

-- =========================
-- DELIVERY
-- =========================

CREATE TABLE delivery
(
    delivery_id        SERIAL PRIMARY KEY,
    status             TEXT   NOT NULL,
    planned_arrival_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by         INT REFERENCES worker (worker_id),
    accepted_by        INT REFERENCES worker (worker_id),
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- SHIPMENT
-- =========================

CREATE TABLE shipment
(
    shipment_id         SERIAL PRIMARY KEY,
    shipment_number     BIGINT NOT NULL,
    status              TEXT   NOT NULL,
    planned_shipment_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by          INT REFERENCES worker (worker_id),
    completed_by        INT REFERENCES worker (worker_id),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- DELIVERY LIST
-- =========================

CREATE TABLE delivery_list
(
    delivery_list_id SERIAL PRIMARY KEY,
    delivery_id      INT REFERENCES delivery (delivery_id) ON DELETE CASCADE,
    supplier_id      INT REFERENCES counterparty (counterparty_id),
    amount           INT  NOT NULL,
    article          TEXT NOT NULL,
    created_by       INT REFERENCES worker (worker_id),
    updated_by       INT REFERENCES worker (worker_id),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- SHIPMENT LIST
-- =========================

CREATE TABLE shipment_list
(
    shipment_list_id SERIAL PRIMARY KEY,
    shipment_id      INT REFERENCES shipment (shipment_id) ON DELETE CASCADE,
    customer_id      INT REFERENCES counterparty (counterparty_id),
    supplier_id      INT REFERENCES counterparty (counterparty_id),
    amount           INT  NOT NULL,
    article          TEXT NOT NULL,
    created_by       INT REFERENCES worker (worker_id),
    updated_by       INT REFERENCES worker (worker_id),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- ITEM
-- =========================

CREATE TABLE item
(
    item_id          SERIAL PRIMARY KEY,
    rfid_id          TEXT NOT NULL,
    delivery_list_id INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
    supplier_id      INT REFERENCES counterparty (counterparty_id),
    name             TEXT NOT NULL,
    article          TEXT NOT NULL,
    created_by       INT REFERENCES worker (worker_id),
    updated_by       INT REFERENCES worker (worker_id),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- EVENT
-- =========================

CREATE TABLE event
(
    event_id   SERIAL PRIMARY KEY,
    rfid_id    TEXT NOT NULL,
    article    TEXT,
    is_in      BOOL,
    error      TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- STATE
-- =========================

CREATE TABLE state
(
    state_id   SERIAL PRIMARY KEY,
    item_id    INT REFERENCES item (item_id) ON DELETE CASCADE,
    state_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- QUANTITY MISTAKE
-- =========================

CREATE TABLE quantity_mistake
(
    quantity_mistake_id SERIAL PRIMARY KEY,
    delivery_list_id    INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
    diff                INT NOT NULL,
    created_by          INT REFERENCES worker (worker_id),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- SCAN MISTAKE
-- =========================

CREATE TABLE scan_mistake
(
    scan_mistake_id  SERIAL PRIMARY KEY,
    event_id         INT REFERENCES event (event_id) ON DELETE CASCADE,
    delivery_list_id INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
    shipment_list_id INT REFERENCES shipment_list (shipment_list_id) ON DELETE CASCADE,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- WORKER ROLE SEED
-- =========================

INSERT INTO worker_role (role_name)
VALUES ('ADMIN'),
       ('WAREHOUSE_WORKER');


-- =========================
-- COUNTERPARTY ROLE SEED
-- =========================

INSERT INTO counterparty_role (role_name)
VALUES ('SUPPLIER'),
       ('CUSTOMER');


-- =========================
-- WORKER
-- =========================

INSERT INTO worker (full_name, username, password_hash, role_id)
VALUES ('Admin User', 'admin1', 'hash1', 1),
       ('Warehouse Worker 1', 'worker1', 'hash2', 2),
       ('Warehouse Worker 2', 'worker2', 'hash3', 2),
       ('admin2 User', 'admin2', 'hash4', 1),
       ('admin3 User', 'admin3', 'hash5', 1) ON CONFLICT DO NOTHING;


-- =========================
-- COUNTERPARTY
-- =========================

INSERT INTO counterparty
(full_name, legal_form, legal_address, bank_name, bik,
 bank_account_number, contact_person, role_id)
VALUES ('Supplier One LLC', 'LLC', 'Address 1', 'Bank 1', '111111111',
        '11111111111111111111', 'Contact 1', 1),

       ('Customer One LLC', 'LLC', 'Address 2', 'Bank 2', '222222222',
        '22222222222222222222', 'Contact 2', 2),

       ('Customer Partner LLC', 'LLC', 'Address 3', 'Bank 3', '333333333',
        '33333333333333333333', 'Contact 3', 1),

       ('Customer Partner LLC', 'LLC', 'Address 4', 'Bank 4', '444444444',
        '44444444444444444444', 'Contact 4', 2),

       ('General Supplier LLC', 'LLC', 'Address 5', 'Bank 5', '555555555',
        '55555555555555555555', 'Contact 5', 1) ON CONFLICT DO NOTHING;


-- =========================
-- DELIVERY
-- =========================

INSERT INTO delivery
    (status, created_by, accepted_by)
VALUES ('NEW', 1, 2),
       ('NEW', 1, 2),
       ('IN_PROGRESS', 2, 3),
       ('COMPLETED', 3, 2),
       ('NEW', 1, 2) ON CONFLICT DO NOTHING;


-- =========================
-- SHIPMENT
-- =========================

INSERT INTO shipment
    (shipment_number, status, created_by, completed_by)
VALUES (2001, 'NEW', 1, NULL),
       (2002, 'NEW', 2, NULL),
       (2003, 'IN_PROGRESS', 3, NULL),
       (2004, 'COMPLETED', 2, 1),
       (2005, 'NEW', 1, NULL) ON CONFLICT DO NOTHING;


-- =========================
-- DELIVERY LIST
-- =========================

INSERT INTO delivery_list
    (delivery_id, supplier_id, amount, article, created_by)
VALUES (1, 1, 10, 'ART-001', 1),
       (1, 2, 15, 'ART-002', 2),
       (2, 1, 20, 'ART-003', 1),
       (3, 3, 5, 'ART-004', 3),
       (4, 4, 12, 'ART-005', 2) ON CONFLICT DO NOTHING;


-- =========================
-- SHIPMENT LIST
-- =========================

INSERT INTO shipment_list
    (shipment_id, customer_id, supplier_id, amount, article, created_by)
VALUES (1, 2, 1, 5, 'ART-001', 1),
       (2, 3, 1, 7, 'ART-002', 2),
       (3, 4, 2, 10, 'ART-003', 3),
       (4, 2, 3, 3, 'ART-004', 1),
       (5, 1, 4, 8, 'ART-005', 2) ON CONFLICT DO NOTHING;


-- =========================
-- ITEM (RFID ITEMS)
-- =========================

INSERT INTO item
    (rfid_id, delivery_list_id, supplier_id, name, article, created_by)
VALUES ('RFID001', 1, 1, 'Item 1', 'ART-001', 1),
       ('RFID002', 2, 2, 'Item 2', 'ART-002', 2),
       ('RFID003', 3, 1, 'Item 3', 'ART-003', 1),
       ('RFID004', 4, 3, 'Item 4', 'ART-004', 3),
       ('RFID005', 5, 4, 'Item 5', 'ART-005', 2) ON CONFLICT DO NOTHING;


-- =========================
-- EVENT
-- =========================

INSERT INTO event (rfid_id, article, is_in, error)
VALUES ('RFID001', 'ART-001', TRUE, NULL),
       ('RFID002', 'ART-002', TRUE, NULL),
       ('RFID003', 'ART-003', FALSE, 'Scan error'),
       ('RFID004', 'ART-004', TRUE, NULL),
       ('RFID005', 'ART-005', FALSE, NULL) ON CONFLICT DO NOTHING;


-- =========================
-- STATE
-- =========================

INSERT INTO state (item_id, state_name)
VALUES (1, 'STORED'),
       (2, 'IN_TRANSIT'),
       (3, 'STORED'),
       (4, 'DAMAGED'),
       (5, 'STORED') ON CONFLICT DO NOTHING;


-- =========================
-- QUANTITY MISTAKE
-- =========================

INSERT INTO quantity_mistake (delivery_list_id, diff, created_by)
VALUES (1, 1, 1),
       (2, -1, 2),
       (3, 2, 3),
       (4, -2, 1),
       (5, 1, 2) ON CONFLICT DO NOTHING;


-- =========================
-- SCAN MISTAKE
-- =========================

INSERT INTO scan_mistake (event_id, delivery_list_id, shipment_list_id)
VALUES (1, 1, 1),
       (2, 2, 2),
       (3, 3, 3),
       (4, 4, 4),
       (5, 5, 5) ON CONFLICT DO NOTHING;