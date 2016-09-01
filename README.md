LRUCache
=========

这是一个使用LRU算法的cache服务 (LRU是指最后使用缓存)

cache中分不同的分组，每个分组有自己的容量和LRU维护

每次对指定key的数字进行写操作时，就会把这个key添加到
链表的开头。如果保存的key的个数达到上限，则会把链表尾的key给删除。

访问的方式为http请求 .
参数可以接受GET或POST方式
反回JSON `{err:0, data:data}`

错误类型:
* `0` 成功
* `1` 错误, 非特定错误
* `2` 分组不存在
* `3` key不存在

##各请求url介绍

###累加
`/counter/incr`
给指定key增加N，如果key不存在默认为0

参数:
* `group` 分组名
* `key` 键名
* `val` 数值

###取出数值最大的的N个key
`/counter/hot`

参数:
* `group` 分组名
* `len` 取出的key的个数

###设置指定的key
/cache/set

参数:
* `group` 分组名
* `key` 键名
* `val` 数值

###取得指定的key的值
/cache/get

参数:
* `group` 分组名
* `key` 键名

###删除指定的key
/cache/del

参数:
* `group` 分组名
* `key` 键名

###创建一个分组
/group/create

参数:
* `group` 分组名
* `cap`   分组容量
* `expire` 过期时间
* `saveTick`   定时保存同期，0 不定时保存
* `statusTick` 定时统状态同期，0 不统计

###删除一个分组
/group/del

参数
* `group` 分组名

###同时请求多个
`/multiple/` 这个地址再接上面的， 如 `/multiple/cache/get` 
参数是

	$data = array();
	//参数与上面的一至
	$data[] = array('group'=>'group1', 'key'=>'key1');
	$data[] = array('group'=>'group2', 'key'=>'key2');
	reqs = json_encode($data)

	程序会返回一个与请求数量一样至的一个数组，里面是每个请求的反回

