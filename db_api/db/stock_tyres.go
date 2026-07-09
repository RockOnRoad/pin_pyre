package db

import (
	"context"
	"database/sql"
	"fmt"
)

const stockTyreColumns = `
	id, own_code, wid, hei, dia, idx, xl, siz, lt, seas, stud, brand, model, suv,
	created_at, updated_at
`

func (s *Store) CreateStockTyre(ctx context.Context, t *StockTyre) (int64, error) {
	const query = `
		INSERT INTO stock_tyres (
			own_code, wid, hei, dia, idx, xl, siz, lt, seas, stud, brand, model, suv
		) VALUES (
			:own_code, :wid, :hei, :dia, :idx, :xl, :siz, :lt, :seas, :stud, :brand, :model, :suv
		)`

	result, err := s.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return 0, fmt.Errorf("create stock tyre: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("stock tyre last insert id: %w", err)
	}

	return id, nil
}

func (s *Store) GetStockTyreByID(ctx context.Context, id int64) (*StockTyre, error) {
	query := `SELECT ` + stockTyreColumns + ` FROM stock_tyres WHERE id = ?`

	var t StockTyre
	if err := s.db.GetContext(ctx, &t, query, id); err != nil {
		return nil, fmt.Errorf("get stock tyre by id: %w", err)
	}

	return &t, nil
}

func (s *Store) GetStockTyreByOwnCode(ctx context.Context, ownCode string) (*StockTyre, error) {
	query := `SELECT ` + stockTyreColumns + ` FROM stock_tyres WHERE own_code = ?`

	var t StockTyre
	if err := s.db.GetContext(ctx, &t, query, ownCode); err != nil {
		return nil, fmt.Errorf("get stock tyre by own code: %w", err)
	}

	return &t, nil
}

func (s *Store) ListStockTyres(ctx context.Context, limit, offset int) ([]StockTyre, error) {
	query := `SELECT ` + stockTyreColumns + ` FROM stock_tyres ORDER BY id LIMIT ? OFFSET ?`

	var tyres []StockTyre
	if err := s.db.SelectContext(ctx, &tyres, query, limit, offset); err != nil {
		return nil, fmt.Errorf("list stock tyres: %w", err)
	}

	return tyres, nil
}

func (s *Store) UpdateStockTyre(ctx context.Context, t *StockTyre) error {
	const query = `
		UPDATE stock_tyres SET
			own_code = :own_code,
			wid = :wid,
			hei = :hei,
			dia = :dia,
			idx = :idx,
			xl = :xl,
			siz = :siz,
			lt = :lt,
			seas = :seas,
			stud = :stud,
			brand = :brand,
			model = :model,
			suv = :suv,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("update stock tyre: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("stock tyre rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) DeleteStockTyre(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM stock_tyres WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete stock tyre: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("stock tyre rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
