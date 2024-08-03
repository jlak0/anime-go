package controller

import "fmt"

func Update() {
	r, err := Parse("[ANi] 这是妳与我的最后战场，或是开创世界的圣战 第二季 - 03 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4]")
	if err == nil {
		fmt.Print(r)
	}
}
