package alphabet

type Alphabet interface {
	Len() int
	Get(...int) (string, error)
}
