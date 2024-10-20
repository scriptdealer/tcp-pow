package services

import (
	"container/list"
	"log/slog"
	"sync"
	"time"
)

const ConnLimitReached int = 0

type RequestRecord struct {
	IP      string
	Count   int
	LastReq time.Time
}

type ConnectionPool struct {
	access sync.Mutex
	log    *slog.Logger
	list   *list.List
	elMap  map[string]*list.Element
	window time.Duration
	size   int
	maxCnn int // max connections per IP address
}

func NewConnectionPool(logger *slog.Logger, size, maxC int, window time.Duration) *ConnectionPool {
	return &ConnectionPool{
		log:    logger,
		list:   list.New(),
		elMap:  make(map[string]*list.Element),
		size:   size,
		maxCnn: maxC,
		window: window,
	}
}

func (cp *ConnectionPool) Reset() {
	cp.access.Lock()
	defer cp.access.Unlock()

	cp.list = list.New()
	cp.elMap = make(map[string]*list.Element)
}

func (cp *ConnectionPool) CountRequests(ip string) int {
	cp.access.Lock()
	defer cp.access.Unlock()

	now := time.Now()
	cp.removeOldRecords(now)

	if elem, exists := cp.elMap[ip]; exists {
		record := elem.Value.(*RequestRecord)
		record.Count++
		record.LastReq = now
		cp.list.MoveToFront(elem)
		if cp.maxCnn < record.Count {
			return cp.maxCnn
		}
		return record.Count
	}

	if len(cp.elMap) >= cp.size {
		return ConnLimitReached
	}

	newRecord := &RequestRecord{IP: ip, Count: 1, LastReq: now}
	elem := cp.list.PushFront(newRecord)
	cp.elMap[ip] = elem
	return 1
}

func (cp *ConnectionPool) removeOldRecords(now time.Time) {
	for e := cp.list.Back(); e != nil; {
		record := e.Value.(*RequestRecord)
		if now.Sub(record.LastReq) > cp.window {
			outdated := e
			e = e.Prev()
			cp.list.Remove(outdated)
			delete(cp.elMap, record.IP)
		} else {
			break
		}
	}
}
