package service

import (
	"github.com/ksrnnb/saml-impl/kvs"
)

const requestIDsKey = "requestIds"

func ListRequestIDs() []string {
	ids := kvs.Get(requestIDsKey)
	if ids == nil {
		return []string{}
	}
	return ids.([]string)
}

func AddRequestID(id string) string {
	ids := ListRequestIDs()
	if ids == nil {
		ids = []string{}
	}
	ids = append(ids, id)
	kvs.Set(requestIDsKey, ids)
	return id
}

func DeleteRequestID(id string) {
	ids := ListRequestIDs()
	for i, v := range ids {
		if v == id {
			ids[i] = ids[len(ids)-1]
			ids = ids[:len(ids)-1]
		}
	}
	kvs.Set(requestIDsKey, ids)
}

func ExistsRequestID(rid string) bool {
	ids := ListRequestIDs()
	for _, id := range ids {
		if id == rid {
			return true
		}
	}
	return false
}
