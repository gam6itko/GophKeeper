package masterkey

type IStorage interface {
	Has() bool
	Load() ([]byte, error)
	Store([]byte) error
}
