package redisservice

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *RedisService) ClosePoll(r *wemeet.ClosePollReq) error {
	key := pollsKey + r.RoomId

	err := s.rc.Watch(s.ctx, func(tx *redis.Tx) error {
		g := tx.HGet(s.ctx, key, r.PollId)

		result, err := g.Result()
		if err != nil {
			return err
		}
		if result == "" {
			return errors.New("not found")
		}

		info := new(wemeet.PollInfo)
		err = protojson.Unmarshal([]byte(result), info)
		if err != nil {
			return err
		}

		info.IsRunning = false
		info.ClosedBy = r.UserId
		marshal, err := protojson.Marshal(info)
		if err != nil {
			return err
		}

		pollVal := map[string]string{
			r.PollId: string(marshal),
		}
		tx.HSet(s.ctx, key, pollVal)

		return nil
	}, key)

	return err
}

func (s *RedisService) CleanUpPolls(roomId string, pollIds []string) error {
	pp := s.rc.Pipeline()

	for _, id := range pollIds {
		key := fmt.Sprintf("%s%s:respondents:%s", pollsKey, roomId, id)
		pp.Del(s.ctx, key)
	}

	roomKey := pollsKey + roomId
	pp.Del(s.ctx, roomKey)

	_, err := pp.Exec(s.ctx)
	if err != nil {
		return err
	}

	return nil
}
