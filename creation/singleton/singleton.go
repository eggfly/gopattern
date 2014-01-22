package singleton

import (
	//"fmt"
	"sync"
)

var S interface{} = "origin value"
var I = 0

var m map[interface{}]interface{} = make(map[interface{}]interface{})
var mRWLock sync.RWMutex

// func init() {
// 	// singleton only initialized by init() function is thread safe
// 	S = "new value"
// 	fmt.Println("init value in init()")
// }

func getMap(key interface{}) interface{} {
	mRWLock.RLock()
	v, _ := m[key]
	mRWLock.RUnlock()
	return v
}
func setMap(key interface{}, value interface{}) {
	mRWLock.Lock()
	m[key] = value
	mRWLock.Unlock()
}
