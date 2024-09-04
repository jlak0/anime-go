// xx_test.go
package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := "[ANi] 這是妳與我的最後戰場，或是開創世界的聖戰 第二季 - 02 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4]"
	expected := AnimeInfo{}

	result, err := Parse(input)
	if err != nil {
		t.Errorf("Parse(%q) unexpected error: %v", input, err)
	} else if result.NameZh != "這是妳與我的最後戰場，或是開創世界的聖戰" {
		t.Errorf("Parse(%q) = %q; expected %q", input, result, expected)
	}
}
