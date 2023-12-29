package friend

type GetRelationInfoResponse struct {
	ID       string `json:"uid"`
	IsFriend int8   `json:"is_friend"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	Sign     string `json:"sign"`
	Desc     string `json:"desc"`
}

type DeleteFriendsRequest struct {
	UIDs []string `json:"uids" binding:"required"`
}

type ListFriendsRequest struct {
	UIds []string `json:"uids" binding:"required"`
}

type UpdateRemarkRequest struct {
	ObjUId string `json:"obj_uid" binding:"required"`
	Remark string `json:"remark"`
}
