<?php
namespace bybzmt\lrucache;

require "./LRUCache.php";
require "./Exception.php";

$server_url = "http://127.0.0.1:8080";

$cache = new LRUCache($server_url);
$cache->addAutoGroup("test1", 100, 60, 60, 300);

$val = $cache->incr("test1", "k1", 1);
var_dump($val);

$val = $cache->get("test1", "k1");
var_dump($val);

$val = $cache->getHot("test1", "10");
var_dump($val);

$val = $cache->set("test1", "k2", 1);
var_dump($val);

$val = $cache->del("test1", "k1");
var_dump($val);

$val = $cache->delGroup("test1");
var_dump($val);

echo "ok\n";
