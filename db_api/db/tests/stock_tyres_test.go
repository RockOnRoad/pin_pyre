package dbTests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func sampleStockTyre(ownCode string) *db.StockTyre {
	return &db.StockTyre{
		OwnCode: ownCode,
		Wid:     strPtr("205"),
		Hei:     strPtr("55"),
		Dia:     intPtr(16),
		Idx:     strPtr("91"),
		XL:      boolPtr(false),
		Siz:     strPtr("205/55R16"),
		Lt:      strPtr("l"),
		Seas:    strPtr("s"),
		Stud:    boolPtr(false),
		Brand:   strPtr("Michelin"),
		Model:   strPtr("Pilot Sport 4"),
		SUV:     boolPtr(false),
	}
}

func TestCreateStockTyre(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateStockTyre(ctx, sampleStockTyre("MIC-205-55-16"))
	if err != nil {
		t.Fatalf("CreateStockTyre: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero id")
	}

	got, err := store.GetStockTyreByID(ctx, id)
	if err != nil {
		t.Fatalf("GetStockTyreByID: %v", err)
	}
	if got.OwnCode != "MIC-205-55-16" {
		t.Fatalf("own_code = %q, want MIC-205-55-16", got.OwnCode)
	}
	if got.Brand == nil || *got.Brand != "Michelin" {
		t.Fatalf("brand = %v, want Michelin", got.Brand)
	}
}

func TestCreateStockTyreDuplicateOwnCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	if _, err := store.CreateStockTyre(ctx, sampleStockTyre("DUP-001")); err != nil {
		t.Fatalf("first create: %v", err)
	}

	if _, err := store.CreateStockTyre(ctx, sampleStockTyre("DUP-001")); err == nil {
		t.Fatal("expected error for duplicate own_code")
	}
}

func TestGetStockTyreByIDNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	_, err := store.GetStockTyreByID(ctx, 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestGetStockTyreByOwnCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	if _, err := store.CreateStockTyre(ctx, sampleStockTyre("OWN-123")); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetStockTyreByOwnCode(ctx, "OWN-123")
	if err != nil {
		t.Fatalf("GetStockTyreByOwnCode: %v", err)
	}
	if got.OwnCode != "OWN-123" {
		t.Fatalf("own_code = %q, want OWN-123", got.OwnCode)
	}
}

func TestListStockTyres(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	for _, code := range []string{"A-001", "B-002", "C-003"} {
		if _, err := store.CreateStockTyre(ctx, sampleStockTyre(code)); err != nil {
			t.Fatalf("create %s: %v", code, err)
		}
	}

	all, err := store.ListStockTyres(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListStockTyres: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d, want 3", len(all))
	}

	page, err := store.ListStockTyres(ctx, 2, 1)
	if err != nil {
		t.Fatalf("ListStockTyres page: %v", err)
	}
	if len(page) != 2 {
		t.Fatalf("page len = %d, want 2", len(page))
	}
	if page[0].OwnCode != "B-002" {
		t.Fatalf("first item own_code = %q, want B-002", page[0].OwnCode)
	}
}

func TestUpdateStockTyre(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateStockTyre(ctx, sampleStockTyre("UPD-001"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	before, err := store.GetStockTyreByID(ctx, id)
	if err != nil {
		t.Fatalf("get before update: %v", err)
	}

	updated := sampleStockTyre("UPD-001")
	updated.ID = id
	updated.Brand = strPtr("Continental")
	updated.Model = strPtr("PremiumContact 6")

	if err := store.UpdateStockTyre(ctx, updated); err != nil {
		t.Fatalf("UpdateStockTyre: %v", err)
	}

	after, err := store.GetStockTyreByID(ctx, id)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if after.Brand == nil || *after.Brand != "Continental" {
		t.Fatalf("brand = %v, want Continental", after.Brand)
	}
	if !after.UpdatedAt.After(before.UpdatedAt) && !after.UpdatedAt.Equal(before.UpdatedAt) {
		t.Fatalf("updated_at did not advance: before=%v after=%v", before.UpdatedAt, after.UpdatedAt)
	}
}

func TestUpdateStockTyreNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	tyre := sampleStockTyre("GHOST")
	tyre.ID = 404

	err := store.UpdateStockTyre(ctx, tyre)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteStockTyre(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	id, err := store.CreateStockTyre(ctx, sampleStockTyre("DEL-001"))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.DeleteStockTyre(ctx, id); err != nil {
		t.Fatalf("DeleteStockTyre: %v", err)
	}

	_, err = store.GetStockTyreByID(ctx, id)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteStockTyreNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	err := store.DeleteStockTyre(ctx, 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}
