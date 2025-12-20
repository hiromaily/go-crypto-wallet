package debug

import (
	"github.com/bookerzzz/grok"
)

func Debug(value any, options ...grok.Option) {
	grok.Value(value, options...)
}
