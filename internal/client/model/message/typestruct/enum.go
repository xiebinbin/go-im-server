package typestruct

type MessageType int8

const (
	TypeText         MessageType = 1
	TypeImage                    = 2
	TypeVoice                    = 3
	TypeVideo                    = 4
	TypeAttachment               = 5
	TypeRedPacket                = 6
	TypeTransfer                 = 7
	TypeCalls                    = 8
	TypeDelApplet                = 9
	TypeMoments                  = 10
	TypePayment                  = 11
	TypeRedEnvelope              = 12
	TypeLocation                 = 13
	TypeMeeting                  = 14
	TypeRemind                   = 15
	TypeDelChat                  = 16
	TypeCmd                      = 18
	TypeCard                     = 19
	TypeCmdA                     = 20
	TypeNotice                   = 21
	TypeVerticalCard             = 23
	TypeDelMessage               = 81
	TypeReadReport               = 82
)
