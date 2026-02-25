package apperror

// IsNotFound checks if an error is a NotFound error
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeNotFound
	}
	return false
}

// IsConflict checks if an error is a Conflict error
func IsConflict(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeConflict
	}
	return false
}

// IsUnauthorized checks if an error is an Unauthorized error
func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeUnauthorized
	}
	return false
}

// IsForbidden checks if an error is a Forbidden error
func IsForbidden(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == CodeForbidden
	}
	return false
}
