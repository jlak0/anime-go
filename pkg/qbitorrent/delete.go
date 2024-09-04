package qbitorrent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Delete(hash string) error {
	sid := AuthInfo.getSid()
	endPoint := "/api/v2/torrents/delete" // 请替换为实际的 base_url

	data := url.Values{}
	data.Set("hashes", hash)
	data.Set("deleteFiles", "true")

	req, err := http.NewRequest("POST", Url+endPoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "SID="+sid) // 请替换为实际的 SID

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(`删除种子请求错误%s`, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`删除种子错误:%d`, resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}
