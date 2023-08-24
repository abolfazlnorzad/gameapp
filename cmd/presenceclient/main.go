package main

import (
	"context"
	"fmt"
	"gameapp/contract/golang/presence"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cl, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	client := presence.NewPresenceServiceClient(cl)
	res, err := client.GetPresence(context.Background(), &presence.GetPresenceRequest{
		UserIds: []uint64{1, 3, 7},
	})

	if err != nil {
		panic(err)
	}

	for _, item := range res.Items {
		fmt.Println("item", item.UserId, item.Timestamp)
	}

}
