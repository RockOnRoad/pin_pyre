package dbTests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func sampleSupplier4Tochki(code string) *db.Supplier4Tochki {
	return &db.Supplier4Tochki{
		CanonicalCode:       strPtr(code),
		Wid:                 strPtr("205"),
		Hei:                 strPtr("55"),
		Dia:                 intPtr(1600),
		Siz:                 strPtr("2055516"),
		Lt:                  strPtr("l"),
		Seas:                strPtr("s"),
		Stud:                boolPtr(false),
		Brand:               strPtr("Michelin"),
		Model:               strPtr("Pilot Sport 4"),
		SUV:                 boolPtr(false),
		AmoSolonkl:          intPtr(4),
		AmoKryarsk2:         intPtr(2),
		PricePerOneSolonkl:  intPtr(850000),
		PricePerOneKryarsk2: intPtr(900000),
	}
}

func TestCreateSupplier4Tochki(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("ABC123"))
	if err != nil {
		t.Fatalf("CreateSupplier4Tochki: %v", err)
	}

	got, err := store.GetSupplier4TochkiByID(ctx, id)
	if err != nil {
		t.Fatalf("GetSupplier4TochkiByID: %v", err)
	}
	if got.CanonicalCode == nil || *got.CanonicalCode != "ABC123" {
		t.Fatalf("canonical_code = %v, want ABC123", got.CanonicalCode)
	}
	if got.AmoSolonkl == nil || *got.AmoSolonkl != 4 {
		t.Fatalf("amo_solonkl = %v, want 4", got.AmoSolonkl)
	}
}

func TestCreateSupplier4TochkiDuplicateCanonicalCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	if _, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("DUP-CODE")); err != nil {
		t.Fatalf("first create: %v", err)
	}

	if _, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("DUP-CODE")); err == nil {
		t.Fatal("expected error for duplicate canonical_code")
	}
}

func TestGetSupplier4TochkiByCanonicalCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	if _, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("CAN-001")); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetSupplier4TochkiByCanonicalCode(ctx, "CAN-001")
	if err != nil {
		t.Fatalf("GetSupplier4TochkiByCanonicalCode: %v", err)
	}
	if got.Brand == nil || *got.Brand != "Michelin" {
		t.Fatalf("brand = %v, want Michelin", got.Brand)
	}
}

func TestGetSupplier4TochkiByIDNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	_, err := store.GetSupplier4TochkiByID(ctx, 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestListSupplier4Tochki(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	for _, code := range []string{"L-001", "L-002", "L-003"} {
		if _, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki(code)); err != nil {
			t.Fatalf("create %s: %v", code, err)
		}
	}

	all, err := store.ListSupplier4Tochki(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListSupplier4Tochki: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d, want 3", len(all))
	}

	page, err := store.ListSupplier4Tochki(ctx, 1, 2)
	if err != nil {
		t.Fatalf("ListSupplier4Tochki page: %v", err)
	}
	if len(page) != 1 {
		t.Fatalf("page len = %d, want 1", len(page))
	}
	if page[0].CanonicalCode == nil || *page[0].CanonicalCode != "L-003" {
		t.Fatalf("page item = %+v, want L-003", page[0])
	}
}

func TestUpdateSupplier4Tochki(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("UPD-001"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	before, err := store.GetSupplier4TochkiByID(ctx, id)
	if err != nil {
		t.Fatalf("get before update: %v", err)
	}

	updated := sampleSupplier4Tochki("UPD-001")
	updated.ID = id
	updated.Brand = strPtr("Nokian")
	updated.PricePerOneSolonkl = intPtr(7900)

	if err := store.UpdateSupplier4Tochki(ctx, updated); err != nil {
		t.Fatalf("UpdateSupplier4Tochki: %v", err)
	}

	after, err := store.GetSupplier4TochkiByID(ctx, id)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if after.Brand == nil || *after.Brand != "Nokian" {
		t.Fatalf("brand = %v, want Nokian", after.Brand)
	}
	if after.PricePerOneSolonkl == nil || *after.PricePerOneSolonkl != 7900 {
		t.Fatalf("price_per_one_solonkl = %v, want 7900", after.PricePerOneSolonkl)
	}
	if !after.UpdatedAt.After(before.UpdatedAt) && !after.UpdatedAt.Equal(before.UpdatedAt) {
		t.Fatalf("updated_at did not advance")
	}
}

func TestUpdateSupplier4TochkiNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	rec := sampleSupplier4Tochki("GHOST")
	rec.ID = 404

	err := store.UpdateSupplier4Tochki(ctx, rec)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteSupplier4Tochki(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateSupplier4Tochki(ctx, sampleSupplier4Tochki("DEL-001"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.DeleteSupplier4Tochki(ctx, id); err != nil {
		t.Fatalf("DeleteSupplier4Tochki: %v", err)
	}

	_, err = store.GetSupplier4TochkiByID(ctx, id)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteSupplier4TochkiNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	err := store.DeleteSupplier4Tochki(ctx, 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}
