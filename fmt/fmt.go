package fmt

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// T transforms and returns a hash according to pattern.
// If the pattern is not valid for the hash an empty string
// is returned. An error will be returned for a bad pattern.
func T(h, p string) (string, error) {
	var t string
	hash := bufio.NewReader(bytes.NewBufferString(h))
	pattern := bufio.NewReader(bytes.NewBufferString(p))

	for {
		pch, _, err := pattern.ReadRune()
		if err != nil {
			break
		}
		if pch != '%' {
			hch, _, err := hash.ReadRune()
			switch {
			case err != nil:
				return "", nil
			case hch != pch:
				return "", nil
			default:
				t += string(hch)
			}
			continue
		}

		var count int
		var kind rune
		var end rune

		c, _, err := pattern.ReadRune()
		if err != nil {
			return "", err
		}
		if c == '*' {
			kind, _, err = pattern.ReadRune()
			if err != nil {
				return "", err
			}
		} else {
			cs := string(c)
			for {
				d, _, err := pattern.ReadRune()
				if err != nil {
					return "", err
				}
				if isDigit(d) {
					cs += string(d)
					continue
				}
				kind = d
				break
			}
			count, err = strconv.Atoi(cs)
			if err != nil {
				return "", err
			}
		}
		end, _, err = pattern.ReadRune()
		readToEnd := err != nil

		var found int
		var h rune
		switch string(kind) {
		case "d":
			for {
				h, _, err = hash.ReadRune()
				if err != nil || !isDigit(h) || !(readToEnd || h != end) {
					break
				}
				found++
				t += string(h)
			}
		case "s":
			for {
				h, _, err := hash.ReadRune()
				if err != nil || !(readToEnd || h != end) {
					break
				}
				found++
				t += string(h)
			}
		case "l":
			for {
				h, _, err := hash.ReadRune()
				if err != nil || !(readToEnd || h != end) {
					break
				}
				found++
				t += strings.ToLower(string(h))
			}
		case "u":
			for {
				h, _, err := hash.ReadRune()
				if err != nil || !(readToEnd || h != end) {
					break
				}
				found++
				t += strings.ToUpper(string(h))
			}
		case "x":
			for {
				h, _, err := hash.ReadRune()
				if err != nil || !isHex(h) || !(readToEnd || h != end) {
					break
				}
				found++
				t += strings.ToLower(string(h))
			}
		case "X":
			for {
				h, _, err := hash.ReadRune()
				if err != nil || !isHex(h) || !(readToEnd || h != end) {
					break
				}
				found++
				t += strings.ToUpper(string(h))
			}
		default:
			return "", errors.New("Bad type in pattern")
		}
		if c != '*' && count != found {
			return "", nil
		}
		if !readToEnd {
			t += string(end)
		}
	}
	return t, nil
}

type F struct {
	Algorithm string
	Pattern   string
	Example   string
}

func Formats() map[string]F {
	return map[string]F{
		"NetNTLMV2": F{"NetNTLMV2", "%*u::%*u:%16x:%32x:%106x", "admin::N46iSNekpT:08ca45b7d7ea58ee:88dcbe4446168966a153a0064958dac6:5c7830315c7830310000000000000b45c67103d07d7b95acd12ffa11230e0000000052920b85f78d013c31cdb3b92f5d765c783030"},
	}
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isHex(c rune) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}
