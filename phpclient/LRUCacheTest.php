<?php

use bybzmt\lrucache\LRUCache;

require_once __DIR__ .  "/Exception.php";
require_once __DIR__ . "/LRUCache.php";

class LRUCacheTest extends PHPUnit_Framework_TestCase
{
	private $_cache;

	public function setUp()
	{
		$server_url = "http://127.0.0.1:8080";

		$cache = new LRUCache($server_url);
		$cache->addAutoGroup("g1", 100, 60, 60, 300);
		$cache->addAutoGroup("g2", 100, 60, 60, 300);
		$cache->addAutoGroup("g3", 100, 60, 60, 300);

		$this->_cache = $cache;
	}

	public function testDelGroup()
	{
		$this->_cache->delGroup("g1");
		$this->_cache->delGroup("g2");
		$this->_cache->delGroup("g3");
	}

	public function testIncr()
	{
		$val = $this->_cache->incr("g1", "k1", 1);

		$this->assertTrue($val == 1);
	}

	public function testGet()
	{
		$val = $this->_cache->get("g1", "k1");

		$this->assertTrue($val == 1);
	}

	public function testSet()
	{
		$val = $this->_cache->set("g1", "k2", 2);

		$this->assertTrue($val);
	}

	public function testGetHot()
	{
		$val = $this->_cache->getHot("g1", "2");

		$this->assertEquals($val, array(
			array("name"=>"k2", "val"=>2),
			array('name'=>'k1', 'val'=>1),
		));
	}

	public function testDel()
	{
		$val = $this->_cache->del("g1", "k1");

		$this->assertTrue($val);
	}

	public function testDels()
	{
		$val = $this->_cache->dels(array(
			array("g1", "k1"),
			array("g1", "k2"),
		));

		$this->assertEquals($val, array(
			'err' => 0,
			'data' => array(
				array('err' => 0, 'data' => null),
				array('err' => 0, 'data' => null),
			)
		));
	}

	public function testSets()
	{
		$val = $this->_cache->sets(array(
			array("g1", "k1", 31),
			array("g1", "k2", 21),
			array("g1", "k3", 11),
		));

		$this->assertEquals($val, array(
			'err' => 0,
			'data' => array(
				array('err' => 0, 'data' => null),
				array('err' => 0, 'data' => null),
				array('err' => 0, 'data' => null),
			)
		));
	}

	public function testIncrs()
	{
		$val = $this->_cache->incrs(array(
			array("g1", "k1", 3),
			array("g1", "k2", 2),
			array("g1", "k3", 1),
		));

		$this->assertEquals($val, array(
			'err' => 0,
			'data' => array(
				array('err' => 0, 'data' => 34),
				array('err' => 0, 'data' => 23),
				array('err' => 0, 'data' => 12),
			)
		));
	}

	public function testGets()
	{
		$val = $this->_cache->gets(array(
			array("g1", "k1"),
			array("g1", "k2"),
			array("g1", "k3"),
		));

		$this->assertEquals($val, array(
			'err' => 0,
			'data' => array(
				array('err' => 0, 'data' => 34),
				array('err' => 0, 'data' => 23),
				array('err' => 0, 'data' => 12),
			)
		));
	}

	public function testGetHots()
	{
		$val = $this->_cache->getHots(array(
			array("g1", "2"),
			array("g2", "2"),
		));

		$this->assertEquals($val, array(
			'err' => 0,
			'data' => array(
				array('err' => 0, 'data' => array(
					array('name' => 'k1', 'val' => 34),
					array('name' => 'k2', 'val' => 23),
				)),
				array('err' => 2, 'data' => 'GroupNotExists'),
			)
		));
	}

}


