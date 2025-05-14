package commonservices

import (
	"encoding/json"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/dbcommons"
	"im-server/commons/tools"
	"time"
)

var confCache *caches.LruCache

func init() {
	confCache = caches.NewLruCache("globalconf_cache", 1000, nil)
	confCache.AddTimeoutAfterCreate(time.Minute)
	confCache.SetValueCreator(func(key interface{}) interface{} {
		if key != nil {
			confDao := dbcommons.GlobalConfDao{}
			conf, err := confDao.FindByKey(key.(string))
			if err == nil {
				return conf.ConfValue
			}
		}
		return nil
	})
}

func GetGlobalConf(key string) string {
	val, exist := confCache.GetByCreator(key, nil)
	if exist && val != nil {
		return val.(string)
	}
	return ""
}

type AddressConf struct {
	Default   []string          `json:"default"`
	NodeConfs map[string]string `json:"confs"`
}

func GetOriginalNavAddress() *AddressConf {
	ret := &AddressConf{
		NodeConfs: map[string]string{},
	}
	nodes := bases.GetCluster().GetAllNodes()
	for _, node := range nodes {
		if val, ok := node.Exts[bases.NodeTag_Nav]; ok {
			navExt := bases.HttpNodeExt{}
			err := tools.JsonUnMarshal([]byte(val), &navExt)
			if err == nil {
				ret.NodeConfs[node.Name] = fmt.Sprintf("%s:%d", purifyNodeIp(node.Ip), navExt.Port)
			}
		}
	}
	return ret
}

func GetProxyNavAddress() *AddressConf {
	ret := &AddressConf{
		NodeConfs: map[string]string{},
	}
	adds := GetGlobalConf(string(dbcommons.GlobalConfKey_NaviAddress))
	if adds != "" {
		err := json.Unmarshal([]byte(adds), ret)
		if err == nil {
			return ret
		}
	}
	return nil
}

func GetNavAddress() *AddressConf {
	ret := GetProxyNavAddress()
	if ret != nil {
		return ret
	}
	//default value
	ret = GetOriginalNavAddress()
	return ret
}

func GetOriginalApiAddress() *AddressConf {
	ret := &AddressConf{
		NodeConfs: map[string]string{},
	}
	nodes := bases.GetCluster().GetAllNodes()
	for _, node := range nodes {
		if val, ok := node.Exts[bases.NodeTag_Api]; ok {
			apiExt := bases.HttpNodeExt{}
			err := tools.JsonUnMarshal([]byte(val), &apiExt)
			if err == nil {
				ret.NodeConfs[node.Name] = fmt.Sprintf("%s:%d", purifyNodeIp(node.Ip), apiExt.Port)
			}
		}
	}
	return ret
}

func GetProxyApiAddress() *AddressConf {
	ret := &AddressConf{}
	adds := GetGlobalConf(string(dbcommons.GlobalConfKey_ApiAddress))
	if adds != "" {
		err := json.Unmarshal([]byte(adds), ret)
		if err == nil {
			return ret
		}
	}
	return nil
}

func GetApiAddress() *AddressConf {
	ret := GetProxyApiAddress()
	if ret != nil {
		return ret
	}
	//default value
	ret = GetOriginalApiAddress()
	return ret
}

func GetOriginalConnectAddress() *AddressConf {
	ret := &AddressConf{
		NodeConfs: map[string]string{},
	}
	nodes := bases.GetCluster().GetAllNodes()
	for _, node := range nodes {
		if val, ok := node.Exts[bases.NodeTag_Connect]; ok {
			connectExt := bases.ConnectNodeExt{}
			err := tools.JsonUnMarshal([]byte(val), &connectExt)
			if err == nil {
				ret.NodeConfs[node.Name] = fmt.Sprintf("%s:%d", purifyNodeIp(node.Ip), connectExt.WsPort)
			}
		}
	}
	return ret
}

func GetProxyConnectAddress() *AddressConf {
	ret := &AddressConf{}
	adds := GetGlobalConf(string(dbcommons.GlobalConfKey_ConnectAddress))
	if adds != "" {
		err := json.Unmarshal([]byte(adds), ret)
		if err == nil {
			return ret
		}
	}
	return nil
}

func GetConnectAddress() *AddressConf {
	ret := GetProxyConnectAddress()
	if ret != nil {
		return ret
	}
	//default value
	ret = GetOriginalConnectAddress()
	return ret
}

func purifyNodeIp(ip string) string {
	return ip
}
