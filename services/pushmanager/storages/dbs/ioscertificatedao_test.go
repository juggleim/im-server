package dbs

import (
	"reflect"
	"testing"
)

func TestIosCertificateJSONContract(t *testing.T) {
	typeOf := reflect.TypeOf(IosCertificateDao{})
	want := map[string]string{
		"Certificate":  "certificate",
		"CertPwd":      "cert_pwd",
		"VoipCert":     "voip_cert",
		"VoipCertPwd":  "voip_cert_pwd",
		"P8PrivateKey": "-",
	}
	for name, tag := range want {
		field, ok := typeOf.FieldByName(name)
		if !ok || field.Tag.Get("json") != tag {
			t.Fatalf("%s json tag = %q, want %q", name, field.Tag.Get("json"), tag)
		}
	}
}
