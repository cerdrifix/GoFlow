package pages

import (
	"GoFlow/engine"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

const message = "Welcome to GoFlow.\n"

type Handler struct {
	logger *log.Logger
	db     *sqlx.DB
}

func New(logger *log.Logger, db *sqlx.DB) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

func (h *Handler) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		defer h.logger.Printf("Request processed in %v", time.Now().Sub(startTime).Milliseconds())
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	//_, err := h.db.ExecContext(r.Context(), fmt.Sprintf("call sp_access_logs_insert('%s')", "home"))
	//if err != nil {
	//	h.logger.Fatalf("Error calling sp_access_logs_insert: %v", err)
	//}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(message + " " + time.Now().Format("2006-01-02 15:04:05.999999999")))
	if err != nil {
		h.logger.Fatalf("Error writing the response: %v", err)
	}
}

func (h *Handler) GetMapDetails(w http.ResponseWriter, r *http.Request) {
	var processMap engine.ProcessMap

	const q = "SELECT * FROM public.fn_maps_getlatestbyname('richiesta_con_approvazione')"
	//h.db.GetContext(r.Context(), &processMap, q)
	rows, err := h.db.QueryxContext(r.Context(), q)
	if err != nil {
		h.logger.Fatalf("Unable to query db")
	}
	for rows.Next() {
		var name string
		var version int
		var data string

		err = rows.Scan(&name, &version, &data)
		if err != nil {
			h.logger.Fatalf("Error during row scan")
		}

		h.logger.Printf("\nMap name: %s\n    version: %d\n    data: %s", name, version, data)

		err = json.Unmarshal([]byte(data), &processMap)
		if err != nil {
			h.logger.Fatalf("Error during unmarshaling \njson: %s\nerror: %v", data, err)
		}
	}

	h.logger.Printf("Json parsed: %#v", processMap)

}
