package interceptors

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"regexp"
	"strings"
)

func ConditionMatchs(conditions []*Condition, senderId, receiverId string, channelType pbobjs.ChannelType, msgType string, msgContent []byte) bool {
	if len(conditions) <= 0 {
		return false
	}
	match := false
	for _, condition := range conditions {
		isMatch := condition.Match(senderId, receiverId, channelType, msgType, msgContent)
		if isMatch {
			return true
		}
	}
	return match
}

type Condition struct {
	ChannelTypeChecker Matcher
	MsgTypeChecker     Matcher
	SenderIdChecker    Matcher
	ReceiverIdChecker  Matcher
	MsgContentChecker  Matcher
}

func (condition *Condition) Match(senderId, receiverId string, channelType pbobjs.ChannelType, msgType string, msgContent []byte) bool {
	ret := true
	//channel_type
	if condition.ChannelTypeChecker != nil {
		isMatch := condition.ChannelTypeChecker.Match(tools.Int642String(int64(channelType)))
		ret = ret && isMatch
	} else {
		ret = ret && true
	}
	//msg_type
	if condition.MsgTypeChecker != nil {
		isMatch := condition.MsgTypeChecker.Match(msgType)
		ret = ret && isMatch
	} else {
		ret = ret && true
	}
	//sender_id
	if condition.SenderIdChecker != nil {
		isMatch := condition.SenderIdChecker.Match(senderId)
		ret = ret && isMatch
	} else {
		ret = ret && true
	}
	//receiver_id
	if condition.ReceiverIdChecker != nil {
		isMatch := condition.ReceiverIdChecker.Match(receiverId)
		ret = ret && isMatch
	} else {
		ret = ret && true
	}
	//msg content
	if condition.MsgContentChecker != nil {
		isMatch := condition.MsgContentChecker.Match(string(msgContent))
		ret = ret && isMatch
	} else {
		ret = ret && true
	}
	return ret
}

type Matcher interface {
	Match(val string) bool
}

func CreateMatcher(val string) Matcher {
	if val == "" || val == "*" {
		return &NilMatcher{}
	} else if strings.HasPrefix(val, "contains(") {
		values, err := extractContainsValues(val)
		if err != nil {
			return &NilMatcher{}
		}
		return NewContainsChecker(values)
	} else if strings.HasPrefix(val, "regex(") {
		value, err := extractRegexValues(val)
		if err != nil {
			return &NilMatcher{}
		}
		return NewRegexChecker(value)
	} else {
		return &EqualMatcher{
			value: val,
		}
	}
}

func extractContainsValues(input string) ([]string, error) {
	re := regexp.MustCompile(`contains\(([^)]+)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no matches found")
	}
	values := strings.Split(matches[1], ",")

	return values, nil
}

func extractRegexValues(input string) (string, error) {
	re := regexp.MustCompile(`regex\(([^)]+)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", fmt.Errorf("no matches found")
	}
	return matches[1], nil
}

// nil matcher
type NilMatcher struct {
}

func (checker *NilMatcher) Match(val string) bool {
	return true
}

// equal matcher
type EqualMatcher struct {
	value string
}

func NewEqualChecker(val string) *EqualMatcher {
	return &EqualMatcher{
		value: val,
	}
}

func (checker *EqualMatcher) Match(val string) bool {
	return checker.value == val
}

// contains matcher
type ContainsMatcher struct {
	values map[string]struct{}
}

func NewContainsChecker(vals []string) *ContainsMatcher {
	m := &ContainsMatcher{
		values: make(map[string]struct{}, len(vals)),
	}
	for _, val := range vals {
		m.values[val] = struct{}{}
	}
	return m
}

func (checker *ContainsMatcher) Match(val string) bool {
	if _, ok := checker.values[val]; ok {
		return true
	}
	return false
}

// regex matcher
type RegexMatcher struct {
	reg *regexp.Regexp
}

func NewRegexChecker(reg string) *RegexMatcher {
	return &RegexMatcher{
		reg: regexp.MustCompile(reg),
	}
}

func (checker *RegexMatcher) Match(val string) bool {
	if checker.reg != nil {
		return checker.reg.MatchString(val)
	}
	return false
}
