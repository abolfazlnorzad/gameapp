package protobuf

import (
	"gameapp/contract/golang/presence"
	"gameapp/dto"
)

func MapGetPresenceResponseToProtobuf(g dto.GetPresenceResponse) *presence.GetPresenceResponse {
	r := &presence.GetPresenceResponse{}

	for _, item := range g.Items {
		r.Items = append(r.Items, &presence.GetPresenceItem{
			UserId:    uint64(item.UserID),
			Timestamp: item.Timestamp,
		})
	}

	return r
}

func MapGetPresenceResponseFromProtobuf(g *presence.GetPresenceResponse) dto.GetPresenceResponse {
	r := dto.GetPresenceResponse{}

	for _, item := range g.Items {
		r.Items = append(r.Items, dto.GetPresenceItem{
			UserID:    uint(item.UserId),
			Timestamp: item.Timestamp,
		})
	}

	return r
}
