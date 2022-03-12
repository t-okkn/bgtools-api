package ws

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	hashids "github.com/speps/go-hashids"

	"unsafe"
)

const (
	// Hashidsに使用する文字列
	alphabet string = "abcdefghkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXY0123456789"

	// Hashidsの最低文字列長
	minLength int = 4
)

// <summary>: Hash(SHA256)値を取得します
func getHash(remote string) string {
	// hash := sha256.Sum256([]byte(str))
	hash := sha256.Sum256(*(*[]byte)(unsafe.Pointer(&remote)))
	return hex.EncodeToString(hash[:])
}

// <summary>: ConnIdを取得します
func getConnId(remote string) (string, error) {
	d := hashids.NewData()
	d.Alphabet = alphabet
	d.MinLength = minLength
	d.Salt = getHash(remote)

	hid, err := hashids.NewWithData(d)
	if err != nil {
		return "", err
	}

	ip, port, err := addressToIpPort(remote)
	if err != nil {
		return "", err
	}

	var num int64 = 0

	for i := 0; i < 4; i++ {
		st := 4 * i
		ed := 4 * (i + 1)
		num += int64(binary.BigEndian.Uint32(ip[st:ed]))
	}
	num += int64(port)

	hashid, err := hid.EncodeInt64([]int64{num})
	if err != nil {
		return "", err
	}

	com := compressedString(ip, port)

	return fmt.Sprintf("%s-%s", hashid, com), nil
}

func addressToIpPort(remote string) ([]byte, uint16, error) {
	h, p, err := net.SplitHostPort(remote)
	if err != nil {
		return []byte{}, 0, err
	}

	ip := net.ParseIP(h)
	if ip == nil {
		e := errors.New("invalid host address")
		return []byte{}, 0, e
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return []byte{}, 0, err
	}

	return []byte(ip.To16()), uint16(port), nil
}

func isCorrectConnId(connid, remote string) bool {
	h, p, err := net.SplitHostPort(remote)
	if err != nil {
		return false
	}

	ip, port := connIdToIp(connid)

	return ip == h && port == p
}

func connIdToIp(connid string) (string, string) {
	cutid := strings.Split(connid, "-")
	if len(cutid) != 2 {
		return "", ""
	}

	splited := strings.Split(cutid[1], ".")
	if len(splited) < 2 {
		return "", ""
	}

	var result strings.Builder
	result.Grow(32)

	for i, s := range splited {
		if i % 2 == 0 {
			result.WriteString(s)

		} else {
			length := len(splited[i-1])
			last := string(splited[i-1][length-1])
			count, _ := strconv.Atoi(s)

			str := make([]string, count)
			for j := 0; j < count; j++ {
				str[j] = last
			}

			result.WriteString(strings.Join(str, ""))
		}
	}

	dec, err := hex.DecodeString(result.String())
	if err != nil {
		return "", ""
	}

	ip := net.IP(dec[2:]).String()
	port := fmt.Sprint(binary.BigEndian.Uint16(dec[:2]))

	return ip, port
}

func compressedString(ip []byte, port uint16) string {
	conv := make([]byte, 2, 18)
	binary.BigEndian.PutUint16(conv, port)

	conv = append(conv, ip...)
	enc := hex.EncodeToString(conv)

	var result strings.Builder
	var buf strings.Builder
	var forward rune

	result.Grow(36)
	buf.Grow(36)

	action := func(r rune) {
		if buf.Len() == 0 {
			result.WriteRune(r)

		} else if buf.Len() <= 3 {
			result.WriteString(buf.String())
			result.WriteRune(r)

		} else {
			result.WriteRune(r)
			result.WriteString(fmt.Sprintf(".%d.", buf.Len()))
		}
	}

	for i, r := range enc {
		if i == 0 {
			forward = r
			continue
		}

		if forward == r {
			buf.WriteRune(forward)

		} else {
			action(forward)
			buf.Reset()
		}

		if i == len([]rune(enc)) - 1 {
			action(r)
		}

		forward = r
	}

	return result.String()
}
