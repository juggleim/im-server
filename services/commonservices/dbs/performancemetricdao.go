package dbs

import "im-server/commons/dbcommons"

type PerformanceMetricDao struct {
	ID          int64   `gorm:"column:id;primaryKey"`
	NodeName    string  `gorm:"column:node_name"`
	CollectTime int64   `gorm:"column:collect_time"`
	MetricType  string  `gorm:"column:metric_type"`
	MetricValue float64 `gorm:"column:metric_value"`
}

func (metric PerformanceMetricDao) TableName() string {
	return "performance_metrics"
}

func (metric PerformanceMetricDao) BatchInsert(items []PerformanceMetricDao) error {
	if len(items) <= 0 {
		return nil
	}
	return dbcommons.GetDb().CreateInBatches(items, 100).Error
}

func (metric PerformanceMetricDao) DeleteBeforeCollectTime(nodeName string, collectTime int64) error {
	return dbcommons.GetDb().Where("node_name = ? and collect_time < ?", nodeName, collectTime).Delete(&PerformanceMetricDao{}).Error
}
