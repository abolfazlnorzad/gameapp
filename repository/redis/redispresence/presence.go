package redispresence

import (
	"context"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
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
