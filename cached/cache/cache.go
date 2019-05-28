package cache

// Cache ...
type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}
