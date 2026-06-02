package handling

func ResponseSuccess[T any](data T, msg string, token string) BaseResponse[T] {
	return BaseResponse[T]{
		Data:    data,
		Message: msg,
		Token:   token,
		Result:  true,
		Code:    200,
	}
}
