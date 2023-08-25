package matchingservice

import (
	"context"
	"fmt"
	"gameapp/contract/broker"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/pkg/timestamp"
	funk "github.com/thoas/go-funk"
	"sync"
	"time"
)

type Service struct {
	repo           Repository
	config         Config
	presenceClient PresenceClient
	pub            broker.Publisher
}

type Repository interface {
	AddUserToWaitingList(userID uint, category entity.Category) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	RemoveUsersFromWaitingList(category entity.Category, userIDs []uint)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request dto.GetPresenceRequest) (dto.GetPresenceResponse, error)
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

func New(config Config, repo Repository, presenceClient PresenceClient, pub broker.Publisher) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient, pub: pub}
}

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"
	err := s.repo.AddUserToWaitingList(req.UserID, req.Category)
	if err != nil {
		return dto.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errmsg.SomethingWentWrong)
	}
	return dto.AddToWaitingListResponse{
		Timeout: s.config.WaitingTimeout,
	}, nil
}

func (s Service) MatchWaitedUsers(ctx context.Context, req dto.MatchWaitedUsersRequest) (dto.MatchWaitedUsersResponse, error) {
	const op = "matchingservice.MatchWaitedUsers"
	// get all user for category
	// get all presence user for category
	// create a final list with merge two last step
	var wg = sync.WaitGroup{}

	for _, category := range entity.GetCategoryList() {

		wg.Add(1)
		go s.Match(ctx, category, &wg)

	}
	wg.Wait()
	return dto.MatchWaitedUsersResponse{}, nil
}

func (s Service) Match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	const op = "matchingservice.Match"
	defer wg.Done()

	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return
	}
	var userIDs = make([]uint, 0)

	for _, member := range list {
		userIDs = append(userIDs, member.UserID)
	}

	if len(userIDs) < 2 {
		fmt.Println("less than two users.")
		return
	}

	presenceList, err := s.presenceClient.GetPresence(ctx, dto.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		fmt.Println("eee  s.presenceClient.GetPresenc", err)
		return
	}

	presenceUserIDs := make([]uint, len(list))
	for _, l := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}

	// TODO - merge presenceList with list based on userID
	// also consider the presence timestamp of each user
	// and remove users from waiting list if the user's timestamp is older than time.Now(-20 seconds)
	//if t < timestamp.Add(-20*time.Second) {
	//	// remove from list
	//}

	var finalList = make([]entity.WaitingMember, 0)
	var toBeRemovedUser = make([]uint, 0)
	for _, member := range list {
		lastOnlineTime, ok := s.getPresenceItem(presenceList, member.UserID)

		if ok && funk.Contains(presenceUserIDs, member.UserID) && lastOnlineTime > timestamp.Add(-20*time.Second) &&
			member.Timestamp > timestamp.Add(-300*time.Second) {
			finalList = append(finalList, member)
		} else {
			// remove from waiting list
			toBeRemovedUser = append(toBeRemovedUser, member.UserID)
		}
	}
	go s.repo.RemoveUsersFromWaitingList(category, toBeRemovedUser)
	matchedUsersToBeRemoved := make([]uint, 0)
	for i := 0; i < len(finalList)-1; i += 2 {
		mu := entity.MatchedUsers{
			Category: category,
			UserID:   []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		fmt.Println("mu", mu)
		fmt.Println("here777777")
		// publish a new event for mu
		go s.pub.Publish(entity.MatchingUsersMatchedEvent, "hello2")
		// remove mu users from waiting list
		matchedUsersToBeRemoved = append(matchedUsersToBeRemoved, mu.UserID...)
	}
	go s.repo.RemoveUsersFromWaitingList(category, matchedUsersToBeRemoved)
}

func (s Service) getPresenceItem(presenceList dto.GetPresenceResponse, userID uint) (int64, bool) {
	for _, item := range presenceList.Items {
		if item.UserID == userID {
			return item.Timestamp, true
		}
	}
	return 0, false
}
