package db

func ZeroToNil(n uint64) any {
	// because NULL is good in SQLite for rows that have FK
	// and not have a record to point to
	if n == 0 {
		return nil
	}
	return n
}
