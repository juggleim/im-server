package dbs

import "im-server/services/message/storages/models"

// mergeCmdCompatMsgs merges two ascending query results. Legacy rows win ties
// because compatibility reads query the legacy table first.
func mergeCmdCompatMsgs(legacy, split []*models.Msg, count int) []*models.Msg {
	if count <= 0 {
		return []*models.Msg{}
	}
	merged := make([]*models.Msg, 0, count)
	legacyIndex := 0
	splitIndex := 0
	for len(merged) < count && (legacyIndex < len(legacy) || splitIndex < len(split)) {
		var item *models.Msg
		if splitIndex >= len(split) || (legacyIndex < len(legacy) && legacy[legacyIndex].SendTime <= split[splitIndex].SendTime) {
			item = legacy[legacyIndex]
			legacyIndex++
		} else {
			item = split[splitIndex]
			splitIndex++
		}
		if item == nil {
			continue
		}
		merged = append(merged, item)
	}
	return merged
}
