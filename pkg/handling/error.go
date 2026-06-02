package handling

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/logger"
)

type ErrorChan struct {
	Source string
	Err    error
}

type ErrorWrapper struct {
	Message    string   `json:"message"` // human readable error
	Validation []string `json:"-"`       //
	Code       int      `json:"-"`       // code
	Err        error    `json:"-"`       // original error
	Filename   string   `json:"-"`
	LineNumber int      `json:"-"`
}

func (w *ErrorWrapper) Error() string {
	// guard against panics
	var messages []string
	var err error = w

	for err != nil {
		var ew *ErrorWrapper
		if errors.As(err, &ew) {
			// Bangun pesan dengan optional filename dan line number
			msg := ew.Message
			if ew.Filename != "" || ew.LineNumber != 0 {
				msg = fmt.Sprintf("%s (%s:%d)", ew.Message, ew.Filename, ew.LineNumber)
			}
			messages = append(messages, msg)
			err = ew.Err
		} else {
			// Error terakhir yang bukan ErrorWrapper
			messages = append(messages, err.Error())
			break
		}

		if err == nil {
			break
		}
	}

	return strings.Join(messages, " -> ")
}

func NewErrorWrapper(code int, msg string, validation []string, err error) *ErrorWrapper {
	// getting previous call stack file and line info
	_, filename, line, _ := runtime.Caller(1)
	return &ErrorWrapper{
		Code:       code,
		Message:    msg,
		Err:        err,
		Validation: validation,
		Filename:   filename,
		LineNumber: line,
	}
}

func ResponseError(ctx context.Context, cfg constant.LogConfig, err error, token string) BaseResponse[any] {
	message := constant.MessageInternalServer
	code := constant.CodeInternalServer
	failed := []string{}

	var ew *ErrorWrapper
	if errors.As(err, &ew) {
		code = ew.Code
		failed = ew.Validation

		switch code {
		case constant.CodeInvalidToken:
			message = constant.MessageInvalidToken
		case constant.CodeNotAcceptable:
			message = constant.MessageNotAcceptable
		case constant.CodeUnprocessableEntity:
			message = constant.MessageInvalidParameter
		case constant.CodeInternalServer:
			message = constant.MessageInternalServer
		case constant.CodeDBError:
			message = constant.MessageInternalServer
		default:
			message = ew.Message
		}
	}

	if code >= 500 {

		log := logger.FromContext(ctx, cfg)
		errMessage := err.Error()

		log.App.ErrorWithoutTrace(errMessage)

	}

	return BaseResponse[any]{
		Token:   token,
		Result:  false,
		Failed:  failed,
		Message: message,
		Code:    code,
	}
}
