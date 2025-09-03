package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"google.golang.org/protobuf/encoding/protojson"
)

func (m *PollModel) ListPolls(roomId string) ([]*wemeet.PollInfo, error) {
	var polls []*wemeet.PollInfo

	result, err := m.rs.GetPollsListByRoomId(roomId)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		// no polls
		return polls, err
	}

	for _, pi := range result {
		info := new(wemeet.PollInfo)
		err = protojson.Unmarshal([]byte(pi), info)
		if err != nil {
			continue
		}

		polls = append(polls, info)
	}

	return polls, nil
}

func (m *PollModel) UserSelectedOption(roomId, pollId, userId string) (uint64, error) {
	allRespondents, err := m.rs.GetPollResponsesByField(roomId, pollId, "all_respondents")
	if err != nil {
		return 0, err
	}

	if allRespondents == "" {
		return 0, err
	}

	var respondents []string
	err = json.Unmarshal([]byte(allRespondents), &respondents)
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(respondents); i++ {
		// format userId:option_id:name
		p := strings.Split(respondents[i], ":")
		if p[0] == userId {
			voted, err := strconv.ParseUint(p[1], 10, 64)
			if err != nil {
				return 0, err
			}
			return voted, err
		}
	}

	return 0, nil
}

func (m *PollModel) GetPollResponsesDetails(roomId, pollId string) (map[string]string, error) {
	result, err := m.rs.GetPollResponsesByPollId(roomId, pollId)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) < 0 {
		return nil, nil
	}

	return result, nil
}

func (m *PollModel) GetResponsesResult(roomId, pollId string) (*wemeet.PollResponsesResult, error) {
	pi, err := m.rs.GetPollInfoByPollId(roomId, pollId)
	if err != nil {
		return nil, err
	}

	info := new(wemeet.PollInfo)
	err = protojson.Unmarshal([]byte(pi), info)
	if err != nil {
		return nil, err
	}
	if info.IsRunning {
		return nil, errors.New("need to wait until poll close")
	}

	res := new(wemeet.PollResponsesResult)
	res.Question = info.Question

	result, err := m.rs.GetPollResponsesByPollId(roomId, pollId)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	var options []*wemeet.PollResponsesResultOptions
	for _, opt := range info.Options {
		f := fmt.Sprintf("%d_count", opt.Id)
		i, _ := strconv.Atoi(result[f])
		rr := &wemeet.PollResponsesResultOptions{
			Id:        uint64(opt.Id),
			Text:      opt.Text,
			VoteCount: uint64(i),
		}
		options = append(options, rr)
	}

	res.Options = options
	i, _ := strconv.Atoi(result["total_resp"])
	res.TotalResponses = uint64(i)

	return res, nil
}

func (m *PollModel) GetPollsStats(roomId string) (*wemeet.PollsStats, error) {
	res := &wemeet.PollsStats{
		TotalPolls:   0,
		TotalRunning: 0,
	}

	result, err := m.rs.GetPollsListByRoomId(roomId)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result) == 0 {
		// no polls
		return nil, nil
	}
	res.TotalPolls = uint64(len(result))

	for _, pi := range result {
		info := new(wemeet.PollInfo)
		err = protojson.Unmarshal([]byte(pi), info)
		if err != nil {
			continue
		}

		if info.IsRunning {
			res.TotalRunning += 1
		}
	}

	return res, nil
}
