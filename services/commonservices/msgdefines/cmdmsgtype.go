package msgdefines

var (
	CmdMsgType_ClearUnread      string = "jg:clearunread"
	CmdMsgType_MarkUnread       string = "jg:markunread"
	CmdMsgType_LogCmd           string = "jg:logcmd"
	CmdMsgType_AddConver        string = "jg:addconver"
	CmdMsgType_DelConvers       string = "jg:delconvers"
	CmdMsgType_TopConvers       string = "jg:topconvers"
	CmdMsgType_Undisturb        string = "jg:undisturb"
	CmdMsgType_ClearTotalUnread string = "jg:cleartotalunread"
	CmdMsgType_TagAddConvers    string = "jg:tagaddconvers"
	CmdMsgType_TagDelConvers    string = "jg:tagdelconvers"
	CmdMsgType_DelConverTags    string = "jg:delconvertags"
	CmdMsgType_CleanMsg         string = "jg:cleanmsg"
	CmdMsgType_DelMsgs          string = "jg:delmsgs"
	CmdMsgType_ReadNtf          string = "jg:readntf"
	CmdMsgType_MsgModify        string = "jg:modify"
	CmdMsgType_MsgExt           string = "jg:msgext"
	CmdMsgType_MsgExSet         string = "jg:msgexset"
	CmdMsgType_GrpReadNtf       string = "jg:grpreadntf"
	CmdMsgType_Recall           string = "jg:recall"
	CmdMsgType_RecallInfo       string = "jg:recallinfo"
	CmdMsgType_TopMsg           string = "jg:topmsg"
)

var cmdMsgMap map[string]bool

func init() {
	cmdMsgMap = make(map[string]bool)
	cmdMsgMap[CmdMsgType_ClearUnread] = true
	cmdMsgMap[CmdMsgType_MarkUnread] = true
	cmdMsgMap[CmdMsgType_LogCmd] = true
	cmdMsgMap[CmdMsgType_AddConver] = true
	cmdMsgMap[CmdMsgType_DelConvers] = true
	cmdMsgMap[CmdMsgType_TopConvers] = true
	cmdMsgMap[CmdMsgType_Undisturb] = true
	cmdMsgMap[CmdMsgType_ClearTotalUnread] = true
	cmdMsgMap[CmdMsgType_TagAddConvers] = true
	cmdMsgMap[CmdMsgType_TagDelConvers] = true
	cmdMsgMap[CmdMsgType_DelConverTags] = true
	cmdMsgMap[CmdMsgType_CleanMsg] = true
	cmdMsgMap[CmdMsgType_DelMsgs] = true
	cmdMsgMap[CmdMsgType_ReadNtf] = true
	cmdMsgMap[CmdMsgType_MsgModify] = true
	cmdMsgMap[CmdMsgType_MsgExt] = true
	cmdMsgMap[CmdMsgType_MsgExSet] = true
	cmdMsgMap[CmdMsgType_GrpReadNtf] = true
	cmdMsgMap[CmdMsgType_Recall] = true
	cmdMsgMap[CmdMsgType_RecallInfo] = true
	cmdMsgMap[CmdMsgType_TopMsg] = true
}

func IsCmdMsgType(msgType string) bool {
	_, ok := cmdMsgMap[msgType]
	return ok
}
