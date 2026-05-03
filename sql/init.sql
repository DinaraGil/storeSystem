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
-- DROP TABLE IF EXISTS state CASCADE;
-- DROP TABLE IF EXISTS quantity_mistake CASCADE;
-- DROP TABLE IF EXISTS scan_mistake CASCADE;

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
    status             TEXT   NOT NULL default 'NEW',
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
--     shipment_number     BIGINT NOT NULL,
    status              TEXT   NOT NULL,
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
    expected_amount           INT  NOT NULL,
    real_amount INT DEFAULT 0,
    status TEXT NOT NULL default 'NEW',
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
    expected_amount           INT  NOT NULL,
    real_amount INT DEFAULT 0,
    status TEXT NOT NULL default 'NEW',
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
    rfid_id          TEXT NOT NULL UNIQUE,
    delivery_list_id INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
    supplier_id      INT REFERENCES counterparty (counterparty_id),
    name             TEXT NOT NULL,
    article          TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'STOCKED',
    created_by       INT REFERENCES worker (worker_id),
    updated_by       INT REFERENCES worker (worker_id),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table stock_balance (
    article TEXT PRIMARY KEY,
    quantity INT NOT NULL DEFAULT 0,
    reserved INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE stock_reservation (
--     reservation_id SERIAL PRIMARY KEY,
--     article TEXT NOT NULL,
--     shipment_list_id INT REFERENCES shipment_list(shipment_list_id) ON DELETE CASCADE,
--     reserved_quantity INT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
-- =========================
-- EVENT
-- =========================

CREATE TABLE event
(
    event_id   SERIAL PRIMARY KEY,
    rfid_id    TEXT,
    article    TEXT,
    scanner INT,
    is_in      BOOL,
    error      TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE report (
    report_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES worker(worker_id),
    report_type TEXT NOT NULL,
    file_name TEXT NOT NULL,
    object_id TEXT NOT NULL,
    bucket_name TEXT NOT NULL,
    date_from DATE,
    date_to DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================
-- STATE
-- =========================

-- CREATE TABLE state
-- (
--     state_id   SERIAL PRIMARY KEY,
--     item_id    INT REFERENCES item (item_id) ON DELETE CASCADE,
--     state_name TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
--


-- =========================
-- QUANTITY MISTAKE
-- =========================

-- CREATE TABLE quantity_mistake
-- (
--     quantity_mistake_id SERIAL PRIMARY KEY,
--     delivery_list_id    INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
--     diff                INT NOT NULL,
--     created_by          INT REFERENCES worker (worker_id),
--     created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- =========================
-- SCAN MISTAKE
-- =========================

-- CREATE TABLE scan_mistake
-- (
--     scan_mistake_id  SERIAL PRIMARY KEY,
--     event_id         INT REFERENCES event (event_id) ON DELETE CASCADE,
--     delivery_list_id INT REFERENCES delivery_list (delivery_list_id) ON DELETE CASCADE,
--     shipment_list_id INT REFERENCES shipment_list (shipment_list_id) ON DELETE CASCADE,
--     created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

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
VALUES ('admin', 'admin1', 'hash1', 1),
       ('Warehouse Worker 1', 'worker1', 'hash2', 2),
       ('Warehouse Worker 2', 'worker2', 'hash3', 2),
       ('admin2 User', 'admin2', 'hash4', 1),
       ('admin3 User', 'admin3', 'hash5', 1) ON CONFLICT DO NOTHING;


-- =========================
-- COUNTERPARTY
-- =========================

INSERT INTO counterparty (
    full_name,
    legal_form,
    inn,
    kpp,
    ogrn,
    legal_address,
    bank_name,
    bik,
    bank_account_number,
    contact_person,
    phone,
    role_id
) VALUES
-- Поставщики (role_id = 1)
(
    'ООО "ТехноПоставка"',
    'ООО',
    '7701123456',
    '770101001',
    '1027700123456',
    'г. Москва, ул. Ленина, д. 10, офис 305',
    'АО "Альфа-Банк"',
    '044525593',
    '40702810123450000001',
    'Иванов Иван Иванович',
    '4951234567',
    1
),
(
    'АО "ПромСнаб"',
    'АО',
    '7802234567',
    '780201002',
    '1037800123457',
    'г. Санкт-Петербург, пр. Невский, д. 45, лит. А',
    'ПАО "Сбербанк"',
    '044030001',
    '40702810700000000002',
    'Петров Петр Петрович',
    '8122345678',
    1
),
-- Покупатели (role_id = 2)
(
    'ООО "РитейлПлюс"',
    'ООО',
    '7705123457',
    '770101005',
    '1077700123460',
    'г. Москва, ул. Новый Арбат, д. 20',
    'АО "ЮниКредит Банк"',
    '044525545',
    '40702810200000000006',
    'Алексеева Анна Викторовна',
    '4957654321',
    2
),
(
    'АО "Торговый Дом"',
    'АО',
    '7806234568',
    '780201006',
    '1087800123461',
    'г. Санкт-Петербург, ул. Садовая, д. 30',
    'ПАО "Промсвязьбанк"',
    '044030002',
    '40702810800000000007',
    'Николаев Николай Николаевич',
    '8123456789',
    2
);

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
    (status, created_by, completed_by)
VALUES ('NEW', 1, NULL),
       ('NEW', 2, NULL) ON CONFLICT DO NOTHING;

-- =========================
-- DELIVERY LIST
-- =========================

INSERT INTO delivery_list
    (delivery_id, supplier_id, expected_amount, article, created_by)
VALUES (1, 1, 10, 'ART-001', 1),
       (1, 2, 15, 'ART-002', 2),
       (2, 1, 20, 'ART-003', 1),
       (3, 3, 5, 'ART-004', 3),
       (4, 4, 12, 'ART-005', 2) ON CONFLICT DO NOTHING;

INSERT INTO shipment_list
    (shipment_id, customer_id, expected_amount, article, created_by)
VALUES (1, 2, 5, 'ART-001', 1),
       (2, 3, 7, 'ART-002', 2),
       (1, 4, 10, 'ART-003', 3),
       (2, 2, 3, 'ART-004', 1),
       (1, 1, 8, 'ART-005', 2) ON CONFLICT DO NOTHING;


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
-- EVENT NOTIFY TRIGGER
-- =========================

CREATE OR REPLACE FUNCTION notify_event_insert() RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify(
    'event_channel',
    json_build_object(
      'event_id', NEW.event_id,
      'rfid_id', NEW.rfid_id,
      'article', NEW.article,
      'scanner', NEW.scanner,
      'is_in', NEW.is_in,
      'error', NEW.error,
      'created_at', NEW.created_at
    )::text
  );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER event_insert_trigger
AFTER INSERT ON event
FOR EACH ROW
EXECUTE FUNCTION notify_event_insert();


-- =========================
-- STATE
-- =========================

-- INSERT INTO state (item_id, state_name)
-- VALUES (1, 'STORED'),
--        (2, 'IN_TRANSIT'),
--        (3, 'STORED'),
--        (4, 'DAMAGED'),
--        (5, 'STORED') ON CONFLICT DO NOTHING;


-- =========================
-- QUANTITY MISTAKE
-- =========================

-- INSERT INTO quantity_mistake (delivery_list_id, diff, created_by)
-- VALUES (1, 1, 1),
--        (2, -1, 2),
--        (3, 2, 3),
--        (4, -2, 1),
--        (5, 1, 2) ON CONFLICT DO NOTHING;
--
--
-- -- =========================
-- -- SCAN MISTAKE
-- -- =========================
--
-- INSERT INTO scan_mistake (event_id, delivery_list_id, shipment_list_id)
-- VALUES (1, 1, 1),
--        (2, 2, 2),
--        (3, 3, 3),
--        (4, 4, 4),
--        (5, 5, 5) ON CONFLICT DO NOTHING;