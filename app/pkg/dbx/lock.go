package dbx

import (
	"context"
	"hash/fnv"

	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Try to obtain an advisory lock
// returns true and an unlock function if lock was aquired
func TryLock(ctx context.Context, trx *Trx, key string) (bool, func()) {
	var locked bool
	if err := trx.Scalar(&locked, "SELECT pg_try_advisory_xact_lock($1)", hash(key)); err != nil {
		log.Error(ctx, errors.Wrap(err, "failed to acquire advisory lock"))
		return false, nil
	}

	unlock := func() {
		trx.MustCommit()
	}

	return locked, unlock
}
