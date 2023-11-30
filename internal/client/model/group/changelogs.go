package group

import (
	//"imsdk/internal/common/dao/group/members/changelogs"
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
	LogType     int    `json:"type"` // 1: increase, 2: all
	Gid         string `json:"gid"`
	MaxSequence int64  `json:"max_sequence"`
	//Logs        []changelogs.ChangeLogs `json:"logs"`
}

func ChangeLogsList(params []ChangeLogsListParams) ([]ChangeLogsListRes, error) {
	var err error
	changelogsDao := changelogs.New()
	result := make([]ChangeLogsListRes, 0)
	for _, param := range params {
		logType := LogTypeIncrease
		logs := make([]changelogs.ChangeLogs, 0)
		var maxSeq int64 = 0
		logAmount := changelogsDao.GetLogAmount(param.Gid, param.MaxSequence)
		if param.MaxSequence == -1 || logAmount >= LogLimit {
			logType = LogTypeAll
			maxSeq, err = changelogsDao.GetMaxSequence(param.Gid)
			if err != nil {
				return result, errno.Add("sys error", errno.SysErr)
			}
		} else {
			logs, err = changelogsDao.List(param.Gid, param.MaxSequence)
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
			Gid:         param.Gid,
			MaxSequence: maxSeq,
			//Logs:        logs,
		}
		result = append(result, item)
	}
	return result, nil
}
