package models_test

import (
	"anime-go/models"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	f := models.Find()
	fmt.Println(len(*f))
	for _, v := range *f {
		fmt.Println(v)
	}
}
