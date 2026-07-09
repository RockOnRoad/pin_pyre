package db

import (
	"context"
	"database/sql"
	"fmt"
)

const supplier4TochkiColumns = `
	id, canonical_code, wid, hei, dia, siz, lt, seas, stud, brand, model, suv,
	amo_solonkl, amo_kryarsk2, price_per_one_solonkl, price_per_one_kryarsk2,
	created_at, updated_at
`

func (s *Store) CreateSupplier4Tochki(ctx context.Context, r *Supplier4Tochki) (int64, error) {
	const query = `
		INSERT INTO supplier_4tochki (
			canonical_code, wid, hei, dia, siz, lt, seas, stud, brand, model, suv,
			amo_solonkl, amo_kryarsk2, price_per_one_solonkl, price_per_one_kryarsk2
		) VALUES (
			:canonical_code, :wid, :hei, :dia, :siz, :lt, :seas, :stud, :brand, :model, :suv,
			:amo_solonkl, :amo_kryarsk2, :price_per_one_solonkl, :price_per_one_kryarsk2
		)`

	result, err := s.db.NamedExecContext(ctx, query, r)
	if err != nil {
		return 0, fmt.Errorf("create supplier 4tochki record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("supplier 4tochki last insert id: %w", err)
	}

	return id, nil
}

func (s *Store) GetSupplier4TochkiByID(ctx context.Context, id int64) (*Supplier4Tochki, error) {
	query := `SELECT ` + supplier4TochkiColumns + ` FROM supplier_4tochki WHERE id = ?`

	var r Supplier4Tochki
	if err := s.db.GetContext(ctx, &r, query, id); err != nil {
		return nil, fmt.Errorf("get supplier 4tochki by id: %w", err)
	}

	return &r, nil
}

func (s *Store) GetSupplier4TochkiByCanonicalCode(ctx context.Context, canonicalCode string) (*Supplier4Tochki, error) {
	query := `SELECT ` + supplier4TochkiColumns + ` FROM supplier_4tochki WHERE canonical_code = ?`

	var r Supplier4Tochki
	if err := s.db.GetContext(ctx, &r, query, canonicalCode); err != nil {
		return nil, fmt.Errorf("get supplier 4tochki by canonical code: %w", err)
	}

	return &r, nil
}

func (s *Store) ListSupplier4Tochki(ctx context.Context, limit, offset int) ([]Supplier4Tochki, error) {
	query := `SELECT ` + supplier4TochkiColumns + ` FROM supplier_4tochki ORDER BY id LIMIT ? OFFSET ?`

	var records []Supplier4Tochki
	if err := s.db.SelectContext(ctx, &records, query, limit, offset); err != nil {
		return nil, fmt.Errorf("list supplier 4tochki records: %w", err)
	}

	return records, nil
}

func (s *Store) UpdateSupplier4Tochki(ctx context.Context, r *Supplier4Tochki) error {
	const query = `
		UPDATE supplier_4tochki SET
			canonical_code = :canonical_code,
			wid = :wid,
			hei = :hei,
			dia = :dia,
			siz = :siz,
			lt = :lt,
			seas = :seas,
			stud = :stud,
			brand = :brand,
			model = :model,
			suv = :suv,
			amo_solonkl = :amo_solonkl,
			amo_kryarsk2 = :amo_kryarsk2,
			price_per_one_solonkl = :price_per_one_solonkl,
			price_per_one_kryarsk2 = :price_per_one_kryarsk2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, r)
	if err != nil {
		return fmt.Errorf("update supplier 4tochki record: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("supplier 4tochki rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) DeleteSupplier4Tochki(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM supplier_4tochki WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete supplier 4tochki record: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("supplier 4tochki rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
