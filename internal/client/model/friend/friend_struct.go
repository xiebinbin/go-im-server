package friend

type GetRelationInfoResponse struct {
	ID       string `json:"uid"`
	IsFriend int8   `json:"is_friend"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

type DeleteFriendsRequest struct {
	UIDs []string `json:"uids" binding:"required"`
}
