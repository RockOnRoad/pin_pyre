package db

import (
	"context"
	"database/sql"
	"fmt"
)

const tyreCodeColumns = `id, supplier, supp_code, own_code, created_at, updated_at`

func (s *Store) CreateTyreCode(ctx context.Context, c *TyreCode) (int64, error) {
	const query = `
		INSERT INTO tyre_codes (supplier, supp_code, own_code)
		VALUES (:supplier, :supp_code, :own_code)`

	result, err := s.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return 0, fmt.Errorf("create tyre code: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("tyre code last insert id: %w", err)
	}

	return id, nil
}

func (s *Store) GetTyreCodeByID(ctx context.Context, id int64) (*TyreCode, error) {
	query := `SELECT ` + tyreCodeColumns + ` FROM tyre_codes WHERE id = ?`

	var c TyreCode
	if err := s.db.GetContext(ctx, &c, query, id); err != nil {
		return nil, fmt.Errorf("get tyre code by id: %w", err)
	}

	return &c, nil
}

func (s *Store) GetTyreCodeBySupplierAndSuppCode(ctx context.Context, supplier, suppCode string) (*TyreCode, error) {
	query := `SELECT ` + tyreCodeColumns + ` FROM tyre_codes WHERE supplier = ? AND supp_code = ?`

	var c TyreCode
	if err := s.db.GetContext(ctx, &c, query, supplier, suppCode); err != nil {
		return nil, fmt.Errorf("get tyre code by supplier and supp code: %w", err)
	}

	return &c, nil
}

func (s *Store) ListTyreCodesByOwnCode(ctx context.Context, ownCode string) ([]TyreCode, error) {
	query := `SELECT ` + tyreCodeColumns + ` FROM tyre_codes WHERE own_code = ? ORDER BY id`

	var codes []TyreCode
	if err := s.db.SelectContext(ctx, &codes, query, ownCode); err != nil {
		return nil, fmt.Errorf("list tyre codes by own code: %w", err)
	}

	return codes, nil
}

func (s *Store) UpdateTyreCode(ctx context.Context, c *TyreCode) error {
	const query = `
		UPDATE tyre_codes SET
			supplier = :supplier,
			supp_code = :supp_code,
			own_code = :own_code,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("update tyre code: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("tyre code rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) DeleteTyreCode(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tyre_codes WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete tyre code: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("tyre code rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
