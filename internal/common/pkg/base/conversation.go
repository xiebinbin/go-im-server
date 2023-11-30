package base

const (
	SyncProcessMaxAmount = 5

	MsgContentTypeTxt       = 1
	MsgContentTypeImg       = 2
	MsgContentTypeHollow    = 18
	MsgContentTypeHollowDel = 81
	MsgContentTypeUnread    = 82 // read report

	MsgTypeText         = 1
	MsgTypeImage        = 2
	MsgTypeVoice        = 3
	MsgTypeVideo        = 4
	MsgTypeAttachment   = 5
	MsgTypeRedPacket    = 6
	MsgTypeTransfer     = 7
	MsgTypeCalls        = 8
	MsgTypeDelApplet    = 9
	MsgTypeMoments      = 10
	MsgTypePay          = 11
	MsgTypeRedEnvelope  = 12
	MsgTypeLocation     = 13
	MsgTypeMeeting      = 14
	MsgTypeRemind       = 15
	MsgTypeDelChat      = 16
	MsgTypeCmd          = 18
	MsgTypeCard         = 19
	MsgTypeCmdA         = 20
	MsgTypeNotification = 21
	MsgTypeCustomize    = 22
	TypeVerticalCard    = 23

	ChatTypeSingle = 1
	ChatTypeGroup  = 2
)

var NotAddUnreadCountType = []int{MsgContentTypeHollow, MsgContentTypeHollowDel,
	MsgContentTypeUnread, MsgTypeRedEnvelope, MsgTypeDelChat, MsgTypeCmdA}
var NotVerifyBlockType = []int{MsgContentTypeHollowDel, MsgContentTypeUnread, MsgTypeMeeting}
var AllowTransType = []int{MsgTypeText, MsgTypeRemind}
var AllowGetType = []int8{
	MsgTypeText, MsgTypeImage, MsgTypeVoice, MsgTypeLocation, MsgTypeAttachment,
	MsgTypeCard, MsgTypeNotification, MsgTypeCustomize,
}
