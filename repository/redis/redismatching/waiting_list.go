package redismatching

import (
	"context"
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/richerror"
	"gameapp/pkg/timestamp"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"strconv"
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

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = "redismatching.GetWaitingListByCategory"
	// TODO - add to config
	minimum := fmt.Sprintf("%d", timestamp.Add(-200000*time.Hour))
	max := strconv.Itoa(int(timestamp.Now()))

	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx, getCategoryKey(category), &redis.ZRangeBy{
		Min:    minimum,
		Max:    max,
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	var result = make([]entity.WaitingMember, 0)

	for _, l := range list {
		userID, _ := strconv.Atoi(l.Member.(string))

		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(l.Score),
			Category:  category,
		})
	}

	return result, nil
}

func getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WaitingListPrefix, category)
}

func (d DB) RemoveUsersFromWaitingList(category entity.Category, userIDs []uint) {
	// TODO - add 5 to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	members := make([]any, 0)
	for _, u := range userIDs {
		members = append(members, strconv.Itoa(int(u)))
	}

	numberOfRemovedMembers, err := d.adapter.Client().ZRem(ctx, getCategoryKey(category), members...).Result()
	if err != nil {
		log.Errorf("remove from waiting list: %v", err)
		// TODO - update metrics
	}

	log.Printf("%d items removed from %s", numberOfRemovedMembers, getCategoryKey(category))
	// TODO - update metrics
}
