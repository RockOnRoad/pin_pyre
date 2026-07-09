package api

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/RockOnRoad/pin_pyre/db_api/db"
)

type Server struct {
	store      *db.Store
	mux        *http.ServeMux
	addr       string
	allowedIPs map[string]struct{}
}

func NewServer(store *db.Store) *Server {
	s := &Server{
		store:      store,
		mux:        http.NewServeMux(),
		addr:       ":8080",
		allowedIPs: make(map[string]struct{}),
	}
	if port := os.Getenv("PORT"); port != "" {
		s.addr = ":" + port
	}

	// Example:
	// ALLOWED_IPS=192.168.1.100,10.0.0.5
	if ips := os.Getenv("ALLOWED_IPS"); ips != "" {
		for _, ip := range strings.Split(ips, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				s.allowedIPs[ip] = struct{}{}
			}
		}
	}

	s.mux.HandleFunc("GET /health", s.handleHealth)

	// Stock tyres
	// GET /stock-tyres?limit=50&offset=0 — list tyres (defaults: limit=50, offset=0); 200 → []StockTyre
	s.mux.HandleFunc("GET /stock-tyres", s.handleListStockTyres)
	// POST /stock-tyres — body: {"own_code":"...", "brand":"...", ...}; required: own_code; 201 → StockTyre, 409 on duplicate
	s.mux.HandleFunc("POST /stock-tyres", s.handleCreateStockTyre)
	// GET /stock-tyres/own-code/{own_code} — lookup by own_code; 200 → StockTyre, 404 if missing
	s.mux.HandleFunc("GET /stock-tyres/own-code/{own_code}", s.handleGetStockTyreByOwnCode)
	// GET /stock-tyres/{id} — lookup by numeric id; 200 → StockTyre, 404 if missing
	s.mux.HandleFunc("GET /stock-tyres/{id}", s.handleGetStockTyre)
	// PUT /stock-tyres/{id} — body: full StockTyre fields (id from path); required: own_code; 200 → StockTyre, 404/409
	s.mux.HandleFunc("PUT /stock-tyres/{id}", s.handleUpdateStockTyre)
	// DELETE /stock-tyres/{id} — 204 on success, 404 if missing
	s.mux.HandleFunc("DELETE /stock-tyres/{id}", s.handleDeleteStockTyre)

	// Tyre codes
	// POST /tyre-codes — body: {"supplier":"...", "supp_code":"...", "own_code":"..."}; 201 → TyreCode, 409 on duplicate
	s.mux.HandleFunc("POST /tyre-codes", s.handleCreateTyreCode)
	// GET /tyre-codes/by-own-code/{own_code} — all codes linked to a stock tyre; 200 → []TyreCode
	s.mux.HandleFunc("GET /tyre-codes/by-own-code/{own_code}", s.handleListTyreCodesByOwnCode)
	// GET /tyre-codes/by-supplier/{supplier}/supp-code/{supp_code} — lookup by supplier pair; 200 → TyreCode, 404 if missing
	s.mux.HandleFunc("GET /tyre-codes/by-supplier/{supplier}/supp-code/{supp_code}", s.handleGetTyreCodeBySupplierAndSuppCode)
	// GET /tyre-codes/{id} — lookup by numeric id; 200 → TyreCode, 404 if missing
	s.mux.HandleFunc("GET /tyre-codes/{id}", s.handleGetTyreCode)
	// PUT /tyre-codes/{id} — body: supplier, supp_code, own_code (id from path); 200 → TyreCode, 404/409
	s.mux.HandleFunc("PUT /tyre-codes/{id}", s.handleUpdateTyreCode)
	// DELETE /tyre-codes/{id} — 204 on success, 404 if missing
	s.mux.HandleFunc("DELETE /tyre-codes/{id}", s.handleDeleteTyreCode)

	// Supplier 4tochki
	// GET /supplier-4tochki?limit=50&offset=0 — list records (defaults: limit=50, offset=0); 200 → []Supplier4Tochki
	s.mux.HandleFunc("GET /supplier-4tochki", s.handleListSupplier4Tochki)
	// POST /supplier-4tochki — body: Supplier4Tochki fields; 201 → Supplier4Tochki, 409 on duplicate canonical_code
	s.mux.HandleFunc("POST /supplier-4tochki", s.handleCreateSupplier4Tochki)
	// GET /supplier-4tochki/by-canonical-code/{canonical_code} — lookup by canonical_code; 200 → Supplier4Tochki, 404 if missing
	s.mux.HandleFunc("GET /supplier-4tochki/by-canonical-code/{canonical_code}", s.handleGetSupplier4TochkiByCanonicalCode)
	// GET /supplier-4tochki/{id} — lookup by numeric id; 200 → Supplier4Tochki, 404 if missing
	s.mux.HandleFunc("GET /supplier-4tochki/{id}", s.handleGetSupplier4Tochki)
	// PUT /supplier-4tochki/{id} — body: full Supplier4Tochki fields (id from path); 200 → Supplier4Tochki, 404/409
	s.mux.HandleFunc("PUT /supplier-4tochki/{id}", s.handleUpdateSupplier4Tochki)
	// DELETE /supplier-4tochki/{id} — 204 on success, 404 if missing
	s.mux.HandleFunc("DELETE /supplier-4tochki/{id}", s.handleDeleteSupplier4Tochki)

	return s
}

func (s *Server) Run() error {
	fmt.Printf("listening on %s\n", s.addr)
	handler := s.ipWhitelistMiddleware(s.mux)
	return http.ListenAndServe(s.addr, handler)
}

func isAllowedIP(ipStr string, allowedIPs map[string]struct{}) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Always allow localhost and private networks.
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}

	// Explicit whitelist.
	_, ok := allowedIPs[ipStr]
	return ok
}

func (s *Server) ipWhitelistMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if !isAllowedIP(ip, s.allowedIPs) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
