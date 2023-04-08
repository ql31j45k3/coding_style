package domain

import "layout_2/internal/libs/response"

var (
	ErrNickNameIsEmpty = response.Err(response.Code100401, "nickName 是必填項目")
)
