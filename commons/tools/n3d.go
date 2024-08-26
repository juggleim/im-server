package tools

import (
	"errors"
	"math"
	"strings"
)

var IdEncryptor *N3d

type N3d struct {
	keyCode int64
	key     string
	radix   int64
	lower   int64
	upper   int64
	dict    [62][62]byte
}

func init() {
	IdEncryptor, _ = NewN3d(1, 4294967295)
}

func NewN3d(lower, upper int64) (*N3d, error) {
	if upper <= lower {
		return nil, errors.New("Parameter is error")
	}
	key := "11EdDIauqcim"
	charMap := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	n := new(N3d)
	n.key = key
	n.lower = lower
	n.upper = upper
	n.radix = int64(len(charMap))

	i, m, ref := 0, 0, len(key)

	for {
		if 0 <= ref {
			if m >= ref {
				break
			}
		} else {
			if m <= ref {
				break
			}
		}

		a := key[i]
		if a > 127 {
			return nil, errors.New("The key is error.")
		}

		n.keyCode += int64(a) * int64(math.Pow(float64(128), float64(i%7)))

		if ref >= 0 {
			i = m + 1
			m = m + 1
		} else {
			i = m - 1
			m = m - 1
		}
	}

	if n.keyCode+n.radix < n.upper {
		return nil, errors.New("The secret key is too short.")
	}

	j, k := 0, n.keyCode-n.radix
	for k < n.keyCode {
		a, b := n.radix, 0
		for a > 0 {
			s := k % a
			n.dict[j][b] = charMap[s]

			charMap[s] = charMap[a-1]

			a--
			b++
		}

		for x := 0; x < len(charMap); x++ {
			charMap[x] = n.dict[j][x]
		}

		j++
		k++
	}
	return n, nil
}

func (n *N3d) Encrypt(num int64) (string, error) {
	if num > n.upper || num < n.lower {
		return "", errors.New("Parameter is error.")
	}

	num = n.keyCode - num
	m := num % n.radix
	mapping := n.dict[m]

	var s int64
	var result strings.Builder
	result.WriteByte(n.dict[0][m])
	for num > n.radix {
		num = (num - m) / n.radix
		m = num % n.radix

		if s += m; s >= n.radix {
			s -= n.radix
		}
		result.WriteByte(mapping[s])
	}
	return result.String(), nil
}

func (n *N3d) Decrypt(str string) (int64, error) {
	if str == "" {
		return 0, errors.New("Parameter is error.")
	}

	chars := []byte(str)
	l, t, s := len(chars), 0, int64(0)
	result := int64(strings.IndexByte(string(n.dict[0][:]), chars[0]))
	if result < 0 {
		return 0, errors.New("Invalid string.")
	}

	mapping := string(n.dict[result][:])

	i, m, ref := 1, 1, l
	for {
		if 1 <= ref {
			if m >= ref {
				break
			}
		} else {
			if m <= ref {
				break
			}
		}

		j := strings.IndexByte(mapping, chars[i])
		if j < 0 {
			return 0, errors.New("Invalid string.")
		}

		s = int64(j - t)
		if s < 0 {
			s += n.radix
		}

		result += s * int64(math.Pow(float64(n.radix), float64(i)))
		t = j

		if 1 <= ref {
			i = m + 1
			m = m + 1
		} else {
			i = m - 1
			m = m - 1
		}
	}

	result = n.keyCode - result

	return result, nil
}

func EncodeInt(num int64) (string, error) {
	// n, err := NewN3d(1, 4294967295)
	// if err != nil {
	// 	return "", err
	// }

	return IdEncryptor.Encrypt(num)
}

func DecodeInt(str string) (int64, error) {
	// n, err := NewN3d(1, 4294967295)
	// if err != nil {
	// 	return 0, err
	// }
	ret, err := IdEncryptor.Decrypt(str)
	if err != nil {
		return 0, err
	}
	//       4294967295
	if ret > 2147483647 {
		return 0, nil
	}
	return ret, err
}
