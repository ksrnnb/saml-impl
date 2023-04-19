package kvs

// TODO: implement TTL
var store = map[string]any{}

func Set(key string, value any) {
	store[key] = value
}

func Get(key string) any {
	return store[key]
}

func Delete(key string) {
	delete(store, key)
}

func Exists(key string) bool {
	_, ok := store[key]
	return ok
}
