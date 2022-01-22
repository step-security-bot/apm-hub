package pkg

import (
	"fmt"
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
	// fmt.Println(searchParams)
	// if strings.HasPrefix(strings.ToLower(searchParams.Type), "kubernetes") {
	// 	_, err := k8s.Search(searchParams)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, err.Error())
	// 	}
	// }
	var searchResults []logs.SearchResults
	for _, backend := range logs.GlobalBackends {
		fmt.Println("executing")
		searchResult, err := backend.Backend.Search(searchParams)
		if err != nil {
			fmt.Println(err)
		}
		searchResults = append(searchResults, searchResult)
	}
	return cc.JSON(http.StatusOK, searchResults)
}
