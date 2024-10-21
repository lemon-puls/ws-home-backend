package page

import "gorm.io/gorm"

type PageParam struct {
	Page    int    `json:"page" form:"page"`         // 当前页码，默认从 1 开始
	Limit   int    `json:"limit" form:"limit"`       // 每页显示的条数
	OrderBy string `json:"order_by" form:"order_by"` // 排序字段
	// 排序方式，asc 升序，desc 降序，默认 desc
	Order string `json:"order" form:"order"`
}

type PageResult struct {
	Total   int64       `json:"total"`   // 总记录数
	Page    int         `json:"page"`    // 当前页码
	Limit   int         `json:"limit"`   // 每页记录数
	Records interface{} `json:"records"` // 当前页记录
}

func Paginate(db *gorm.DB, param PageParam, result interface{}) (*PageResult, error) {
	// 计算偏移量
	offset := (param.Page - 1) * param.Limit

	// 统计总记录数
	var total int64
	if err := db.Model(result).Count(&total).Error; err != nil {
		return nil, err
	}

	orderStr := ""
	if param.OrderBy != "" {
		orderStr = param.OrderBy + " "
		if param.Order == "asc" {
			orderStr += "asc"
		} else {
			orderStr += "desc"
		}
	}

	if orderStr != "" {
		db = db.Order(orderStr)
	}

	// 查询当前页数据
	if err := db.Limit(param.Limit).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}

	// 组装分页结果
	pageResult := &PageResult{
		Total:   total,
		Page:    param.Page,
		Limit:   param.Limit,
		Records: result,
	}

	return pageResult, nil
}
