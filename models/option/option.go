package option

import (
	"github.com/lisijie/goblog/models"
	"github.com/lisijie/goblog/util"
	"strconv"
)

//得到Options的所有配置信息
func GetOptions() map[string]string {
	rs, _ := util.Factory.Get("cache")
	cache := rs.(*util.LruCache) //类型动态转换/查询,把接口对象rs转换为*util.LruCache类型
	if !cache.IsExist("options") {
		var result []*models.Option
		new(models.Option).Query().All(&result)
		options := make(map[string]string)
		for _, v := range result {
			options[v.Name] = v.Value
		}
		cache.Put("options", options, 0)
	}
	v := cache.Get("options")
	return v.(map[string]string)
}

func FlushOptions() {
	rs, _ := util.Factory.Get("cache")
	cache := rs.(*util.LruCache)
	cache.Delete("options")
}

func Get(key string) string {
	options := GetOptions()
	if v, ok := options[key]; ok {
		return v
	}
	return ""
}

func GetInt(key string) int {
	v, _ := strconv.Atoi(Get(key))
	return v
}
