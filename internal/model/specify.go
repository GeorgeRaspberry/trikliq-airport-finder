package model

import (
	"gorm.io/gorm"
)

type OrderBy struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type SumBy struct {
	Key string `json:"key"`
	As  string `json:"as"`
}

const (
	defaultLimit  = 10
	maxLimit      = defaultLimit * 10
	defaultOffset = 0
)

func Specify(request Request) (orderBy string, err error) {
	//err is unused for now, add if needed

	if request.Metadata.OrderBy.Key != "" &&
		request.Metadata.OrderBy.Type != "" {
		orderBy = request.Metadata.OrderBy.Key + " " + request.Metadata.OrderBy.Type
	}

	return
}

func Paginate(request Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		page := request.Metadata.Page
		if page <= 0 {
			page = 1
		}

		pageSize := request.Metadata.Limit
		switch {
		case pageSize > maxLimit:
			pageSize = maxLimit
		case pageSize <= 0:
			pageSize = defaultLimit
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

/*

func Specify(request Request) (orderBy string, limit, offset int, err error) {
	orderBy = request.Metadata.OrderBy.Key + " " + request.Metadata.OrderBy.Type

	limit = request.Metadata.Limit
	if limit == 0 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	page := request.Metadata.Page
	if page <= 0 {
		page = 1
	}

	if page > limit {
		err = fmt.Errorf("illegal page value %d for limit %d", offset, limit)
		return
	}

	offset = (page - 1) * limit

	return
}
*/
