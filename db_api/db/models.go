package db

import "time"

type StockTyre struct {
	ID        int64     `db:"id" json:"id"`
	OwnCode   string    `db:"own_code" json:"own_code"`
	Wid       *string   `db:"wid" json:"wid,omitempty"`
	Hei       *string   `db:"hei" json:"hei,omitempty"`
	Dia       *int      `db:"dia" json:"dia,omitempty"`
	Idx       *string   `db:"idx" json:"idx,omitempty"`
	XL        *bool     `db:"xl" json:"xl,omitempty"`
	Siz       *string   `db:"siz" json:"siz,omitempty"`
	Lt        *string   `db:"lt" json:"lt,omitempty"`
	Seas      *string   `db:"seas" json:"seas,omitempty"`
	Stud      *bool     `db:"stud" json:"stud,omitempty"`
	Brand     *string   `db:"brand" json:"brand,omitempty"`
	Model     *string   `db:"model" json:"model,omitempty"`
	SUV       *bool     `db:"suv" json:"suv,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type TyreCode struct {
	ID        int64     `db:"id" json:"id"`
	Supplier  string    `db:"supplier" json:"supplier"`
	SuppCode  string    `db:"supp_code" json:"supp_code"`
	OwnCode   string    `db:"own_code" json:"own_code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Supplier4Tochki struct {
	ID                  int64     `db:"id" json:"id"`
	CanonicalCode       *string   `db:"canonical_code" json:"canonical_code,omitempty"`
	Wid                 *string   `db:"wid" json:"wid,omitempty"`
	Hei                 *string   `db:"hei" json:"hei,omitempty"`
	Dia                 *int      `db:"dia" json:"dia,omitempty"`
	Siz                 *string   `db:"siz" json:"siz,omitempty"`
	Lt                  *string   `db:"lt" json:"lt,omitempty"`
	Seas                *string   `db:"seas" json:"seas,omitempty"`
	Stud                *bool     `db:"stud" json:"stud,omitempty"`
	Brand               *string   `db:"brand" json:"brand,omitempty"`
	Model               *string   `db:"model" json:"model,omitempty"`
	SUV                 *bool     `db:"suv" json:"suv,omitempty"`
	AmoSolonkl          *int      `db:"amo_solonkl" json:"amo_solonkl,omitempty"`
	AmoKryarsk2         *int      `db:"amo_kryarsk2" json:"amo_kryarsk2,omitempty"`
	PricePerOneSolonkl  *int      `db:"price_per_one_solonkl" json:"price_per_one_solonkl,omitempty"`
	PricePerOneKryarsk2 *int      `db:"price_per_one_kryarsk2" json:"price_per_one_kryarsk2,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}
