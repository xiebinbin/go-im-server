package base

const (
	TopicSocketMessagePush  = "imsdk-socket-msg-forward"
	TopicMsgToOffline       = "imsdk-msg-to-offline"
	TopicMsgToWaitReProcess = "imsdk-msg-to-wait-process"
	TopicMsgChunksProcess   = "imsdk-msg-chunks-process"

	PushTypeWebSocket = 1
	PushTypeFcm       = 2
	PushTypeApns      = 3

	QueueTypeBaseMsg  = "base_msg"
	QueueTypeSliceMsg = "slice_msg"
)
