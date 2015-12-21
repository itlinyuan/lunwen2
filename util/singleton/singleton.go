package singleton

import (
	"sync"
)

type SingletonInitFunc func() (interface{}, error)

//单例模式
// Example use:
// var configSelectorSingleton = NewSingleton(init)
// func configSelector() (configSelector, error) {
//     s, err := configSelectorSingleton.Get()
//     if err != nil {
//         return nil, err
//     }
//     return s.(configSelector), nil
// }
type Singleton interface {
	// 返回封装好的单例对象
	Get() (interface{}, error)
}

// 封装构造函数（调用创建一个新的单例与给定的init函数实例化。）
// init is not called until the first invocation of Get().  If init errors, it will be called again on the next invocation of Get().
func NewSingleton(init SingletonInitFunc) Singleton {
	return &singletonImpl{init: init}
}

type singletonImpl struct {
	sync.Mutex

	// 单例对象
	data interface{}
	// 构造函数
	init SingletonInitFunc
	// 初始化没错误的话为true
	initialized bool
}

func (s *singletonImpl) Get() (interface{}, error) {
	// Don't lock in the common case
	if s.initialized {
		return s.data, nil
	}

	s.Lock()
	defer s.Unlock()

	if s.initialized {
		return s.data, nil
	}

	var err error
	s.data, err = s.init()
	if err != nil {
		return nil, err
	}

	s.initialized = true
	return s.data, nil
}
