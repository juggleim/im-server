package msgdefines

/*
	Flags:
	1: is cmd msg
	2: is count msg
	4: is state msg
	8: is store msg
	16: is modified msg
	32: is merged msg
	64: is undisturb msg
	128: is no affect sender's conversation msg
	256: is ext msg
	512: is reaction msg
	1024: is stream msg
	2048: is encrypted msg
*/

func IsCmdMsg(flag int32) bool {
	return (flag & 0x1) == 1
}

func SetCmdMsg(flag int32) int32 {
	return flag | 0x1
}

func IsCountMsg(flag int32) bool {
	return (flag & 0x2) == 2
}

func SetCountMsg(flag int32) int32 {
	return flag | 0x2
}

func IsStateMsg(flag int32) bool {
	return (flag & 0x4) == 4
}

func SetStateMsg(flag int32) int32 {
	return flag | 0x4
}

func IsStoreMsg(flag int32) bool {
	return (flag & 0x8) == 8
}

func SetStoreMsg(flag int32) int32 {
	return flag | 0x8
}

func IsModifiedMsg(flag int32) bool {
	return (flag & 0x10) == 16
}

func SetModifiedMsg(flag int32) int32 {
	return flag | 0x10
}

func IsMergedMsg(flag int32) bool {
	return (flag & 0x20) == 32
}

func SetMergedMsg(flag int32) int32 {
	return flag | 0x20
}

func IsUndisturbMsg(flag int32) bool {
	return (flag & 0x40) == 64
}

func SetUndisturbMsg(flag int32) int32 {
	return flag | 0x40
}

func IsNoAffectSenderConver(flag int32) bool {
	return (flag & 0x80) == 128
}

func SetNoAffectSenderConver(flag int32) int32 {
	return flag | 0x80
}

func IsExtMsg(flag int32) bool {
	return (flag & 0x100) == 256
}

func SetExtMsg(flag int32) int32 {
	return flag | 0x100
}

func IsReactionMsg(flag int32) bool {
	return (flag & 0x200) == 512
}

func SetReactionMsg(flag int32) int32 {
	return flag | 0x200
}

func IsStreamMsg(flag int32) bool {
	return (flag & 0x400) == 1024
}

func SetStreamMsg(flag int32) int32 {
	return flag | 0x400
}

func IsEncryptedMsg(flag int32) bool {
	return (flag & 0x800) == 2048
}

func SetEncryptedMsg(flag int32) int32 {
	return flag | 0x800
}
