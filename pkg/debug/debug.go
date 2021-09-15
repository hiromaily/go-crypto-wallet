package debug

import (
	"github.com/bookerzzz/grok"
)

func Debug(value interface{}, options ...grok.Option) {
	grok.Value(value, options...)
}
