package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
)

type SensitiveWordDao struct {
	ID       int64  `gorm:"primary_key"`
	Word     string `gorm:"word"`
	WordType int    `gorm:"word_type"`
	AppKey   string `gorm:"app_key"`
}

func (word SensitiveWordDao) TableName() string {
	return "sensitivewords"
}

func (word SensitiveWordDao) BatchUpsert(items []SensitiveWordDao) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`word`,`word_type`,`app_key`)values", word.TableName())
	params := []interface{}{}

	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?);")
		} else {
			buffer.WriteString("(?,?,?),")
		}
		params = append(params, item.Word, item.WordType, item.AppKey)
	}

	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (word SensitiveWordDao) UpdateWord(appkey string, wordStr string, wordType int) error {
	return dbcommons.GetDb().Model(&word).Where("app_key=? and word=?", appkey, wordStr).Update("word_type", wordType).Error
}

func (word SensitiveWordDao) DeleteWords(appkey string, words ...string) error {
	return dbcommons.GetDb().Where("app_key=? and word in (?)", appkey, words).Delete(&word).Error
}

func (word SensitiveWordDao) QrySensitiveWords(appkey string, limit, startId int64) ([]*SensitiveWordDao, error) {
	var items []*SensitiveWordDao
	err := dbcommons.GetDb().Where("app_key=? and id>?", appkey, startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}

func (word SensitiveWordDao) QrySensitiveWordsWithPage(appkey string, page, size int64, str string, wordType int32) (words []*SensitiveWordDao, total int, err error) {
	var items []*SensitiveWordDao
	db := dbcommons.GetDb().Table(word.TableName()).Where("app_key=?", appkey)
	if str != "" {
		db = db.Where("word like ?", fmt.Sprintf("%%%s%%", str))
	}
	if wordType != 0 {
		db = db.Where("word_type=?", wordType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id asc").Offset((page - 1) * size).Limit(size).Find(&items).Error
	if err != nil {
		return nil, 0, nil
	}

	return items, total, err
}

func (word SensitiveWordDao) Total(appkey string) int {
	var count int
	err := dbcommons.GetDb().Table(word.TableName()).Where("app_key=?", appkey).Count(&count).Error
	if err != nil {
		count = 0
	}
	return count
}
