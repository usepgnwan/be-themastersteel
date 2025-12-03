package helpers

import (
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Paginate struct {
	Model interface{}
	Type  *string // ini untuk tipe yang soft delete atau engga
}

func (p *Paginate) Paginate(db *gorm.DB, e echo.Context) ResponsePaginate {
	page := 1
	limit := 10

	if pg := e.QueryParam("page"); pg != "" {
		if p, err := strconv.Atoi(pg); err == nil && p > 0 {
			page = p
		}
	}
	if l := e.QueryParam("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	offset := (page - 1) * limit

	totalRecords := record(db, p)

	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Create a slice of the model type
	sliceType := reflect.SliceOf(modelType)
	slicePtr := reflect.New(sliceType).Interface() // Create a pointer to a slice

	// Query the database into the slice
	if p.Type == nil {
		db.Model(p.Model).Limit(limit).Offset(offset).Find(slicePtr)
	} else {
		db.Unscoped().Model(p.Model).Limit(limit).Offset(offset).Find(slicePtr)
	}
	// execute
	db.Commit()
	// Get the length of the slice
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	length := sliceValue.Len()

	totalPages := int(totalRecords) / limit
	if int(totalRecords)%limit != 0 {
		totalPages++
	}

	from := offset + 1
	to := offset + length
	if to > int(totalRecords) {
		to = int(totalRecords)
	}

	// Construct the response
	response := ResponsePaginate{
		Total:       totalRecords,
		Rows:        sliceValue.Interface(), // Convert reflect.Value back to interface{}
		CurrentPage: page,
		PerPage:     limit,
		From:        from,
		To:          to,
		LastPage:    totalPages,
	}

	return response
}
func record(db *gorm.DB, p *Paginate) int64 {
	var total int64
	if p.Type == nil {
		db.Count(&total)
	} else {
		db.Unscoped().Count(&total)
	}
	return total
}
