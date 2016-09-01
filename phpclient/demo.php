<?php
namespace bybzmt\lrucache;

require "./LRUCache.php";
require "./Exception.php";

$server_url = "http://127.0.0.1:8080";

$cache = new LRUCache($server_url);

//下面需要使用到的组需要先添加进来
$cache->addAutoGroup("g1", 100, 60, 60, 300);
$cache->addAutoGroup("g2", 100, 60, 60, 300);
$cache->addAutoGroup("g3", 100, 60, 60, 300);

//增加数值
$val = $cache->incr("g1", "k1", 1);
var_dump($val);
//批量增加
$val = $cache->incrs(array(
	array("g1", "k1",1),
	array("g2", "k1",1),
	array("g3", "k1",1),
));
var_dump($val);

//到一个key
$val = $cache->get("g1", "k1");
var_dump($val);
//批量取key
$val = $cache->gets(array(
	array("g2", "k1"),
	array("g3", "k1"),
));
var_dump($val);

//得到数值最大的N条
$val = $cache->getHot("g1", "2");
var_dump($val);
$val = $cache->getHots(array(
	array("g2", "2"),
	array("g3", "2"),
));
var_dump($val);

//设置key
$val = $cache->set("g1", "k2", 1);
var_dump($val);
//批量设置key
$val = $cache->sets(array(
	array("g2", "k1",1),
	array("g3", "k1",1),
));
var_dump($val);

//删降key
$val = $cache->del("g1", "k1");
var_dump($val);
//批量删除key
$val = $cache->dels(array(
	array("g2", "k1"),
	array("g3", "k1"),
));
var_dump($val);

//删除组
$val = $cache->delGroup("g1");
var_dump($val);
//批量删除组
$val = $cache->delGroups(array(
	array("g2"),
	array("g3"),
));
var_dump($val);

echo "ok\n";
