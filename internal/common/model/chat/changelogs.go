package chat

import (
	"context"
	"imsdk/internal/common/dao/chat/changelogs"
	"imsdk/pkg/errno"
)

type ChangeLogsListParams struct {
	Gid         string `json:"gid"`
	MaxSequence int64  `json:"max_sequence"`
}

const (
	LogTypeIncrease = 1
	LogTypeAll      = 2
	LogLimit        = 200
)

type ChangeLogsListRes struct {
	LogType     int                     `json:"type"` // 1: increase, 2: all
	ChatId      string                  `json:"chat_id"`
	MaxSequence int64                   `json:"max_sequence"`
	Logs        []changelogs.ChangeLogs `json:"logs"`
}

type ChangeLogsReq struct {
	Chat []ChangeLogsReqItem `json:"chat" binding:"required"`
}

type ChangeLogsReqItem struct {
	ChatID      string `json:"chat_id"`
	MaxSequence int64  `json:"max_sequence"`
}

type ChangeLogsResp struct {
	ChatID      string                 `json:"chat_id"`
	Type        int8                   `json:"type"`
	MaxSequence int64                  `json:"max_sequence"`
	Logs        map[string]interface{} `json:"logs"`
}

func ChangeLogsList(ctx context.Context, params ChangeLogsReq) ([]ChangeLogsListRes, error) {
	var err error
	changelogsDao := changelogs.New()
	result := make([]ChangeLogsListRes, 0)
	for _, param := range params.Chat {
		logType := LogTypeIncrease
		logs := make([]changelogs.ChangeLogs, 0)
		var maxSeq int64 = 0
		logAmount := changelogsDao.GetLogAmount(param.ChatID, param.MaxSequence)
		if param.MaxSequence == -1 || logAmount >= LogLimit {
			logType = LogTypeAll
			maxSeq, err = changelogsDao.GetMaxSequence(param.ChatID)
			if err != nil {
				return result, errno.Add("sys error", errno.SysErr)
			}
		} else {
			logs, err = changelogsDao.List(param.ChatID, param.MaxSequence)
			if err != nil {
				return result, errno.Add("sys error", errno.SysErr)
			}
			maxSeq = param.MaxSequence
			if len(logs) > 0 {
				maxSeq = logs[0].Sequence
			}
			if param.MaxSequence == 0 && len(logs) == 0 {
				continue
			}
		}
		item := ChangeLogsListRes{
			LogType:     logType,
			ChatId:      param.ChatID,
			MaxSequence: maxSeq,
			Logs:        logs,
		}
		result = append(result, item)
	}
	return result, nil
}
