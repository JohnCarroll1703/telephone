package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Pagination struct {
	PerPage    int
	Page       int
	Pagination bool
	Total      int64
}

type PaginationCtxType string

const (
	PaginationCtx PaginationCtxType = "pagination"
)

func paginationHandlerFunc(c *gin.Context) {
	perPagePath, pagePath := formQueriesFromValues(c)

	if perPagePath == "0" && pagePath == "0" {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), PaginationCtx,
			&Pagination{}))

		return
	}

	page, err := strconv.Atoi(pagePath)
	if err != nil {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), PaginationCtx,
			&Pagination{}))

		return
	}

	if page == 0 || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPagePath)
	if err != nil {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), PaginationCtx,
			&Pagination{}))

		return
	}

	if perPage > 100 || perPage <= 0 {
		perPage = 100
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), PaginationCtx,
		&Pagination{
			PerPage:    perPage,
			Page:       page,
			Pagination: true,
		}))
}

func formQueriesFromValues(c *gin.Context) (string, string) {
	perPagePath, pagePath := c.Query("per_page"), c.Query("page")

	if perPagePath == "" {
		perPagePath = "0"
	}

	if pagePath == "" {
		pagePath = "0"
	}

	return perPagePath, pagePath
}
