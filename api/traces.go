package api

type TraceParams struct {
	Limit int
	Page  string
}

type TraceResults struct {
	Total    int
	NextPage string
}

type TraceAPI interface {
	Search(q TraceParams) (r TraceResults, err error)
}
