package translate

import (
	"imsdk/internal/client/model/common"
	"imsdk/pkg/app"
)

type PushTagType string

const (
	ImPushVoice        PushTagType = "push_voice"
	ImPushImage        PushTagType = "push_image"
	ImPushAttachment   PushTagType = "push_attachment"
	ImPushAudio        PushTagType = "push_audio"
	ImPushVideo        PushTagType = "push_video"
	ImPushConversation PushTagType = "push_conversation"
	ImPushRemindAll    PushTagType = "push_remind_all"
	MomentsBase        PushTagType = "push_moments"
	ImPushLocation     PushTagType = "push_location"
	ImRevokeMessage    PushTagType = "push_revoke_message"
	PaymentTransfer    PushTagType = "push_transfer"
	PaymentRedPacket   PushTagType = "push_red_packet"

	PaymentPaySuc          PushTagType = "pay_suc"
	PaymentCoinReceive     PushTagType = "coin_receive"
	PaymentWithdrawIng     PushTagType = "withdraw_ing"
	PaymentWithdrawSuc     PushTagType = "withdraw_suc"
	PaymentWithdrawRefund  PushTagType = "withdraw_refund"
	PaymentRedPacketRefund PushTagType = "red_packet_refund"

	MeetingInviteVoice     PushTagType = "meeting_invite_voice"
	MeetingInviteVideo     PushTagType = "meeting_invite_video"
	MeetingInviteManyVoice PushTagType = "meeting_invite_many_voice"
	MeetingInviteManyVideo PushTagType = "meeting_invite_many_video"
	MeetingWaitingCall     PushTagType = "meeting_waiting_call"
	MeetingOtherBusy       PushTagType = "meeting_other_busy"
	MeetingNotAnswered     PushTagType = "meeting_not_answered"
	MeetingDeclined        PushTagType = "meeting_declined"
	MeetingOtherDeclined   PushTagType = "meeting_other_declined"
	MeetingCallerCancelled PushTagType = "meeting_caller_cancelled"
	MeetingCallFailed      PushTagType = "meeting_call_failed"
	MeetingOtherCancel     PushTagType = "meeting_other_cancel"
	MeetingChangeVoice     PushTagType = "meeting_change_voice"
	MeetingChangeVideo     PushTagType = "meeting_change_video"
	MeetingCallDuration    PushTagType = "meeting_call_duration"
)

func GetTransByTag(pushTag PushTagType, targetLang string) string {
	defLang, _ := app.Config().GetChildConf("global", "default", "lang")
	resStr := common.GetTransByTag(string(pushTag), targetLang)
	if resStr == "" {
		common.GetTransByTag(string(pushTag), defLang.(string))
	}
	return resStr
}
