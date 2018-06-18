package dbx

// TryLock tries to obtain a lock on the database
func (db *Database) TryLock() (bool, error) {
	var locked bool
	row := db.conn.QueryRow("SELECT pg_try_advisory_lock(123987);")
	err := row.Scan(&locked)
	if err != nil {
		return false, err
	}
	return locked, nil
}

// Unlock releases the lock on the database
func (db *Database) Unlock() error {
	_, err := db.conn.Exec("SELECT pg_advisory_unlock(123987);")
	return err
}
