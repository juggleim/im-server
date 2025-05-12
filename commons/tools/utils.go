package tools

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var snameRegexp *regexp.Regexp

func init() {
	rand.Seed(time.Now().UnixNano())
	snameRegexp = regexp.MustCompile("([a-z]|[0-9])([A-Z])")
}

func HashStr(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

type SegmentatedLocks struct {
	num   int
	locks []*sync.RWMutex
}

func NewSegmentatedLocks(num int) *SegmentatedLocks {
	locks := make([]*sync.RWMutex, num)
	for i := 0; i < num; i++ {
		locks[i] = &sync.RWMutex{}
	}
	return &SegmentatedLocks{
		num:   num,
		locks: locks,
	}
}

func (seg *SegmentatedLocks) GetLocks(strs ...string) *sync.RWMutex {
	key := strings.Join(strs, "_")
	hash := HashStr(key)
	return seg.locks[hash%uint32(seg.num)]
}

func RandInt(n int) int {
	return rand.Intn(n)
}

var chars string = "abcdefghigklmnopqrstuvwxyz0123456789"

func RandString(l int) string {
	ret := ""
	charLen := len(chars)
	for i := 0; i < l; i++ {
		charIndex := RandInt(charLen)
		char := chars[charIndex]
		ret = ret + string(char)
	}
	return ret
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

func ToJson(obj interface{}) string {
	bs, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bs)
}

func ToJsonBs(obj interface{}) []byte {
	bs, err := json.Marshal(obj)
	if err != nil {
		return []byte{}
	}
	return bs
}

func MapToStruct[T any](m map[string]interface{}) T {
	var t T
	data, _ := json.Marshal(m)
	_ = json.Unmarshal(data, &t)

	return t
}

func CamelToSnake(str string) string {
	snakeCase := snameRegexp.ReplaceAllStringFunc(str, func(match string) string {
		return match[:1] + "_" + match[1:]
	})
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}

func GetLocalMac() (mac string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		fmt.Println(inter.Name)
		mac := inter.HardwareAddr
		fmt.Println("MAC ===== ", mac)
	}
	fmt.Println("MAC = ", mac)
	return mac
}

func GetIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interfaces ipAddress: %v\n", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	fmt.Println("ips = ", ips)
	return ips
}

func RandStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func Array2Map(arr []string) (map[string]bool, bool) {
	ret := make(map[string]bool)
	for _, item := range arr {
		ret[item] = true
	}
	return ret, len(ret) > 0
}

func DistinctStringArray(arr []string) []string {
	ret := []string{}
	tmpMap := make(map[string]bool)
	for _, item := range arr {
		if _, exist := tmpMap[item]; !exist {
			ret = append(ret, item)
			tmpMap[item] = true
		}
	}
	return ret
}

func PureStr(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func TruncateText(str string, length int) string {
	charArr := []rune(str)
	if len(charArr) > length {
		return string(charArr[:length])
	}
	return str
}

func CreateDirs(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
