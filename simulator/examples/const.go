package examples

import (
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"
)

var (
	Xl_user1  string = "CgZhcHBrZXkaIPIePJzT2dTjP6OPNuZUoUoMpNWy353LbioWXcjNJdL7"
	Xl_user2  string = "CgZhcHBrZXkaIJw28GY+oWd5spQFnF6RsnrFVMoPjR4zq3kUijlj5ufC"
	Xl_user3  string = "CgZhcHBrZXkaIP/kXX++OmyxFT+NgFZ5ogXKUELGpk93NsDeLghwWt+w"
	Xl_user4  string = "CgZhcHBrZXkaILVSqLtTcmMqif1j1nJeFZDo70V3q4IaTe+0mImbpewt"
	Xl_user5  string = "CgZhcHBrZXkaIIhr2I+s7wC4Kv4tSIfQ8SCaTMJuNfwx9ooWCzLe3qvM"
	Xl_user6  string = "CgZhcHBrZXkaIHhH+w534LA6cbkfupMf1IOVxhc97ewxqJF+/QyGYdRu"
	Xl_user7  string = "CgZhcHBrZXkaIJtKN8wa35kszSEJhz9KSnBSFoxNvDWgIAow4QZhDGDN"
	Xl_user8  string = "CgZhcHBrZXkaIKuygv9u5N7QBj9FIgX20cswd/n4Gk8OMJ2bUt8eCVkI"
	Xl_user9  string = "CgZhcHBrZXkaII+4uBNPsb04hCKa6uMKKX7qCtKRDk3GYPlOAfaF1GUA"
	Xl_user10 string = "CgZhcHBrZXkaIOT+NsG/4ChJaQjt+JeFxTzGGjwsgSmU9nckF6MbMsl7"

	Token1 string = "CgZhcHBrZXkaICAvo1UH53CiwR/aurQXCDBpogz9OGlWbbWDpDsMJ4dn"
	Token2 string = "CgZhcHBrZXkaIHyiKFX87ojypRsjRqk/IPYTkqTNEiuvvABITR/imPaH"
	Token3 string = "CgZhcHBrZXkaIDIBXriAVM4RyD7VLFv8vrR1+efi6LycPMuKqbQ/oVdF"
	Token4 string = "CgZhcHBrZXkaIKTH3MaZdkgLYMLsYpVmt/UT3jQkd2UgGX35LjN26ouz"
	Token5 string = "CgZhcHBrZXkaINYLoPeDJyh0HuZdk3Vx+dNs5RBD2/McgZiDjjyXS2Pm"
)

func OnMessage(msg *pbobjs.DownMsg) {
	fmt.Println("***************received msg*******************")
	fmt.Println("channel_type:", msg.ChannelType, "sender:", msg.SenderId, "\ttarget:", msg.TargetId, "\tmsg_type:", msg.MsgType, "\tmsg_content:", string(msg.MsgContent))
	bs, _ := json.Marshal(msg.TargetUserInfo)
	fmt.Println("sender_info:", string(bs))
}

func OnDisconnect(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody) {
	fmt.Println("*****************disconnect msg*************")
	fmt.Println("disconnect_code:", code, "\tdisconnect_msg:", tools.ToJson(disMsg))
}
