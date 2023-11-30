package sdk

type StatusType = int8
type JoinChatType = uint8
type JoinGroupType = uint8

const (
	EmptyString                       = ""
	StatusNormal        StatusType    = 1
	StatusForbidden     StatusType    = -1
	JoinGroupTypeSelf   JoinGroupType = 1
	JoinGroupTypeInvite JoinGroupType = 2
	JoinChatTypeSelf    JoinChatType  = 1
	JoinChatTypeInvite  JoinChatType  = 2
)
