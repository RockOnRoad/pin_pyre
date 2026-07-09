-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE stock_tyres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    own_code VARCHAR(255) NOT NULL UNIQUE,

    wid VARCHAR(255),
    hei VARCHAR(255),
    dia INTEGER,
    idx VARCHAR(255),
    xl BOOLEAN CHECK (xl IN (0, 1)),
    siz VARCHAR(255),
    lt VARCHAR(255) CHECK (lt IN ('l', 'lt', 'm', 't', 'z')),
    seas VARCHAR(255) CHECK (seas IN ('w', 's', 'as')),
    stud BOOLEAN CHECK (stud IN (0, 1)),
    brand VARCHAR(255),
    model VARCHAR(255),
    suv BOOLEAN CHECK (suv IN (0, 1)),

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tyre_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    supplier VARCHAR(255) NOT NULL,
    supp_code VARCHAR(255) NOT NULL,
    own_code VARCHAR(255) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (own_code)
        REFERENCES stock_tyres(own_code)
        ON DELETE CASCADE,

    UNIQUE (supplier, supp_code)
);

CREATE INDEX idx_stock_tyre_codes_own_code
    ON tyre_codes(own_code);

CREATE INDEX idx_stock_tyre_codes_supp_code
    ON tyre_codes(supp_code);


-- +goose Down
DROP INDEX IF EXISTS idx_stock_tyre_codes_own_code;
DROP INDEX IF EXISTS idx_stock_tyre_codes_supp_code;

DROP TABLE IF EXISTS tyre_codes;
DROP TABLE IF EXISTS stock_tyres;