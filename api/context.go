package api

import (
	"github.com/flanksource/kommons"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	Kommons *kommons.Client
}

// var ctx Context
