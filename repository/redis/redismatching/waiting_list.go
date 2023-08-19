package redismatching

import (
	"context"
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/richerror"
	"github.com/redis/go-redis/v9"
	"time"
)

const WaitingListPrefix = "waitinglist"

func (d DB) AddUserToWaitingList(userID uint, category entity.Category) error {
	const op = "redismatching.AddToWaitingList"
	_, err := d.adapter.Client().
		ZAdd(context.Background(),
			fmt.Sprintf("%s:%s", WaitingListPrefix, category),
			redis.Z{Score: float64(time.Now().UnixMicro()),
				Member: fmt.Sprintf("%d", userID),
			}).Result()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
