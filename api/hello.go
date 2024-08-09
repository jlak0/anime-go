package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	// 创建一个响应对象
	response := Response{
		Message: "Hello, World!",
		Status:  200,
	}

	// 设置响应头的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 将响应对象编码为 JSON，并写入响应
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
