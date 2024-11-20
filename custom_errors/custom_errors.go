package custom_errors

const (
	_ = iota
	ERR_CODE_AUTHORIZATION
	ERR_CODE_Database
	ERR_CODE_Marshal
	ERR_CODE_BodyRequired
	ERR_CODE_NotFound
	ERR_CODE_NotValid
	ERR_CODE_REQUIRED_FIELD
	ERR_CODE_Exists
	ERR_CODE_Redis
	ERR_CODE_RabbitMQ
)

var ERRORS = map[int]string{
	ERR_CODE_AUTHORIZATION:  "authorization error",
	ERR_CODE_Database:       "database error",
	ERR_CODE_Marshal:        "marshal error",
	ERR_CODE_BodyRequired:   "body required",
	ERR_CODE_NotFound:       "not found",
	ERR_CODE_NotValid:       "not valid",
	ERR_CODE_REQUIRED_FIELD: "required field is missing",
	ERR_CODE_Exists:         "already exists",
	ERR_CODE_Redis:          "redis error",
	ERR_CODE_RabbitMQ:       "rabbitmq error",
}
