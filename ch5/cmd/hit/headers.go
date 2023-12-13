package main

import (
	"errors"
	"strings"
)

type HttpHead struct {
	key string
	val string
}

var headers []HttpHead

func headersString() string {
	var strs = []string{}
	for _, v := range headers {
		strs = append(strs, v.String())
	}
	return strings.Join(strs, ", ")
}

func toHttpHead(s *string) *HttpHead {
	var val HttpHead
	setHttpHead(&val, *s)
	return &val
}

func setHttpHead(h *HttpHead, s string) error {
	before, after, found := strings.Cut(s, ":")
	if found {
		h.key = strings.TrimSpace(before)
		h.val = strings.TrimSpace(after)
		return nil
	}
	return errors.New("invalid http head")
}

func (h *HttpHead) Set(s string) error {
	err := setHttpHead(h, s)
	if err == nil {
		headers = append(headers, *h)
	}
	return err
}

func (h *HttpHead) String() string {
	return strings.Join([]string{h.key, h.val}, ":")
}

func isValidHead(s string) bool {
	var val HttpHead
	err := val.Set(s)
	return err == nil
}
