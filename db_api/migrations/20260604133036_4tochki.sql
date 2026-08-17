-- +goose Up

CREATE TABLE supplier_4tochki (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    canonical_code VARCHAR(255) UNIQUE,
    wid VARCHAR(255),
    hei VARCHAR(255),
    dia INTEGER,
    siz VARCHAR(255),
    lt VARCHAR(255) CHECK (lt IN ('l', 'lt', 'm', 't', 'z')),
    seas VARCHAR(255) CHECK (seas IN ('w', 's', 'as')),
    stud BOOLEAN CHECK (stud IN (0, 1)),
    brand VARCHAR(255),
    model VARCHAR(255),
    suv BOOLEAN CHECK (suv IN (0, 1)),
    amo_solonkl INTEGER,
    amo_kryarsk2 INTEGER,
    price_per_one_solonkl INTEGER,
    price_per_one_kryarsk2 INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS supplier_4tochki;
