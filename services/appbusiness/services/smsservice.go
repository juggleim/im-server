package services

func CheckPhoneSmsCode(phone, code string) bool {
	return code == "000000"
}
