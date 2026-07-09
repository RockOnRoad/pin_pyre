package api

import (
	"encoding/json"
	"net/http"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreateStockTyre(w http.ResponseWriter, r *http.Request) {
	var t db.StockTyre
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if t.OwnCode == "" {
		writeError(w, http.StatusBadRequest, "own_code is required")
		return
	}

	id, err := s.store.CreateStockTyre(r.Context(), &t)
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "stock tyre already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create stock tyre")
		return
	}

	created, err := s.store.GetStockTyreByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load created stock tyre")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleCreateTyreCode(w http.ResponseWriter, r *http.Request) {
	var c db.TyreCode
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if c.Supplier == "" || c.SuppCode == "" || c.OwnCode == "" {
		writeError(w, http.StatusBadRequest, "supplier, supp_code, and own_code are required")
		return
	}

	id, err := s.store.CreateTyreCode(r.Context(), &c)
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "tyre code already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create tyre code")
		return
	}

	created, err := s.store.GetTyreCodeByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load created tyre code")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleCreateSupplier4Tochki(w http.ResponseWriter, r *http.Request) {
	var rec db.Supplier4Tochki
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	id, err := s.store.CreateSupplier4Tochki(r.Context(), &rec)
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "supplier 4tochki record already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create supplier 4tochki record")
		return
	}

	created, err := s.store.GetSupplier4TochkiByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load created supplier 4tochki record")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) handleListStockTyres(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)

	tyres, err := s.store.ListStockTyres(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list stock tyres")
		return
	}

	writeJSON(w, http.StatusOK, tyres)
}

func (s *Server) handleGetStockTyre(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	tyre, err := s.store.GetStockTyreByID(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "stock tyre not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get stock tyre")
		return
	}

	writeJSON(w, http.StatusOK, tyre)
}

func (s *Server) handleGetStockTyreByOwnCode(w http.ResponseWriter, r *http.Request) {
	ownCode := r.PathValue("own_code")
	if ownCode == "" {
		writeError(w, http.StatusBadRequest, "own_code is required")
		return
	}

	tyre, err := s.store.GetStockTyreByOwnCode(r.Context(), ownCode)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "stock tyre not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get stock tyre")
		return
	}

	writeJSON(w, http.StatusOK, tyre)
}

func (s *Server) handleUpdateStockTyre(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var t db.StockTyre
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if t.OwnCode == "" {
		writeError(w, http.StatusBadRequest, "own_code is required")
		return
	}

	t.ID = id
	if err := s.store.UpdateStockTyre(r.Context(), &t); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "stock tyre not found")
			return
		}
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "stock tyre already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update stock tyre")
		return
	}

	updated, err := s.store.GetStockTyreByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load updated stock tyre")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteStockTyre(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := s.store.DeleteStockTyre(r.Context(), id); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "stock tyre not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete stock tyre")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetTyreCode(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	code, err := s.store.GetTyreCodeByID(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "tyre code not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get tyre code")
		return
	}

	writeJSON(w, http.StatusOK, code)
}

func (s *Server) handleGetTyreCodeBySupplierAndSuppCode(w http.ResponseWriter, r *http.Request) {
	supplier := r.PathValue("supplier")
	suppCode := r.PathValue("supp_code")
	if supplier == "" || suppCode == "" {
		writeError(w, http.StatusBadRequest, "supplier and supp_code are required")
		return
	}

	code, err := s.store.GetTyreCodeBySupplierAndSuppCode(r.Context(), supplier, suppCode)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "tyre code not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get tyre code")
		return
	}

	writeJSON(w, http.StatusOK, code)
}

func (s *Server) handleListTyreCodesByOwnCode(w http.ResponseWriter, r *http.Request) {
	ownCode := r.PathValue("own_code")
	if ownCode == "" {
		writeError(w, http.StatusBadRequest, "own_code is required")
		return
	}

	codes, err := s.store.ListTyreCodesByOwnCode(r.Context(), ownCode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tyre codes")
		return
	}

	writeJSON(w, http.StatusOK, codes)
}

func (s *Server) handleUpdateTyreCode(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var c db.TyreCode
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if c.Supplier == "" || c.SuppCode == "" || c.OwnCode == "" {
		writeError(w, http.StatusBadRequest, "supplier, supp_code, and own_code are required")
		return
	}

	c.ID = id
	if err := s.store.UpdateTyreCode(r.Context(), &c); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "tyre code not found")
			return
		}
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "tyre code already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update tyre code")
		return
	}

	updated, err := s.store.GetTyreCodeByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load updated tyre code")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteTyreCode(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := s.store.DeleteTyreCode(r.Context(), id); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "tyre code not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete tyre code")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListSupplier4Tochki(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)

	records, err := s.store.ListSupplier4Tochki(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list supplier 4tochki records")
		return
	}

	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleGetSupplier4Tochki(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	rec, err := s.store.GetSupplier4TochkiByID(r.Context(), id)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "supplier 4tochki record not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get supplier 4tochki record")
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleGetSupplier4TochkiByCanonicalCode(w http.ResponseWriter, r *http.Request) {
	canonicalCode := r.PathValue("canonical_code")
	if canonicalCode == "" {
		writeError(w, http.StatusBadRequest, "canonical_code is required")
		return
	}

	rec, err := s.store.GetSupplier4TochkiByCanonicalCode(r.Context(), canonicalCode)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "supplier 4tochki record not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get supplier 4tochki record")
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleUpdateSupplier4Tochki(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var rec db.Supplier4Tochki
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	rec.ID = id
	if err := s.store.UpdateSupplier4Tochki(r.Context(), &rec); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "supplier 4tochki record not found")
			return
		}
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "supplier 4tochki record already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update supplier 4tochki record")
		return
	}

	updated, err := s.store.GetSupplier4TochkiByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load updated supplier 4tochki record")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteSupplier4Tochki(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := s.store.DeleteSupplier4Tochki(r.Context(), id); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusNotFound, "supplier 4tochki record not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete supplier 4tochki record")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
