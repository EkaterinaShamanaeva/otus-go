package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	//content, err := ioutil.ReadAll(r) // -> io ? // читать по строке (content, lines,...)
	//if err != nil {
	//	return
	//}

	//lines := strings.Split(string(content), "\n")
	//var lines []string
	i := 0
	for {
		cnt, errCnt := bufio.NewReader(r).ReadString('\n')
		fmt.Println(cnt)
		if errCnt != nil {
			break
		}
		var user User
		if err = json.Unmarshal([]byte(cnt), &user); err != nil {
			return
		} // -> easyJSON?
		result[i] = user
		i++
		//lines = append(lines, cnt)
	}

	//for i, line := range lines {
	//	var user User
	//	if err = json.Unmarshal([]byte(line), &user); err != nil {
	//return
	//} // -> easyJSON?
	//	result[i] = user
	//}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched { // change ?
			// num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			// num++
			// result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
