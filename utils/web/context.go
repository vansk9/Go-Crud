package web

import "net/http"

func GetUserID(r *http.Request) uint {
	if uid, ok := r.Context().Value("user_id").(uint); ok {
		return uid
	}
	return 0
}
