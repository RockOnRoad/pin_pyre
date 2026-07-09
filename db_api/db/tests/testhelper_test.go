package dbTests

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

const testSchema = `
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
    FOREIGN KEY (own_code) REFERENCES stock_tyres(own_code) ON DELETE CASCADE,
    UNIQUE (supplier, supp_code)
);

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
`

func newTestStore(t *testing.T) *db.Store {
	t.Helper()

	conn, err := sqlx.Connect("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	conn.SetMaxOpenConns(1)

	if _, err := conn.Exec(testSchema); err != nil {
		t.Fatalf("apply test schema: %v", err)
	}

	t.Cleanup(func() { _ = conn.Close() })

	return db.NewStore(conn)
}

func strPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }

func boolPtr(b bool) *bool { return &b }
