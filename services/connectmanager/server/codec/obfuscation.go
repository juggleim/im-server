package codec

import (
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/imcontext"
)

func GetObfuscationCodeFromCtx(ctx imcontext.WsHandleContext) ([8]byte, bool) {
	obfuscationCodeObj := imcontext.GetContextAttr(ctx, imcontext.StateKey_ObfuscationCode)
	if obfuscationCodeObj != nil {
		obfuscationCode := obfuscationCodeObj.([8]byte)
		return obfuscationCode, true
	}
	return [8]byte{0, 0, 0, 0, 0, 0, 0, 0}, false
}

var fixedConnMsgBytes []byte

func getFixedConnMsgBytes() []byte {
	if len(fixedConnMsgBytes) != 8 {
		connMsg := &ConnectMsgBody{
			ProtoId: ProtoId,
		}
		bs, err := tools.PbMarshal(connMsg)
		if err == nil {
			fixedConnMsgBytes = bs[:8]
		}
	}
	return fixedConnMsgBytes
}
func CalObfuscationCode(connectData []byte) [8]byte {
	fixedConnBytes := getFixedConnMsgBytes()
	code := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	if len(connectData) > 8 && len(fixedConnBytes) == 8 {
		for i := 0; i < 8; i++ {
			code[i] = connectData[i] ^ fixedConnBytes[i]
		}
	}
	return code
}

func DoObfuscation(code [8]byte, data []byte) {
	dataLen := len(data)
	for i := 0; i < dataLen; i++ {
		data[i] = data[i] ^ code[i%8]
	}
}
