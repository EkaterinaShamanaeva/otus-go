package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

type email []string

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUserEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUserEmails(r io.Reader) (email, error) {
	result := make(email, 0, 100000)

	scanner := bufio.NewScanner(r)
	var p fastjson.Parser
	for scanner.Scan() {
		v, errF := p.Parse(scanner.Text())
		if errF != nil {
			return nil, errF
		}
		result = append(result, string(v.GetStringBytes("Email")))
	}
	return result, nil
}

func countDomains(u email, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, mail := range u {
		matched := strings.Contains(mail, "."+domain)

		if matched {
			result[strings.ToLower(strings.SplitN(mail, "@", 2)[1])]++
		}
	}

	return result, nil
}
