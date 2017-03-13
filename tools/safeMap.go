package tools

import (
	"sync"
	"xsbPro/log"
)

var (
	Cache *SafeMap
)

func InitCache() {
	Cache = NewSafeMap()
}

type NginxStatistics struct {
	URL     string
	Success int
	Fail    int
}

func newNginxStatistics(URL string, success, fail int) *NginxStatistics {
	return &NginxStatistics{
		URL:     URL,
		Success: success,
		Fail:    fail,
	}
}

type SafeMap struct {
	lock *sync.RWMutex
	bm   map[string]*NginxStatistics
}

// NewBeeMap return new safemap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[string]*NginxStatistics),
	}
}


//增加一个成功访问计数
func (m *SafeMap) AddOneSuccess(URL string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if value, ok := m.bm[URL]; !ok {
		m.bm[URL] = newNginxStatistics(URL, 1, 0)
	} else {
		value.Success += 1
	}
}

//增加一个失败访问计数
func (m *SafeMap) AddOneFail(URL string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if value, ok := m.bm[URL]; !ok {
		m.bm[URL] = newNginxStatistics(URL, 0, 1)
	} else {
		value.Fail += 1
	}
}

//清零所有api计数
func (m *SafeMap) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range m.bm {
		m.bm[k] = &NginxStatistics{
			URL:     v.URL,
			Success: 0,
			Fail:    0,
		}
	}
}

// Items returns all items in safemap.
func (m *SafeMap) Items() map[string]*NginxStatistics {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[string]*NginxStatistics)
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}



func SaveResultsToCache(infos []*NginxLogInfo) {
	for _, info := range (infos) {
		if len(info.URL) <= 0 {
			log.SysF("结构含有不正确的URL信息: %v", *info)
			continue
		}
		if info.Code == "200" || info.Code == "301" || info.Code == "302" || info.Code == "303" || info.Code == "304" {
			Cache.AddOneSuccess(info.URL)
		} else {
			Cache.AddOneFail(info.URL)
		}
	}
}
