package qbitorrent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Rename(hash, oldName, newName string) error {
	apiUrl := Url + "/api/v2/torrents/renameFile"
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("oldPath", oldName)
	data.Set("newPath", newName)
	bodyData := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiUrl, bodyData)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "SID="+AuthInfo.Sid) // 请替换为实际的 SID

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		return nil
	} else {
		return fmt.Errorf(`改名错误 代码:%d`, resp.StatusCode)
	}

}
