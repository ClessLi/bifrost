package code

//go:generate codegen -type=int -fullname=Bifrost
//go:generate codegen -type=int -fullname=Bifrost -doc -output ../../../docs/guide/zh-CN/api/error_code_generated.md

// Common: basic errors.
// Code must start with 1xxxxx.
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound

	// ErrRequestTimeout - 408: Request timeout.
	ErrRequestTimeout
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

// common: data repository errors.
const (
	// ErrDataRepository - 500: Data Repository error.
	ErrDataRepository int = iota + 100201
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100301

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader

	// ErrUserOrPasswordIncorrect - 401: User or Password was incorrect.
	ErrUserOrPasswordIncorrect

	// ErrPermissionDenied - 403: Permission denied.
	ErrPermissionDenied

	// ErrAuthnClientInitFailed - 500: The `Authentication` client initialization failed.
	ErrAuthnClientInitFailed

	// ErrAuthClientNotInit - 500: The `Authentication` and `Authorization` client not initialized.
	ErrAuthClientNotInit

	// ErrConnToAuthServerFailed - 500: Failed to connect to `Authentication` and `Authorization` server.
	ErrConnToAuthServerFailed
)

// common: encode/decode errors.
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100401

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)
