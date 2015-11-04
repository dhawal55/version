package version

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Interface that returns version information.
type Versioner interface {
	GetVersion() string
}

type versionService struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}

func New(v Versioner) *http.ServeMux {
	service := &versionService{Version: v.GetVersion(), Checksum: GetChecksum()}
	return service.registerRoutes()
}

func (v *versionService) registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/version", CorsHandler(v))

	return mux
}

func (v *versionService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(v)
}

func GetChecksum() string {
	file, err := os.Open(os.Args[0])
	if err != nil {
		return "Error getting checksum"
	}
	defer file.Close()

	//expensive, but the hit is once on start and our binaries are small
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "Error getting checksum"
	}

	var result []byte
	return fmt.Sprintf("%x", hash.Sum(result))
}

func CorsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
		}
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
