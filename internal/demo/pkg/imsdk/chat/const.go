package chat

type actionType = string

const (
	SignatureVersion                        = 2
	ActionCreateChat             actionType = "createChat"
	ActionJoinChat                          = "joinChat"
	ActionRemoveChatMember                  = "removeMember"
	ActionUpdateChatName                    = "updateChatName"
	ActionUpdateChatAvatar                  = "updateChatAvatar"
	ActionChangeChatMemberStatus            = "changeChatMemberStatus"
)
