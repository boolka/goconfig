package entry

type Entry interface {
	Get(string) (any, bool)
}
