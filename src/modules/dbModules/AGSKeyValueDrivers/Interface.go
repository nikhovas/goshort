package AGSKeyValueDrivers

type AGSKeyValueInterface interface {
	Get(key string) ([]byte, error)
	SetIfNotExists(key string, data []byte) error
	GetNewUniqueNumber() (int, error)
}
