package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = strings.ToLower(domain)
	br := bufio.NewReader(r)
	var user User
	for {
		line, isPrefix, err := br.ReadLine()
		if isPrefix {
			buf := append([]byte{}, line...)
			for isPrefix {
				var extra []byte
				extra, isPrefix, err = br.ReadLine()
				buf = append(buf, extra...)
			}
			line = buf
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if bytes.IndexByte(line, 0x40) == -1 {
			continue
		}

		if err := easyjson.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		if user.Email == "" {
			continue
		}

		at := strings.IndexByte(user.Email, '@')
		if at <= 0 || at == len(user.Email)-1 {
			continue
		}

		eDomain := strings.ToLower(user.Email[at+1:])

		if strings.HasSuffix(eDomain, "."+domain) {
			result[eDomain]++
		}
	}
	return result, nil
}
