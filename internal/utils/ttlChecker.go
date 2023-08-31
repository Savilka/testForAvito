package utils

import (
	"fmt"
	"testForAvito/internal/storage/postgres"
	"time"
)

func CheckTtlForSegments(db *postgres.Storage) error {
	const op = "storage.postgres.DeleteExpiredSegmentsFromUsers"

	for true {
		err := db.DeleteExpiredSegmentsFromUsers()
		if err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}

		time.Sleep(time.Hour)
	}

	return nil
}
