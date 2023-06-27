package dicts

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	Estate       = estate{make(map[string]bool)}
	Swears       = swears{make(map[string]bool)}
	Chats        = chats{make(map[int64]struct{})}
	Admins       = chats{make(map[int64]struct{})}
	RegexPattern *regexp.Regexp
	Mu           sync.Mutex
)

type estate struct {
	Map map[string]bool
}

type swears struct {
	Map map[string]bool
}

type chats struct {
	Map map[int64]struct{}
}

type admins struct {
	Map map[int64]struct{}
}

func InitPattern() {
	regexPattern := "(" + strings.Join(keys(Estate.Map), "|") + ")"
	RegexPattern = regexp.MustCompile(regexPattern)
	fmt.Println(regexPattern)
}

func (e estate) Check(text string) bool {
	text = strings.ToLower(text)
	fmt.Println(text)
	return RegexPattern.MatchString(text)
}

func (s swears) Check(text string) bool {
	words := strings.Fields(text)
	for _, word := range words {
		if s.Map[word] {
			return true
		}
	}
	return false
}

func (c chats) Check(chatId int64) bool {
	_, ok := c.Map[chatId]
	return ok
}

func (a admins) Check(userId int64) bool {
	_, ok := a.Map[userId]
	return ok
}

func keys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
