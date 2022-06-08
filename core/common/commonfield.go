package common

import "time"

const TimeFormat = "2006-01-02 15:04:05"

// CommonField omitempty 加上这个标签是为了更新的时候不指定 ctime 等字段不会出现 ctime 字段被覆盖的问题
type CommonField struct {
	Ctime    string `bson:"ctime,omitempty"`    // 创建时间
	Mtime    string `bson:"mtime,omitempty"`    // 修改时间
	Creator  string `bson:"creator,omitempty"`  // 创建人
	Modifier string `bson:"modifier,omitempty"` // 修改人
	DTime    string `bson:"dtime"`              // 删除时间
	Deleter  string `bson:"deleter"`            // 删除人
}

type CreateCommonField struct {
	Ctime   string `bson:"ctime"`
	Creator string `bson:"creator"`
}

type ModifyCommonField struct {
	Mtime    string `bson:"mtime"`
	Modifier string `bson:"modifier"`
}

type DeleteCommonField struct {
	DTime   string `bson:"dtime"` // 删除时间
	Deleter string `bson:"deleter"`
}

func GenCreateCommonField(userName string) *CommonField {
	now := time.Now().Format(TimeFormat)
	return &CommonField{
		Ctime:    now,
		Creator:  userName,
		Mtime:    now,
		Modifier: userName,
	}
}

func GenModifyCommonField(userName string) *CommonField {
	now := time.Now().Format(TimeFormat)
	return &CommonField{
		Mtime:    now,
		Modifier: userName,
	}
}

func GenDeleteCommonField(userName string) *DeleteCommonField {
	now := time.Now().Format(TimeFormat)
	d := DeleteCommonField{
		DTime:   now,
		Deleter: userName,
	}
	return &d
}
