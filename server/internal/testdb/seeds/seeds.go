package seeds

import "database/sql"

// Exec wraps tx.Exec and panics on error. Seeds are test infrastructure
// - if they fail, the test cannot run and should abort immediately.
func Exec(tx *sql.Tx, query string) {
	if _, err := tx.Exec(query); err != nil {
		panic("seed exec failed: " + err.Error())
	}
}
