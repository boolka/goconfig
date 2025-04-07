package entry

import "context"

type Entry interface {
	Get(context.Context, string) (any, bool)
}
