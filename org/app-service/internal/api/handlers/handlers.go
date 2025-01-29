package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/config"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/go-playground/validator/v10"
)

func Success(payload interface{}) *ApiResponse {
	return BuildResponse(
		xconstants.STATUS_CODE_SERVICE_SUCCESS,
		xconstants.MSG_SUCCESS,
		[]string{},
		payload)
}

func Accepted() *ApiResponse {
	return BuildResponse(
		xconstants.STATUS_CODE_SERVICE_SUCCESS,
		xconstants.MSG_SUCCESS,
		[]string{},
		nil)
}

func Deleted() *ApiResponse {
	return BuildResponse(
		xconstants.STATUS_CODE_DELETE_SUCCESS,
		xconstants.MSG_SUCCESS,
		[]string{},
		nil)
}

func ValidationErrors(errs error) *ApiResponse {
	payload := []FieldValidationError{}
	for _, err := range errs.(validator.ValidationErrors) {
		errObj := &FieldValidationError{}
		errObj.Field = err.Field()
		errObj.Namespace = err.Namespace()
		errObj.Kind = err.Kind().String()
		errObj.Value = err.Value()
		errObj.Error = fmt.Sprintf("%s %s", err.Tag(), err.Param())
		payload = append(payload, *errObj)
	}
	if config.DevModeFlag {
		str, _ := json.MarshalIndent(payload, "", "  ")
		logger.Error("ValidationErrors:")
		logger.Error(string(str))
	}
	return BuildResponse(
		xconstants.STATUS_CODE_VALIDATION_ERROR,
		xconstants.MSG_VALIDATION_ERROR,
		[]string{xconstants.MSG_VALIDATION_ERROR},
		payload)
}
