package qbitorrent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Data struct {
	Availability int     `json:"availability"`
	Index        int     `json:"index"`
	IsSeed       bool    `json:"is_seed"`
	Name         string  `json:"name"`
	PieceRange   [2]int  `json:"piece_range"`
	Priority     int     `json:"priority"`
	Progress     float64 `json:"progress"`
	Size         int64   `json:"size"`
}

func GetFileName(hash string) (string, error) {
	apiUrl := Url + "/api/v2/torrents/files"
	data := url.Values{}
	data.Add("hash", hash)
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf(`改名错误:%s`, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "SID="+AuthInfo.getSid()) // 请替换为实际的 SID
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf(`改名错误:%s`, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(`改名错误 代码:%d`, resp.StatusCode)
	}
	var d []Data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(`改名错误:%s`, err)
	}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return "", fmt.Errorf(`改名错误:%s`, err)
	}

	if len(d) == 1 {
		return d[0].Name, nil
	} else if len(d) == 0 {
		return "", fmt.Errorf(`无文件`)
	} else {
		return "", fmt.Errorf(`非单一文件`)
	}

}
