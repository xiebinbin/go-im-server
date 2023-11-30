package group

type IdRequest struct {
	GroupID string `json:"id" binding:"required"`
}

type IdsRequest struct {
	GroupIDs []string `json:"ids" binding:"required"`
}
type ListRequest struct {
	GroupIDs []string `json:"ids"`
}

type MembersRequest struct {
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid"`
}

type QuitRequest struct {
	GroupID   string `json:"id" binding:"required"`
	IsDelChat int    `json:"is_del_chat"`
}

type QuitAllRequest struct {
	IsDelChat int `json:"is_del_chat"`
}

type JoinRequest struct {
	GroupID string `json:"id" binding:"required"`
	QrID    string `json:"qr_code"`
}

type AgreeJoinRequest struct {
	UID     string   `json:"uid"`
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid" binding:"required"`
}

type InviteJoinRequest struct {
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid" binding:"required"`
}

type KickOutRequest struct {
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid" binding:"required"`
}

type UpdateNameRequest struct {
	GroupID string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

type UpdateAvatarRequest struct {
	GroupID  string `json:"id" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	IsNotice int8   `json:"is_notice"`
}

type UpdateAliasRequest struct {
	GroupID string `json:"id" binding:"required"`
	ObjUId  string `json:"obj_uid"`
	Alias   string `json:"alias" binding:"required"`
}

type UpdateNoticeRequest struct {
	GroupID string `json:"id" binding:"required"`
	Notice  string `json:"notice"`
}

type AddAdministratorRequest struct {
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid" binding:"required"`
}

type RemoveAdministratorRequest struct {
	GroupID string   `json:"id" binding:"required"`
	ObjUid  []string `json:"obj_uid" binding:"required"`
}

type TransferRequest struct {
	GroupID string `json:"id" binding:"required"`
	ObjUId  string `json:"obj_uid" binding:"required"`
}

type GetNoticeRequest struct {
	GroupID  string `json:"id" binding:"required"`
	NoticeID string `json:"notice_id"`
}

type DisbandRequest struct {
	GroupID string `json:"id" binding:"required"`
}
