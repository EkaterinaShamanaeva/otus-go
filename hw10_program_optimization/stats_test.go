//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetUsers(t *testing.T) {
	userFirst := User{
		ID:       1,
		Name:     "Howard Mendoza",
		Username: "0Oliver",
		Email:    "aliquid_qui_ea@Browsedrive.gov",
		Phone:    "6-866-899-36-79",
		Password: "InAQJvsq",
		Address:  "Blackbird Place 25",
	}
	userSecond := User{
		ID:       2,
		Name:     "Justin Oliver Jr. Sr. I II III IV V MD DDS PhD DVM",
		Username: "oPerez",
		Email:    "MelissaGutierrez@Twinte.biz",
		Phone:    "106-05-18",
		Password: "f00GKr9i",
		Address:  "Oak Valley Lane 19",
	}
	userThird := User{
		ID:       3,
		Name:     "Brian Olson",
		Username: "non_quia_id",
		Email:    "FrancesEllis@Quinu.edu",
		Phone:    "237-75-34",
		Password: "cmEPhX8",
		Address:  "Butterfield Junction 74",
	}
	// expectedUsers := users{userFirst, userSecond, userThird}
	expectedUsers := email{userFirst.Email, userSecond.Email, userThird.Email}

	var file *os.File
	file, err := os.OpenFile("testdata/get_users_test.dat", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	resultGetUsers, errGetUsers := getUsers(file)

	if err = file.Close(); err != nil {
		log.Fatal(err)
	}

	require.NoError(t, errGetUsers)
	require.Equal(t, expectedUsers, resultGetUsers)
}

func TestCountDomains(t *testing.T) {
	userFirst := User{
		ID:       1,
		Name:     "Howard Mendoza",
		Username: "0Oliver",
		Email:    "aliquid_qui_ea@Browsedrive.gov",
		Phone:    "6-866-899-36-79",
		Password: "InAQJvsq",
		Address:  "Blackbird Place 25",
	}
	userSecond := User{
		ID:       2,
		Name:     "Jesse Vasquez",
		Username: "qRichardson",
		Email:    "mLynch@broWsecat.com",
		Phone:    "9-373-949-64-00",
		Password: "SiZLeNSGn",
		Address:  "Fulton Hill 80",
	}
	userThird := User{
		ID:       3,
		Name:     "Clarence Olson",
		Username: "RachelAdams",
		Email:    "RoseSmith@Browsecat.com",
		Phone:    "988-48-97",
		Password: "71kuz3gA5w",
		Address:  "Monterey Park 39",
	}
	// users := users{userFirst, userSecond, userThird}
	users := email{userFirst.Email, userSecond.Email, userThird.Email}
	domainCom := "com"

	expected := make(DomainStat)
	expected["browsecat.com"] = 2

	result, err := countDomains(users, domainCom)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}
