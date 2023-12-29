package user

type GetConnectInfoResponse struct {
	UId      string `json:"uid"`
	AUId     string `json:"auid"`
	IsOnline int8   `json:"is_online"`
}

type GetUserInfoRequest struct {
	Ids []string `json:"uids" binding:"required"`
}

type UpdateNameRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

type UpdateGenderRequest struct {
	Gender int8 `json:"gender" binding:"required"`
}

type UpdateSignRequest struct {
	Sign string `json:"sign" binding:"required"`
}
