package qbitorrent

import (
	"fmt"
	"net/http"
	"net/url"

	"strings"
)

func Add(hash string, path string) error {
	endPoint := "/api/v2/torrents/add" // 请替换为实际的 base_url

	data := url.Values{}
	data.Set("urls", "magnet:?xt=urn:btih:"+hash)
	data.Set("category", "AnimeGo")
	data.Set("savepath", "AnimeGo"+path)

	req, err := http.NewRequest("POST", Url+endPoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "SID="+AuthInfo.Sid) // 请替换为实际的 SID

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(`添加种子请求错误%s`, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`添加种子错误`)
	}
	defer resp.Body.Close()

	return nil
}
