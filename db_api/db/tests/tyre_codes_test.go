package dbTests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func createStockTyreForTest(t *testing.T, store *db.Store, ownCode string) {
	t.Helper()
	if _, err := store.CreateStockTyre(context.Background(), sampleStockTyre(ownCode)); err != nil {
		t.Fatalf("create stock tyre %s: %v", ownCode, err)
	}
}

func TestCreateTyreCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "MIC-205-55-16")

	code := &db.TyreCode{
		Supplier: "4tochki",
		SuppCode: "ABC123",
		OwnCode:  "MIC-205-55-16",
	}

	id, err := store.CreateTyreCode(ctx, code)
	if err != nil {
		t.Fatalf("CreateTyreCode: %v", err)
	}

	got, err := store.GetTyreCodeByID(ctx, id)
	if err != nil {
		t.Fatalf("GetTyreCodeByID: %v", err)
	}
	if got.Supplier != "4tochki" || got.SuppCode != "ABC123" {
		t.Fatalf("got %+v, want supplier=4tochki supp_code=ABC123", got)
	}
}

func TestCreateTyreCodeDuplicate(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "MIC-205-55-16")

	code := &db.TyreCode{Supplier: "4tochki", SuppCode: "ABC123", OwnCode: "MIC-205-55-16"}
	if _, err := store.CreateTyreCode(ctx, code); err != nil {
		t.Fatalf("first create: %v", err)
	}

	if _, err := store.CreateTyreCode(ctx, code); err == nil {
		t.Fatal("expected error for duplicate supplier+supp_code")
	}
}

func TestCreateTyreCodeForeignKeyViolation(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	code := &db.TyreCode{Supplier: "4tochki", SuppCode: "NOPE", OwnCode: "MISSING"}
	if _, err := store.CreateTyreCode(ctx, code); err == nil {
		t.Fatal("expected foreign key error")
	}
}

func TestGetTyreCodeBySupplierAndSuppCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "MIC-205-55-16")

	if _, err := store.CreateTyreCode(ctx, &db.TyreCode{
		Supplier: "4tochki", SuppCode: "XYZ999", OwnCode: "MIC-205-55-16",
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := store.GetTyreCodeBySupplierAndSuppCode(ctx, "4tochki", "XYZ999")
	if err != nil {
		t.Fatalf("GetTyreCodeBySupplierAndSuppCode: %v", err)
	}
	if got.OwnCode != "MIC-205-55-16" {
		t.Fatalf("own_code = %q, want MIC-205-55-16", got.OwnCode)
	}
}

func TestListTyreCodesByOwnCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "MIC-205-55-16")

	for _, suppCode := range []string{"A1", "A2", "A3"} {
		if _, err := store.CreateTyreCode(ctx, &db.TyreCode{
			Supplier: "4tochki", SuppCode: suppCode, OwnCode: "MIC-205-55-16",
		}); err != nil {
			t.Fatalf("create %s: %v", suppCode, err)
		}
	}

	codes, err := store.ListTyreCodesByOwnCode(ctx, "MIC-205-55-16")
	if err != nil {
		t.Fatalf("ListTyreCodesByOwnCode: %v", err)
	}
	if len(codes) != 3 {
		t.Fatalf("len = %d, want 3", len(codes))
	}
}

func TestUpdateTyreCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "OLD-CODE")
	createStockTyreForTest(t, store, "NEW-CODE")

	id, err := store.CreateTyreCode(ctx, &db.TyreCode{
		Supplier: "4tochki", SuppCode: "UPD1", OwnCode: "OLD-CODE",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	before, err := store.GetTyreCodeByID(ctx, id)
	if err != nil {
		t.Fatalf("get before update: %v", err)
	}

	updated := &db.TyreCode{
		ID:       id,
		Supplier: "other",
		SuppCode: "UPD2",
		OwnCode:  "NEW-CODE",
	}
	if err := store.UpdateTyreCode(ctx, updated); err != nil {
		t.Fatalf("UpdateTyreCode: %v", err)
	}

	after, err := store.GetTyreCodeByID(ctx, id)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if after.Supplier != "other" || after.SuppCode != "UPD2" || after.OwnCode != "NEW-CODE" {
		t.Fatalf("after update = %+v", after)
	}
	if !after.UpdatedAt.After(before.UpdatedAt) && !after.UpdatedAt.Equal(before.UpdatedAt) {
		t.Fatalf("updated_at did not advance")
	}
}

func TestUpdateTyreCodeNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	err := store.UpdateTyreCode(ctx, &db.TyreCode{ID: 999, Supplier: "x", SuppCode: "y", OwnCode: "z"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteTyreCode(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "MIC-205-55-16")

	id, err := store.CreateTyreCode(ctx, &db.TyreCode{
		Supplier: "4tochki", SuppCode: "DEL1", OwnCode: "MIC-205-55-16",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := store.DeleteTyreCode(ctx, id); err != nil {
		t.Fatalf("DeleteTyreCode: %v", err)
	}

	_, err = store.GetTyreCodeByID(ctx, id)
	if err == nil {
		t.Fatal("expected not found after delete")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteTyreCodeNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	err := store.DeleteTyreCode(ctx, 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want sql.ErrNoRows", err)
	}
}

func TestDeleteStockTyreCascadesTyreCodes(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()
	createStockTyreForTest(t, store, "CASCADE-001")

	tyreID, err := store.GetStockTyreByOwnCode(ctx, "CASCADE-001")
	if err != nil {
		t.Fatalf("get stock tyre: %v", err)
	}

	if _, err := store.CreateTyreCode(ctx, &db.TyreCode{
		Supplier: "4tochki", SuppCode: "C1", OwnCode: "CASCADE-001",
	}); err != nil {
		t.Fatalf("create tyre code: %v", err)
	}

	if err := store.DeleteStockTyre(ctx, tyreID.ID); err != nil {
		t.Fatalf("DeleteStockTyre: %v", err)
	}

	codes, err := store.ListTyreCodesByOwnCode(ctx, "CASCADE-001")
	if err != nil {
		t.Fatalf("ListTyreCodesByOwnCode: %v", err)
	}
	if len(codes) != 0 {
		t.Fatalf("expected cascade delete, got %d tyre codes", len(codes))
	}
}
