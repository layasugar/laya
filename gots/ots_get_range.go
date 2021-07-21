package gots

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

// 主键范围查询多条
func (c *tableStoreClient) GetRange(start, end int64) *getRangeRowsRequest {
	var ir = &getRangeRowsRequest{
		dwSearchOrm: c.dwSearchOrm,
		tableName:   c.tableName,
		startValue:  start,
		endValue:    end,
	}

	return ir
}

func (c *getRangeRowsRequest) Fields(columns []string) *getRangeRowsRequest {
	c.columns = columns
	return c
}

func (c *getRangeRowsRequest) Order(direction string) *getRangeRowsRequest {
	c.direction = direction
	return c
}

func (c *getRangeRowsRequest) Limit(limit int32) *getRangeRowsRequest {
	c.limit = limit
	return c
}

func (c *getRangeRowsRequest) Do() (*getRangeRowsResp, error) {
	var err error
	if c.err != nil {
		return nil, err
	}

	if c.dwSearchOrm.err != nil {
		return nil, c.dwSearchOrm.err
	}

	if c.limit == 0 {
		c.limit = 10
	}

	if c.direction == "" {
		c.direction = directionForward
	}

	if c.startValue == 0 {
		return nil, errors.New("pk开始值不能为空")
	}

	if c.endValue == 0 {
		return nil, errors.New("pk结束值不能为空")
	}

	getRangeRequest := &tablestore.GetRangeRequest{}
	rangeRowQueryCriteria := &tablestore.RangeRowQueryCriteria{}
	rangeRowQueryCriteria.TableName = c.tableName
	rangeRowQueryCriteria.MaxVersion = 1
	rangeRowQueryCriteria.Limit = c.limit

	if c.direction == directionForward {
		rangeRowQueryCriteria.Direction = tablestore.FORWARD
	} else {
		rangeRowQueryCriteria.Direction = tablestore.BACKWARD
	}

	startPK := new(tablestore.PrimaryKey)
	startPK.AddPrimaryKeyColumn(pkName, c.startValue)
	rangeRowQueryCriteria.StartPrimaryKey = startPK
	endPK := new(tablestore.PrimaryKey)
	endPK.AddPrimaryKeyColumn(pkName, c.endValue)
	rangeRowQueryCriteria.EndPrimaryKey = endPK

	if len(c.columns) > 0 {
		rangeRowQueryCriteria.ColumnsToGet = c.columns
	}

	getRangeRequest.RangeRowQueryCriteria = rangeRowQueryCriteria

	getRangeResp, err := c.dwSearchOrm.client.GetRange(getRangeRequest)
	if err != nil {
		return nil, err
	}
	var res getRangeRowsResp
	if getRangeResp != nil && getRangeResp.Rows != nil && len(getRangeResp.Rows) > 0 {
		for _, rowsResult := range getRangeResp.Rows {
			if rowsResult.Columns != nil && len(rowsResult.Columns) > 0 {
				res.Columns = append(res.Columns, rowsResult.Columns)
			}
		}
	}

	if getRangeResp.NextStartPrimaryKey != nil && len(getRangeResp.NextStartPrimaryKey.PrimaryKeys) > 0 {
		for _, key := range getRangeResp.NextStartPrimaryKey.PrimaryKeys {
			res.NextPk = append(res.NextPk, key)
		}
	}
	return &res, nil
}
