package server

import (
	"container/list"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	body := "set a \"hello world\""
	l := list.New()
	open := false
	buffer := ""
	for _, c := range body {
		if '"' == c {
			if open {
				l.PushBack(buffer)
				buffer = ""
			}
			open = !open
			continue
		}
		if ' ' == c && !open {
			if len(buffer) > 0 {
				l.PushBack(buffer)
				buffer = ""
			}
			continue
		}
		buffer += string(c)
	}
	//result = make([]string, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
