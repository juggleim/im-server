package apis

import (
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/admingateway/services"
	"im-server/services/commonservices"
	"net/http"
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func QryMsgStatistic(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	statTypeStrArr := ctx.QueryArray("stat_type")
	statTypes := []commonservices.StatType{}
	for _, statTypeStr := range statTypeStrArr {
		intVal, err := tools.String2Int64(statTypeStr)
		if err == nil && intVal > 0 {
			statTypes = append(statTypes, commonservices.StatType(intVal))
		}
	}
	if len(statTypes) <= 0 {
		statTypes = append(statTypes, commonservices.StatType_Up)
		statTypes = append(statTypes, commonservices.StatType_Down)
		statTypes = append(statTypes, commonservices.StatType_Dispatch)
	}

	channelTypeStr := ctx.Query("channel_type")
	var channelType int64 = 0
	if channelTypeStr != "" {
		intVal, err := tools.String2Int64(channelTypeStr)
		if err == nil && intVal > 0 {
			channelType = intVal
		}
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := commonservices.QryMsgStatistic(appkey, statTypes, pbobjs.ChannelType(channelType), start, end)
	services.SuccessHttpResp(ctx, items)
}

func QryUserActivities(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := commonservices.QryUserActivities(appkey, start, end)
	services.SuccessHttpResp(ctx, items)
}

func QryConnectCount(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	wg := &sync.WaitGroup{}
	rpcNodes := bases.GetCluster().GetAllNodes()
	timeMarkMap := map[int64]int64{}
	lock := &sync.RWMutex{}
	for _, rpcNode := range rpcNodes {
		wg.Add(1)
		nodeName := rpcNode.Name
		go func() {
			defer wg.Done()
			code, respObj, err := services.SyncApiCall(ctx, "qry_connect_count", "", nodeName, &pbobjs.QryConnectCountReq{
				Start: start,
				End:   end,
			}, func() proto.Message {
				return &pbobjs.QryConnectCountResp{}
			})
			if err == nil && code == services.AdminErrorCode_Success && respObj != nil {
				resp := respObj.(*pbobjs.QryConnectCountResp)
				for _, item := range resp.Items {
					lock.Lock()
					var count int64 = 0
					if oldCount, exist := timeMarkMap[item.TimeMark]; exist {
						count = oldCount + item.Count
					} else {
						count = item.Count
					}
					timeMarkMap[item.TimeMark] = count
					lock.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	connectCountItems := []*commonservices.ConcurrentConnectItem{}
	for timeMark, count := range timeMarkMap {
		connectCountItems = append(connectCountItems, &commonservices.ConcurrentConnectItem{
			TimeMark: timeMark,
			Count:    count,
		})
	}
	sort.Slice(connectCountItems, func(i, j int) bool {
		return connectCountItems[i].TimeMark < connectCountItems[j].TimeMark
	})
	retItems := []interface{}{}
	for _, item := range connectCountItems {
		retItems = append(retItems, item)
	}
	ret := &commonservices.Statistics{
		Items: retItems,
	}
	services.SuccessHttpResp(ctx, ret)
}

func QryUserRegiste(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &services.ApiErrorMsg{
			Code: services.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := commonservices.QryUserRegiste(appkey, start, end)
	services.SuccessHttpResp(ctx, ret)
}
