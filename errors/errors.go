package errors

const (
	ECONFLICT       = "conflict"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	EFORBIDDEN      = "forbidden"
	EEXPECTED       = "expected"
	ETIMEOUT        = "timeout"
	ECONNECTIONDIAL = "connection dial"
)

type Error struct {
	Code    string
	Message string
	Op      string
	Err     error
	Detail  []byte
}
