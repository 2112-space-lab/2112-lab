package fx

import "fmt"

// FlattenErrorsIfAny iterates over input errors and wrap them all into a single error
func FlattenErrorsIfAny(errs ...error) error {
	var errSum error
	errCount := 0
	for _, err := range errs {
		if err == nil {
			continue
		}
		errCount++
		if errSum != nil {
			errSum = fmt.Errorf("[%d:%s] %w", errCount, err.Error(), errSum)
		} else {
			errSum = fmt.Errorf("[%d:%w]", errCount, err)
		}
	}
	return errSum
}

// FlattenErrorsAsStringIfAny iterates over input errors and wrap them all into a single string error
func FlattenErrorsAsStringIfAny(errs ...error) string {
	errMessage := ""
	flattenErr := FlattenErrorsIfAny(errs...)
	if flattenErr != nil {
		errMessage = flattenErr.Error()
	}
	return errMessage
}

// FlattenErrorsIfAnyWithPath iterates over input errors and wrap them all into a single error with the provided function path
func FlattenErrorsIfAnyWithPath(path string, errs ...error) error {
	err := FlattenErrorsIfAny(errs...)
	if err != nil {
		return fmt.Errorf("%s [%w]", path, err)
	}
	return nil
}
