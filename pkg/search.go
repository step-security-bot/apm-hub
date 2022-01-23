package pkg

import (
	"github.com/flanksource/commons/logger"
	"net/http"

	"github.com/flanksource/flanksource-ui/apm-hub/api"
	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
	"github.com/labstack/echo/v4"
)

func Search(c echo.Context) error {
	cc := c.(*api.Context)
	searchParams := new(logs.SearchParams)
	err := c.Bind(searchParams)
	if err != nil {
		cc.Error(err)
	}
	var searchResults []logs.SearchResults
	for _, backend := range logs.GlobalBackends {
		searchResult, err := backend.Backend.Search(searchParams)
		if err != nil {
			logger.Errorf("error executing error: %v", err)
			continue
		}
		searchResults = append(searchResults, searchResult)
	}
	return cc.JSON(http.StatusOK, searchResults)
}
