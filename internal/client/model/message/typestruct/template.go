package typestruct

type Template struct {
	TemId    string   `json:"temId"`
	Operator string   `json:"operator"`
	Target   []string `json:"target"`
}

type ModGroupNameTemplate struct {
	*Template
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

type RedEnvelopeFinishedTemplate struct {
	*Template
	Duration int64 `json:"duration"`
	LuckUId  int64 `json:"luckUid"`
}

type MeetingVideoCallEndSubTemplate struct {
	*Template
	Duration int64 `json:"duration"`
}

type OpenRedEnvelopTemplate struct {
	*Template
	RecordId   string `json:"recordId"`
	OutTradeNo string `json:"outTradeNo"`
}

type TempIdEnum string

const (
	TempIDAddAdministrator       TempIdEnum = "add-administrator"
	TempIDRemoveAdministrator               = "remove-administrator"
	TempIDTransferOwner                     = "transfer-owner"
	TempIDAddFriend                         = "add-friend"
	TempIDApplyAddFriend                    = "apply-add-friend"
	TempIDAgreeAddFriend                    = "agree-add-friend"
	TempIDRevokeMsg                         = "revoke-msg"
	TempIDRevokeMsgAll                      = "revoke-msg-all"
	TempIDJoinGroup                         = "join-group"
	TempIDQuitGroup                         = "quit-group"
	TempIDModGroupName                      = "mod-group-name"
	TempIDModGroupNotice                    = "mod-group-notice"
	TempIDModGroupAvatar                    = "mod-group-avatar"
	TempIDInviteJoinGroup                   = "invite-join-group"
	TempIDKickOutGroupMember                = "kick-out-group-member"
	TempIDMomentsComment                    = "moments-comment"
	TempIDMomentsUpvote                     = "moments-upvote"
	TempIDCommentsReply                     = "comments-reply"
	TempIDMomentsRepost                     = "moments-repost"
	TempIDOpenRedEnvelope                   = "open-red-envelope"
	TempIDOpenRedEnvelopeSelf               = "open-red-envelope-self"
	TempIDRedEnvelopeFinished               = "red-envelope-finished"
	TempIDMeetingCreateGroup                = "meeting-create-group"
	TempIDMeetingCreateVoiceCall            = "meeting-create-voice-call"
	TempIDMeetingCreateVideoCall            = "meeting-create-video-call"
	TempIDMeetingVoiceCallEnd               = "meeting-voice-call-end"
	TempIDMeetingVideoCallEnd               = "meeting-video-call-end"
	TempIDMeetingVoiceCallEndSub            = "meeting-voice-call-end-sub"
	TempIDMeetingVideoCallEndSub            = "meeting-video-call-end-sub"
)
