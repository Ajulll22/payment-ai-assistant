package validation

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func fileSizeValidation(fl validator.FieldLevel) bool {
	param := fl.Param()

	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok || len(files) == 0 {
		return true
	}

	fileSize, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	fileSize = fileSize * 1024 * 1024
	for _, file := range files {
		if file.Size > int64(fileSize) {
			return false
		}
	}

	return true
}

func fileTypeValidation(fl validator.FieldLevel) bool {
	fileType := fl.Param()

	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok || len(files) == 0 {
		return true
	}

	allowedExtensions := []string{}
	switch fileType {
	case "image":
		allowedExtensions = constant.ImageFormat[:]
	case "report":
		allowedExtensions = constant.ReportFormat[:]
	}

	for _, file := range files {
		validExtension := false
		for _, ext := range allowedExtensions {
			if strings.HasSuffix(file.Filename, ext) {
				validExtension = true
				break
			}
		}

		if !validExtension {
			return false
		}
	}

	return true
}

func RegisterCustomValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("filesize", fileSizeValidation)
		v.RegisterValidation("filetype", fileTypeValidation)
	}
}
