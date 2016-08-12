<?php
namespace bybzmt\lrucache;

/**
 * LRU算法缓存操作类
 */
class LRUCache
{
	private $_server_url;
	private $_groups = array();

	public function __construct($server_url)
	{
		$this->_server_url = $server_url;
	}

	/**
	 * 增加指定key的值, 如果key不存在默认为0
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function incr($group, $key, $num)
	{
		$res = $this->_doCall("/counter/incr", array(
			'group' => $group,
			'key' => $key,
			'val' => $num,
		));

		switch ($res['err']) {
		case 0:
			return $res['data'];
		case 1:
			throw new Exception("LRUCache Error:".$res['data']);
		case 2:
			$this->_autoGroup($group);
			return $this->incr($group, $key, $num);
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 增加指定key的值, 如果key不存在默认为0
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function incrs(array $groups, array $keys, array $nums)
	{
		$params= array();
		foreach ($groups as $idx => $group) {
			$params[] = array(
				'group' => $groups[$idx],
				'key' => $keys[$idx],
				'val' => $nums[$idx],
			);
		}

		$res = $this->_doCall("/counter/incr", array(
			'reqs' => json_encode($params),
		));

		foreach ($res['data'] as $idx => $_res) {
			switch ($res['err']) {
			case 2:
				$this->_autoGroup($groups[$idx]);
				$res[$idx] = $this->incr($groups[$idx], $keys[$idx], $nums[$idx]);
			}
		}

		return $res;
	}

	/**
	 * 将组里所有数据排序，取出最大的N个
	 */
	public function getHot($group, $len)
	{
		$res = $this->_doCall("/counter/hot", array(
			'group' => $group,
			'len' => $len,
		));

		switch ($res['err']) {
		case 0:
			return $res['data'];
		case 1:
			throw new Exception("LRUCache Error:".$res['data']);
		case 2:
			return array();
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 得到指定key的值
	 */
	public function get($group, $key)
	{
		$res = $this->_doCall("/cache/get", array(
			'group' => $group,
			'key' => $key,
		));

		switch ($res['err']) {
		case 0:
			return $res['data'];
		case 1:
			throw new Exception("LRUCache Error:".$res['data']);
		case 2:
		case 3:
			return null;
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 得到指定key的值
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function gets(array $groups, array $keys)
	{
		$params= array();
		foreach ($groups as $idx => $group) {
			$params[] = array(
				'group' => $group,
				'key' => $keys[$idx],
			);
		}

		$res = $this->_doCall("/cache/get", array(
			'reqs' => json_encode($params),
		));

		return $res;
	}

	/**
	 * 设定指定key的值
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function set($group, $key, $val)
	{
		$res = $this->_doCall("/cache/set", array(
			'group' => $group,
			'key' => $key,
			'val' => $val,
		));

		switch ($res['err']) {
		case 0:
			return true;
		case 1:
			throw new Exception("LRUCache Error:".$res['data']);
		case 2:
			$this->_autoGroup($group);
			return $this->set($group, $key);
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 设定指定key的值
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function sets(array $groups, array $keys, array $vals)
	{
		$params= array();
		foreach ($groups as $idx => $group) {
			$params[] = array(
				'group' => $group,
				'key' => $keys[$idx],
				'val' => $vals[$idx],
			);
		}

		$res = $this->_doCall("/cache/set", array(
			'reqs' => json_encode($params),
		));

		foreach ($res['data'] as $idx => $_res) {
			switch ($res['err']) {
			case 2:
				$this->_autoGroup($groups[$idx]);
				$res[$idx] = $this->set($groups[$idx], $keys[$idx], $vals[$idx]);
			}
		}

		return $res;
	}

	/**
	 * 删除指定key
	 */
	public function del($group, $key)
	{
		$res = $this->_doCall("/cache/del", array(
			'group' => $group,
			'key' => $key,
		));

		switch ($res['err']) {
		case 0:
			return true;
		case 1:
			throw new Exception("LRUCache Error:".$res['data']);
		case 2:
		case 3:
			return false;
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 删除指定key
	 * 如果group不存在，自动尝试通过addAutoGroup添加的数据自动创建
	 */
	public function dels(array $groups, array $keys)
	{
		$params= array();
		foreach ($groups as $idx => $group) {
			$params[] = array(
				'group' => $group,
				'key' => $keys[$idx],
			);
		}

		$res = $this->_doCall("/cache/set", array(
			'reqs' => json_encode($params),
		));

		return $res;
	}

	private function _autoGroup($group)
	{
		if (!isset($this->_groups[$group])) {
			throw new Exception("LRUCache Group Not Exists");
		}

		call_user_func_array(array($this, 'createGroup'), $this->_groups[$group]);
	}

	/**
	 * 添加需要自动创建的组
	 */
	public function addAutoGroup($group, $cap, $saveTick, $statusTick, $expire)
	{
		$this->_groups[$group] = array($group, $cap, $saveTick, $statusTick, $expire);
	}

	/**
	 * 手动创建组
	 */
	public function createGroup($group, $cap, $saveTick, $statusTick, $expire)
	{
		$res = $this->_doCall("/group/create", array(
			'group' => $group,
			'cap' => $cap,
			'saveTick' => $saveTick,
			'statusTick' => $statusTick,
			'expire' => $expire,
		));

		switch ($res['err']) {
		case 0:
			return true;
		case 1:
		case 2:
			return false;
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	/**
	 * 删除组
	 */
	public function delGroup($group)
	{
		$res = $this->_doCall("/group/del", array(
			'group' => $group,
		));

		switch ($res['err']) {
		case 0:
			return true;
		case 1:
		case 2:
			return false;
		default:
			throw new Exception("LRUCache Unknown Error.");
		}
	}

	//进行网络请求
	private function _doCall($req, array $params)
	{
		$opts = array('http' =>
			array(
				'method'  => 'POST',
				'header'  => 'Content-type: application/x-www-form-urlencoded',
				'content' => http_build_query($params),
			)
		);

		$context = stream_context_create($opts);

		$result = file_get_contents($this->_server_url . $req, false, $context);
		if (!$result) {
			throw new Exception("LRUCache Server Error.");
		}

		return json_decode($result, true);
	}

}
