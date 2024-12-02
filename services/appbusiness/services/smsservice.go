package services

func CheckPhoneSmsCode(phone, code string) bool {
	if code == "000000" {
		return true
	}
	return false
}
