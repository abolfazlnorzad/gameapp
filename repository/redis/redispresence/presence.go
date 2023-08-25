package redispresence

import (
	"context"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/pkg/timestamp"
	"time"
)

func (d DB) Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error {
	const op = "redispresence.presence"
	_, err := d.adapter.Client().Set(ctx, key, timestamp, expTime).Result()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errmsg.SomethingWentWrong)
	}
	return nil
}

func (d DB) GetPresence(ctx context.Context, prefixKey string, userIDs []uint) (map[uint]int64, error) {
	// TODO - implement me
	// TODO - How to get multiple redis key at once?

	m := make(map[uint]int64)

	for _, u := range userIDs {
		m[u] = timestamp.Add(time.Millisecond * -100)
	}

	return m, nil
}
