package routes

import (
	"net/url"
)

type Route interface {
	URL(pairs ...string) (*url.URL, error)
}