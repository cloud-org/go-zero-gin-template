package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func ObjectIdFromHexErr(id string) error {
	return fmt.Errorf("obejctId: %s 转换失败", id)
}

func NoDocumentErr(id string) error {
	return fmt.Errorf("id: %v 数据不存在", id)
}

//GetCond add dtime cond
func GetCond(cond bson.M) bson.M {
	cond["dtime"] = ""
	return cond
}

type Stage struct {
	SortStage  bson.M
	SkipStage  bson.M
	LimitStage bson.M
	Page       int
	PageSize   int
}

// StageNeedPagination 判断聚合查询是否需要分页
func StageNeedPagination(stage *Stage) bool {
	if stage.SkipStage != nil && stage.LimitStage != nil {
		return true
	}

	return false
}

//GetStage 除了 match 之外的其他条件 sort limit 等
func GetStage(page int, pageSize int) *Stage {

	order := "desc"

	sortStage := bson.M{
		"$sort": bson.M{
			"ctime": func() int {
				if order == "desc" {
					return -1
				}
				return 1
			}(),
		},
	}

	if page == 0 || pageSize == 0 {
		return &Stage{
			SortStage: sortStage,
		}
	}

	// 跳过
	skipStage := bson.M{
		"$skip": (page - 1) * pageSize,
	}

	// 截取
	limitStage := bson.M{
		"$limit": pageSize,
	}

	return &Stage{
		SortStage:  sortStage,
		SkipStage:  skipStage,
		LimitStage: limitStage,
		Page:       page,
		PageSize:   pageSize,
	}
}
