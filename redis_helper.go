/*
*	Description : redis 相关方法  TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Rds Redis客户端
type Rds struct {
	SSHUser       string
	SSHPassword   string
	SSHAddr       string
	RedisHost     string
	RedisPost     string
	RedisPassword string

	// redis DB
	RedisDB int

	// 单个连接
	Conn redis.Conn

	//	最大闲置数，用于redis连接池
	RedisMaxIdle int

	//	最大连接数
	RedisMaxActive int

	//	单条连接Timeout
	RedisIdleTimeoutSec int

	// 连接池
	Pool *redis.Pool
}

// SSHConnInfo ssh连接通道
type SSHConnInfo struct {
	SSHUser     string
	SSHPassword string
	SSHAddr     string
}

// NewSSHInfo 新建ssh连接通道
func NewSSHInfo(addr, user, password string) *SSHConnInfo {
	return &SSHConnInfo{
		SSHUser:     user,
		SSHPassword: password,
		SSHAddr:     addr,
	}
}

// NewRedis 新建Redis客户端对象
func NewRedis(host, port, password string, db int, vs ...any) (*Rds, error) {
	var sshConnInfo SSHConnInfo
	for _, v := range vs {
		switch vv := v.(type) {
		case *SSHConnInfo:
			sshConnInfo = *vv
		case SSHConnInfo:
			sshConnInfo = vv
		}
	}
	rds := &Rds{
		SSHUser:       sshConnInfo.SSHUser,
		SSHPassword:   sshConnInfo.SSHPassword,
		SSHAddr:       sshConnInfo.SSHAddr,
		RedisHost:     host,
		RedisPost:     port,
		RedisPassword: password,
		RedisDB:       db,
	}
	err := rds.RedisConn()
	return rds, err
}

// NewRedisPool 新建Redis连接池对象
func NewRedisPool(host, port, password string, db, maxIdle, maxActive, idleTimeoutSec int, vs ...any) *Rds {
	var sshConnInfo SSHConnInfo
	for _, v := range vs {
		switch vv := v.(type) {
		case *SSHConnInfo:
			sshConnInfo = *vv
		case SSHConnInfo:
			sshConnInfo = vv
		}
	}
	return &Rds{
		SSHUser:             sshConnInfo.SSHUser,
		SSHPassword:         sshConnInfo.SSHPassword,
		SSHAddr:             sshConnInfo.SSHAddr,
		RedisHost:           host,
		RedisPost:           port,
		RedisPassword:       password,
		RedisDB:             db,
		RedisMaxIdle:        maxIdle,
		RedisMaxActive:      maxActive,
		RedisIdleTimeoutSec: idleTimeoutSec,
	}
}

// RedisConn redis连接
func (r *Rds) RedisConn() (err error) {
	host := fmt.Sprintf("%s:%s", r.RedisHost, r.RedisPost)
	if r.SSHPassword != "" && r.SSHUser != "" && r.SSHAddr != "" {
		if sshClient, err := SSHClient(r.SSHUser, r.SSHPassword, r.SSHAddr); err == nil {
			var conn net.Conn
			conn, err = sshClient.Dial("tcp", host)
			r.Conn = redis.NewConn(conn, -1, -1)
		}
	} else {
		r.Conn, err = redis.Dial("tcp", host)
	}
	if err != nil {
		return
	}
	if r.Conn == nil {
		err = errors.New("redis conn is null")
		return
	}
	if r.RedisPassword != "" {
		if _, authErr := r.Conn.Do("AUTH", r.RedisPassword); authErr != nil {
			err = fmt.Errorf("redis auth password error: %s", authErr)
			return
		}
	}
	if r.RedisDB < 1 {
		r.RedisDB = 0
	}
	_, err = r.Conn.Do("select", fmt.Sprintf("%d", r.RedisDB))
	return
}

// RedisPool 连接池连接
// 返回redis连接池  *redis.Pool.Get() 获取redis连接
func (r *Rds) RedisPool() error {
	host := fmt.Sprintf("%s:%s", r.RedisHost, r.RedisPost)

	r.Pool = &redis.Pool{
		MaxIdle:     r.RedisMaxIdle,
		MaxActive:   r.RedisMaxActive,
		IdleTimeout: time.Duration(r.RedisIdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			var (
				c   redis.Conn
				err error
			)
			if r.SSHPassword != "" && r.SSHUser != "" && r.SSHAddr != "" {
				//ssh Client
				sshClient, err := SSHClient(r.SSHUser, r.SSHPassword, r.SSHAddr)
				if err != nil {
					return nil, err
				}
				conn, err := sshClient.Dial("tcp", host)
				if err != nil {
					return nil, err
				}
				c = redis.NewConn(conn, 60, 60)
				//if  sshClient != nil {
				//	var conn net.Conn
				//	conn, err = sshClient.Dial("tcp", host)
				//	c = redis.NewConn(conn, -1, -1)
				//}
				//if err != nil{
				//	return nil, err
				//}
			} else {
				c, err = redis.Dial("tcp", host)
				if err != nil {
					return nil, fmt.Errorf("redis connection error: %s", err)
				}
			}
			if c == nil {
				return nil, fmt.Errorf("redis connection is null")
			}
			//验证redis密码
			if r.RedisPassword != "" {
				if _, authErr := c.Do("AUTH", r.RedisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}
			_, err = c.Do("select", fmt.Sprintf("%d", r.RedisDB))
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
	return nil
}

// GetConn 获取redis连接
func (r *Rds) GetConn() redis.Conn {
	if r.Conn != nil {
		return r.Conn
	}
	rc := r.Pool.Get()
	if rc != nil {
		return rc
	}
	return nil
}

// SelectDB 切换redis db
func (r *Rds) SelectDB(dbNumber int) error {
	rc := r.GetConn()
	if rc == nil {
		return errors.New("redis conn is nil")
	}
	_, err := rc.Do("select", fmt.Sprintf("%d", dbNumber))
	return err
}

// RedisDELKeys Del key
// 使用常见： 并发删除大量key
func RedisDELKeys(rds *Rds, keys string, jobNumber int) {
	CPUMax()
	rds.RedisMaxActive = rds.RedisMaxActive + jobNumber*2
	rds.RedisMaxIdle = rds.RedisMaxIdle + jobNumber*2

	_ = rds.RedisPool()
	conn := rds.Pool.Get()
	queue := NewQueue()
	res, err := redis.Strings(conn.Do("keys", keys))
	if err != nil {
		Error(err)
	}
	_ = conn.Close()

	for _, v := range res {
		_ = queue.Add(&Task{Url: v})
	}
	allNumber := queue.Size()

	var wg sync.WaitGroup
	for job := 0; job < jobNumber; job++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Info("启动第", i, "个任务")

			for {
				if queue.IsEmpty() || queue.Size() < 2 {
					break
				}

				task := queue.Poll()
				Info("第", i, "个任务取的值： ", task.Url)
				c := rds.Pool.Get()
				s, err := redis.Int64(c.Do("DEL", task.Url))
				if err != nil || s == 0 {
					Info("redis command:  err : ", err)
				} else {
					Info("删除成功 ！！！")
				}
				_ = c.Close()
				Info(fmt.Sprintf("[进度] %d/%d  %f %%", allNumber-queue.Size(),
					allNumber, (float64(allNumber-queue.Size())/float64(allNumber))*100))
			}
			Info("第", i, "个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// 使用List实现消息队列

// MqProducer Redis消息队列生产方
func (r *Rds) MqProducer(mqName string, data any) error {
	args := redis.Args{}.Add(mqName)
	args = args.Add(data)
	_, err := r.GetConn().Do("LPUSH", args...)
	return err
}

// MqConsumer Redis消息队列消费方
func (r *Rds) MqConsumer(mqName string) (reply any, err error) {
	if r.MqLen(mqName) < 1 {
		return nil, fmt.Errorf("data len is 0")
	}
	return r.GetConn().Do("RPOP", mqName)
}

// MqLen Redis消息队列消息数量
func (r *Rds) MqLen(mqName string) int64 {
	number, err := redis.Int64(r.GetConn().Do("LLEN", mqName))
	if err != nil {
		number = 0
	}
	return number
}

func (r *Rds) ToString(reply any, err error) (string, error) {
	return redis.String(reply, err)
}

func (r *Rds) ToInt(reply any, err error) (int, error) {
	return redis.Int(reply, err)
}

func (r *Rds) ToInt64(reply any, err error) (int64, error) {
	return redis.Int64(reply, err)
}

func (r *Rds) ToBool(reply any, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func (r *Rds) ToBytes(reply any, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

func (r *Rds) ToByteSlices(reply any, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

func (r *Rds) ToFloat64(reply any, err error) (float64, error) {
	return redis.Float64(reply, err)
}

func (r *Rds) ToFloat64s(reply any, err error) ([]float64, error) {
	return redis.Float64s(reply, err)
}

func (r *Rds) ToInt64Map(reply any, err error) (map[string]int64, error) {
	return redis.Int64Map(reply, err)
}

func (r *Rds) ToInt64s(reply any, err error) ([]int64, error) {
	return redis.Int64s(reply, err)
}

func (r *Rds) ToIntMap(reply any, err error) (map[string]int, error) {
	return redis.IntMap(reply, err)
}

func (r *Rds) ToInts(reply any, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

func (r *Rds) ToStringMap(reply any, err error) (map[string]string, error) {
	return redis.StringMap(reply, err)
}

func (r *Rds) ToStrings(reply any, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

// Keys ===============================================================================================================

// GetAllKeys 获取所有的key
func (r *Rds) GetAllKeys(match string) (ksyList map[string]int) {
	//初始化拆分值
	matchSplit := match
	//match :匹配值，没有则匹配所有 *
	if match == "" {
		match = "*"
	} else {
		match = fmt.Sprintf("*%s*", match)
	}
	//cursor :初始游标为0
	cursor := "0"
	ksyList = make(map[string]int)
	ksyList, cursor = r.addGetKey(ksyList, cursor, match, matchSplit)
	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if cursor != "0" {
		for {
			ksyList, cursor = r.addGetKey(ksyList, cursor, match, matchSplit)
			if cursor == "0" {
				break
			}
		}
	}
	return
}

// addGetKey 内部方法
// 针对分组的key进行分组合并处理
func (r *Rds) addGetKey(ksyList map[string]int, cursor, match, matchSplit string) (map[string]int, string) {
	countNumber := "10000"
	res, err := redis.Values(r.Conn.Do("scan", cursor, "MATCH", match, "COUNT", countNumber))
	InfoTimes(3, "[Redis Log] execute :", "scan ", cursor, " MATCH ", match, " COUNT ", countNumber)
	if err != nil {
		Error("GET error", err.Error())
	}
	//获取	match 含有多少:
	cfNumber := strings.Count(match, ":")
	//获取新的游标
	newCursor := string(res[0].([]byte))
	allKey := res[1]
	allKeyData := allKey.([]any)
	for _, v := range allKeyData {
		keyData := string(v.([]byte))
		//没有:的key 则不集合
		if strings.Count(keyData, ":") == cfNumber || keyData == match {
			ksyList[keyData] = 0
			continue
		}
		//有:需要集合
		keyDataNew, _ := fenGeYinHaoOne(keyData, matchSplit)
		ksyList[keyDataNew] = ksyList[keyDataNew] + 1
	}
	return ksyList, newCursor
}

// fenGeYinHaoOne 对查询出来的key进行拆分，集合，分组处理
func fenGeYinHaoOne(str string, matchSplit string) (string, int) {
	likeKey := ""
	if matchSplit != "" {
		likeKey = fmt.Sprintf("%s", matchSplit)
	}
	str = strings.Replace(str, likeKey, "", 1)
	fg := strings.Split(str, ":")
	if len(fg) > 0 {
		return fmt.Sprintf("%s%s", likeKey, fg[0]), len(fg)
	}
	return "", len(fg)
}

// SearchKeys  搜索key
func (r *Rds) SearchKeys(match string) (ksyList map[string]int) {
	ksyList = make(map[string]int)
	if match == "" {
		return
	} else {
		match = fmt.Sprintf("*%s*", match)
	}
	cursor := "0"
	ksyList = make(map[string]int)
	ksyList, cursor = r.addSearchKey(ksyList, cursor, match)
	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if cursor != "0" {
		for {
			ksyList, cursor = r.addSearchKey(ksyList, cursor, match)
			if cursor == "0" {
				break
			}
		}
	}
	return
}

// addGetKey 内部方法获取key
func (r *Rds) addSearchKey(ksyList map[string]int, cursor, match string) (map[string]int, string) {
	countNumber := "10000"
	res, err := redis.Values(r.Conn.Do("scan", cursor, "MATCH", match, "COUNT", countNumber))
	InfoTimes(3, "[Redis Log] execute :", "scan ", cursor, " MATCH ", match, " COUNT ", countNumber)
	if err != nil {
		Error("GET error", err.Error())
	}
	//获取新的游标
	newCursor := string(res[0].([]byte))
	allKey := res[1]
	allKeyData := allKey.([]any)
	for _, v := range allKeyData {
		keyData := string(v.([]byte))
		ksyList[keyData] = 0
	}
	return ksyList, newCursor
}

// Type 获取key的类型
func (r *Rds) Type(key string) string {
	InfoFTimes(3, "[Redis Log] execute : TYPE %s", key)
	res, err := redis.String(r.Conn.Do("TYPE", key))
	if err != nil {
		ErrorTimes(3, "GET error", err.Error())
	}
	return res
}

// Ttl 获取key的过期时间
func (r *Rds) Ttl(key string) int64 {
	InfoFTimes(3, "[Redis Log] execute : TTL %s", key)
	res, err := redis.Int64(r.Conn.Do("TTL", key))
	if err != nil {
		ErrorTimes(3, "GET error", err.Error())
	}
	return res
}

// DUMP 检查给定 key 是否存在。
func (r *Rds) DUMP(key string) bool {
	InfoFTimes(3, "[Redis Log] execute : DUMP %s", key)
	data, err := redis.String(r.Conn.Do("DUMP", key))
	if err != nil || data == "0" {
		ErrorTimes(3, "GET error", err.Error())
		return false
	}
	return true
}

// Rename 修改key名称
func (r *Rds) Rename(name, newName string) bool {
	arg := redis.Args{}.Add(name).Add(newName)
	InfoFTimes(3, "[Redis Log] execute : RENAME %s %v", name, newName)
	_, err := r.Conn.Do("RENAME", arg...)
	if err != nil {
		ErrorTimes(3, "GET error", err.Error())
		return false
	}
	return true
}

// Expire 更新key ttl
func (r *Rds) Expire(key string, ttl int64) bool {
	arg := redis.Args{}.Add(key).Add(ttl)
	InfoFTimes(3, "[Redis Log] execute : EXPIRE %s %v", key, ttl)
	_, err := r.Conn.Do("EXPIRE", arg...)
	if err != nil {
		ErrorTimes(3, err.Error())
		return false
	}
	return true
}

// ExpireAt 指定key多久过期 接收的是unix时间戳
func (r *Rds) ExpireAt(key string, date int64) bool {
	arg := redis.Args{}.Add(key).Add(date)
	InfoFTimes(3, "[Redis Log] execute : EXPIREAT %s %v", key, date)
	_, err := r.Conn.Do("EXPIREAT", arg...)
	if err != nil {
		ErrorTimes(3, err.Error())
		return false
	}
	return true
}

// DelKey 删除key
func (r *Rds) DelKey(key string) bool {
	InfoFTimes(3, "[Redis Log] execute : DEL %s", key)
	_, err := r.Conn.Do("DEL", key)
	if err != nil {
		ErrorTimes(3, err.Error())
		return false
	}
	return true
}

var RdsNotConnError = fmt.Errorf("未连接redis")

// String =============================================================================================================

// Get GET 获取String value
func (r *Rds) Get(key string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : GET %s", key)
	return redis.String(r.Conn.Do("GET", key))
}

// Set SET新建String
func (r *Rds) Set(key string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : SET %s %v", key, value)
	_, err := r.Conn.Do("SET", arg...)
	return err
}

// SetEx SETEX 新建String 含有时间
func (r *Rds) SetEx(key string, ttl int64, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(ttl).Add(value)
	InfoFTimes(3, "[Redis Log] execute : SETEX  %s %v %v", key, ttl, value)
	_, err := r.Conn.Do("SETEX", arg...)
	return err
}

// PSetEx PSETEX key milliseconds value
// 这个命令和 SETEX 命令相似，但它以毫秒为单位设置 key 的生存时间，而不是像 SETEX 命令那样，以秒为单位。
func (r *Rds) PSetEx(key string, ttl int64, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(ttl).Add(value)
	InfoFTimes(3, "[Redis Log] execute : PSETEX %s %v %v", key, ttl, value)
	_, err := r.Conn.Do("PSETEX", arg...)
	return err
}

// SetNx key value
// 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
func (r *Rds) SetNx(key string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : SETNX %s %v", key, value)
	arg := redis.Args{}.Add(key).Add(value)
	_, err := r.Conn.Do("SETNX", arg...)
	return err
}

// SetRange SETRANGE key offset value
// 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (r *Rds) SetRange(key string, offset int64, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(offset).Add(value)
	InfoFTimes(3, "[Redis Log] execute : SETRANGE %s %v %v", key, offset, value)
	_, err := r.Conn.Do("SETRANGE", arg...)
	return err
}

// Append APPEND key value
// 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (r *Rds) Append(key string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : APPEND %s %v", key, value)
	_, err := redis.String(r.Conn.Do("APPEND", arg...))
	return err
}

// SetBit SETBIT key offset value
// 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// value : 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 注意 offset 不能太大，越大key越大
func (r *Rds) SetBit(key string, offset, value int64) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(offset).Add(value)
	InfoFTimes(3, "[Redis Log] execute : SETBIT %s %d %d", key, offset, value)
	_, err := r.Conn.Do("SETBIT", arg...)
	return err
}

// BitCount BITCOUNT key [start] [end]
// 计算给定字符串中，被设置为 1 的比特位的数量。
func (r *Rds) BitCount(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : BITCOUNT %s", key)
	return redis.Int64(r.Conn.Do("BITCOUNT", key))
}

// GetBit GETBIT key offset
// 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
func (r *Rds) GetBit(key string, offset int64) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(offset)
	InfoFTimes(3, "[Redis Log] execute : GETBIT %s %d", key, offset)
	return redis.Int64(r.Conn.Do("GETBIT", arg...))
}

// TODO StringBITOP BITOP operation destkey key [key ...]
// 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
// BITOP AND destkey key [key ...] ，对一个或多个 key 求逻辑并，并将结果保存到 destkey 。
// BITOP OR destkey key [key ...] ，对一个或多个 key 求逻辑或，并将结果保存到 destkey 。
// BITOP XOR destkey key [key ...] ，对一个或多个 key 求逻辑异或，并将结果保存到 destkey 。
// BITOP NOT destkey key ，对给定 key 求逻辑非，并将结果保存到 destkey 。

// Decr key
// 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
func (r *Rds) Decr(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : DECR %s", key)
	return redis.Int64(r.Conn.Do("DECR", key))
}

// DecrBy DECRBY key decrement
// 将 key 所储存的值减去减量 decrement 。
func (r *Rds) DecrBy(key, decrement string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(decrement)
	InfoFTimes(3, "[Redis Log] execute : DECRBY %s %v", key, decrement)
	return redis.Int64(r.Conn.Do("DECRBY", arg...))
}

// GetRange GETRANGE key start end
// 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
func (r *Rds) GetRange(key string, start, end int64) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(end)
	InfoFTimes(3, "[Redis Log] execute : GETRANGE %s %v %v", key, start, end)
	return redis.String(r.Conn.Do("GETRANGE", arg...))
}

// GetSet GETSET key value
// 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
func (r *Rds) GetSet(key string, value any) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : GETSET %s %v", key, value)
	return redis.String(r.Conn.Do("GETSET", arg...))
}

// Incr INCR key
// 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
func (r *Rds) Incr(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : INCR %s", key)
	return redis.Int64(r.Conn.Do("INCR", key))
}

// IncrBy INCRBY key increment
// 将 key 所储存的值加上增量 increment 。
func (r *Rds) IncrBy(key, increment string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(increment)
	InfoFTimes(3, "[Redis Log] execute : INCRBY %s %v", key, increment)
	return redis.Int64(r.Conn.Do("INCRBY", arg...))
}

// IncrByFloat INCRBYFLOAT key increment
// 为 key 中所储存的值加上浮点数增量 increment 。
func (r *Rds) IncrByFloat(key, increment float64) (float64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(increment)
	InfoFTimes(3, "[Redis Log] execute : INCRBYFLOAT %s %v", key, increment)
	return redis.Float64(r.Conn.Do("INCRBYFLOAT", arg...))
}

// MGet MGET key [key ...]
// 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
func (r *Rds) MGet(key []any) ([]string, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}
	for _, value := range key {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : MGET %s %s", key, strings.Join(Any2Strings(key), " "))
	return redis.Strings(r.Conn.Do("MGET", args...))
}

// MSet MSET key value [key value ...]
// 同时设置一个或多个 key-value 对。
// 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，
// 请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。
// MSET 是一个原子性(atomic)操作，所有给定 key 都会在同一时间内被设置，某些给定 key 被更新而另一些给定 key 没有改变的情况，不可能发生。
func (r *Rds) MSet(values []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}
	for _, value := range values {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : MSET %s", strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("MSET", args...)
	return err
}

// MSetNx MSETNX key value [key value ...]
// 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，
// 所有字段要么全被设置，要么全不被设置。
func (r *Rds) MSetNx(values []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}
	for _, value := range values {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : MSETNX %s", strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("MSETNX", args...)
	return err
}

// TODO StringSTRLEN STRLEN key
// 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。

// List ===============================================================================================================

// LRange LRANGE 获取List value
func (r *Rds) LRange(key string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : LRANGE %s 0 -1", key)
	return redis.Values(r.Conn.Do("LRANGE", key, 0, -1))
}

// LRangeST LRANGE key start stop
// 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
func (r *Rds) LRangeST(key string, start, stop int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop)
	InfoFTimes(3, "[Redis Log] execute : LRANGE %s %v %v", key, start, stop)
	return redis.Values(r.Conn.Do("LRANGE", arg...))
}

// LPush LPUSH 新创建list 将一个或多个值 value 插入到列表 key 的表头
func (r *Rds) LPush(key string, values []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : LPUSH %s %s", key, strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("LPUSH", args...)
	return err
}

// RPush RPUSH key value [value ...]
// 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行
// RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
// 新创建List  将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (r *Rds) RPush(key string, values []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : RPUSH %s %s", key, strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("RPUSH", args...)
	return err
}

// TODO ListBLPOP BLPOP key [key ...] timeout
// BLPOP 是列表的阻塞式(blocking)弹出原语。
// 它是 LPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPOP 命令阻塞，直到等待超时或发现可弹出元素为止。

// TODO ListBRPOP BRPOP key [key ...] timeout
// BRPOP 是列表的阻塞式(blocking)弹出原语。
// 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。

// TODO ListBRPOPLPUSH BRPOPLPUSH source destination timeout
// BRPOPLPUSH 是 RPOPLPUSH 的阻塞版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH 一样。
// 当列表 source 为空时， BRPOPLPUSH 命令将阻塞连接，直到等待超时，或有另一个客户端对 source 执行 LPUSH 或 RPUSH 命令为止。

// LIndex LINDEX key index
// 返回列表 key 中，下标为 index 的元素。
func (r *Rds) LIndex(key string, index int64) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(index)
	InfoFTimes(3, "[Redis Log] execute : LINDEX %s %v", key, index)
	return redis.String(r.Conn.Do("LINDEX", arg...))
}

// LInsert LINSERT key BEFORE|AFTER pivot value
// 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
// 当 pivot 不存在于列表 key 时，不执行任何操作。
// 当 key 不存在时， key 被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// direction : 方向 bool true:BEFORE(前)    false: AFTER(后)
func (r *Rds) LInsert(direction bool, key, pivot, value string) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	directionStr := "AFTER"
	if direction {
		directionStr = "BEFORE"
	}
	arg := redis.Args{}.Add(key).Add(directionStr).Add(pivot).Add(value)
	InfoFTimes(3, "[Redis Log] execute : LINSERT %s %v %v %v", key, directionStr, pivot, value)
	_, err := r.Conn.Do("LINSERT", arg...)
	return err
}

// LLen LLEN key
// 返回列表 key 的长度。
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
func (r *Rds) LLen(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : LLEN %s", key)
	return redis.Int64(r.Conn.Do("LLEN", key))
}

// ListLPOP LPOP key
// 移除并返回列表 key 的头元素。
func (r *Rds) ListLPOP(key string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : LPOP %s", key)
	return redis.String(r.Conn.Do("LPOP", key))
}

// LPusHx LPUSHX key value
// 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
// 和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
func (r *Rds) LPusHx(key string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : LPUSHX %s %v", key, value)
	_, err := r.Conn.Do("LPUSHX", arg...)
	return err
}

// LRem LREM key count value
// 根据参数 count 的值，移除列表中与参数 value 相等的元素。
// count 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
// count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
// count = 0 : 移除表中所有与 value 相等的值。
func (r *Rds) LRem(key string, count int64, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(count).Add(value)
	InfoFTimes(3, "[Redis Log] execute : LREM %s %v %v", key, count, value)
	_, err := r.Conn.Do("LREM", arg...)
	return err
}

// LSet LSET key index value
// 将列表 key 下标为 index 的元素的值设置为 value 。
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
func (r *Rds) LSet(key string, index int64, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(index).Add(value)
	InfoFTimes(3, "[Redis Log] execute : LSET %s %v %v", key, index, value)
	_, err := r.Conn.Do("LSET", arg...)
	return err
}

// LTrim LTRIM key start stop
// 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
func (r *Rds) LTrim(key string, start, stop int64) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop)
	InfoFTimes(3, "[Redis Log] execute : LTRIM %s %v %v", key, start, stop)
	_, err := r.Conn.Do("LTRIM", arg...)
	return err
}

// RPop RPOP key
// 移除并返回列表 key 的尾元素。
func (r *Rds) RPop(key string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : RPOP %s", key)
	return redis.String(r.Conn.Do("RPOP", key))
}

// RPopLPush RPOPLPUSH source destination
// 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
// 将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
// 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
// 举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination
// 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ，
// destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作。
// 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
func (r *Rds) RPopLPush(key, destination string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(destination)
	InfoFTimes(3, "[Redis Log] execute : RPOPLPUSH %s %v", key, destination)
	return redis.String(r.Conn.Do("RPOPLPUSH", arg...))
}

// RPushX RPUSHX key value
// 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
func (r *Rds) RPushX(key string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : RPUSHX %s %v", key, value)
	_, err := r.Conn.Do("RPUSHX", arg...)
	return err
}

// Hash ===============================================================================================================

// HGetAll HGETALL 获取Hash value
func (r *Rds) HGetAll(key string) (map[string]string, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.StringMap(r.Conn.Do("HGETALL", key))
}

// HGetAllInt HGETALL 获取Hash value
func (r *Rds) HGetAllInt(key string) (map[string]int, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.IntMap(r.Conn.Do("HGETALL", key))
}

// HGetAllInt64 HGETALL 获取Hash value
func (r *Rds) HGetAllInt64(key string) (map[string]int64, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.Int64Map(r.Conn.Do("HGETALL", key))
}

// HSet HSET 新建Hash 单个field
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Rds) HSet(key, field string, value any) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(value)
	InfoFTimes(3, "[Redis Log] execute : HSET %s %v %v", key, field, value)
	return redis.Int64(r.Conn.Do("HSET", arg...))
}

// HMSet HMSET 新建Hash 多个field
// HMSET key field value [field value ...]
// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
func (r *Rds) HMSet(key string, values map[any]any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for k, v := range values {
		args = args.Add(k)
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : HMSET %s %s ", key, strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("HMSET", args...)
	return err
}

// HSetNx HSETNX key field value
// 给hash追加field value
// 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
func (r *Rds) HSetNx(key, field string, value any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(value)
	InfoFTimes(3, "[Redis Log] execute : HSETNX %s %v %v", key, field, value)
	_, err := r.Conn.Do("HSETNX", arg...)
	return err
}

// HDel HDEL key field [field ...] 删除哈希表
// key 中的一个或多个指定域，不存在的域将被忽略。
func (r *Rds) HDel(key string, fields []string) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range fields {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : HDEL %s %s", key, strings.Join(fields, " "))
	_, err := r.Conn.Do("HDEL", args...)
	return err
}

// HExIsTs HEXISTS key field 查看哈希表
// key 中，给定域 field 是否存在。
func (r *Rds) HExIsTs(key, fields string) bool {
	if r.Conn == nil {
		ErrorTimes(3, "[Redis Log] error :", RdsNotConnError)
		return false
	}
	arg := redis.Args{}.Add(key).Add(fields)
	InfoFTimes(3, "[Redis Log] execute : HEXISTS %s %s", key, fields)
	res, err := redis.Int(r.Conn.Do("HEXISTS", arg...))
	if err != nil {
		ErrorTimes(3, "[Redis Log] error :", err)
		return false
	}
	if res == 0 {
		return false
	}
	return true
}

// HGet HGET key field 返回哈希表
// key 中给定域 field 的值。
func (r *Rds) HGet(key, fields string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(fields)
	InfoFTimes(3, "[Redis Log] execute : HGET %s %s", key, fields)
	return redis.String(r.Conn.Do("HGET", arg...))
}

// HIncrBy HINCRBY key field increment
// 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0
func (r *Rds) HIncrBy(key, field string, increment int64) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(increment)
	InfoFTimes(3, "[Redis Log] execute : HINCRBY %s %v %v", key, field, increment)
	return redis.Int64(r.Conn.Do("HINCRBY", arg...))
}

// HIncrByFloat HINCRBYFLOAT key field increment
// 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Rds) HIncrByFloat(key, field string, increment float64) (float64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(increment)
	InfoFTimes(3, "[Redis Log] execute : HINCRBYFLOAT %s %v %v", key, field, increment)
	return redis.Float64(r.Conn.Do("HINCRBYFLOAT", arg...))
}

// HKeys HKEYS key 返回哈希表
// key 中的所有域。
func (r *Rds) HKeys(key string) ([]string, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HKEYS %s", key)
	return redis.Strings(r.Conn.Do("HKEYS", key))
}

// HLen HLEN key 返回哈希表
// key 中域的数量。
func (r *Rds) HLen(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HLEN %s", key)
	return redis.Int64(r.Conn.Do("HLEN", key))
}

// HMGet HMGET key field [field ...]
// 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
func (r *Rds) HMGet(key string, fields []string) ([]string, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range fields {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : HMGET %s %s", key, strings.Join(fields, " "))
	return redis.Strings(r.Conn.Do("HMGET", args...))
}

// HVaLs HVALS key
// 返回哈希表 key 中所有域的值。
func (r *Rds) HVaLs(key string) ([]string, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : HVALS %s", key)
	return redis.Strings(r.Conn.Do("HVALS", key))
}

// TODO HSCAN
// 搜索value hscan test4 0 match *b*

// Set ================================================================================================================

// SMemeRs SMEMBERS key
// 返回集合 key 中的所有成员。
// 获取Set value 返回集合 key 中的所有成员。
func (r *Rds) SMemeRs(key string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : SMEMBERS %s", key)
	return redis.Values(r.Conn.Do("SMEMBERS", key))
}

// SAdd SADD 新创建Set  将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
func (r *Rds) SAdd(key string, values []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	InfoFTimes(3, "[Redis Log] execute : SADD %s %s", key, strings.Join(Any2Strings(values), " "))
	_, err := r.Conn.Do("SADD", args...)
	return err
}

// SCard SCARD key
// 返回集合 key 的基数(集合中元素的数量)。
func (r *Rds) SCard(key string) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : SCARD %s", key)
	_, err := redis.Int64(r.Conn.Do("SCARD ", key))
	return err
}

// SDiff SDIFF key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (r *Rds) SDiff(keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	InfoFTimes(3, "[Redis Log] execute : SDIFF %s", strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SDIFF", args...))
}

// SDiffStore SDIFFSTORE destination key [key ...]
// 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Rds) SDiffStore(key string, keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, k := range keys {
		args = args.Add(k)
	}
	InfoFTimes(3, "[Redis Log] execute : SDIFFSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SDIFFSTORE", args...))
}

// SInter SINTER key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
func (r *Rds) SInter(keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	InfoFTimes(3, "[Redis Log] execute : SINTER %s", strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SINTER", args...))
}

// SInterStore SINTERSTORE destination key [key ...]
// 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Rds) SInterStore(key string, keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, k := range keys {
		args = args.Add(k)
	}
	InfoFTimes(3, "[Redis Log] execute : SINTERSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SINTERSTORE", args...))
}

// SIsMember SISMEMBER key member
// 判断 member 元素是否集合 key 的成员。
// 返回值:
// 如果 member 元素是集合的成员，返回 1 。
// 如果 member 元素不是集合的成员，或 key 不存在，返回 0 。
func (r *Rds) SIsMember(key string, value any) (resBool bool, err error) {
	if r.Conn == nil {
		return false, RdsNotConnError
	}
	resBool = false
	arg := redis.Args{}.Add(key).Add(value)
	InfoFTimes(3, "[Redis Log] execute : SISMEMBER %s %v", key, value)
	res, err := redis.Int64(r.Conn.Do("SISMEMBER", arg...))
	if err != nil {
		return
	}
	if res == 1 {
		resBool = true
		return
	}
	return
}

// SMove SMOVE source destination member
// 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则，
// member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
// 返回值: 成功移除，返回 1 。失败0
func (r *Rds) SMove(key, destination string, member any) (resBool bool, err error) {
	if r.Conn == nil {
		return false, RdsNotConnError
	}
	resBool = false
	arg := redis.Args{}.Add(key).Add(destination).Add(member)
	InfoFTimes(3, "[Redis Log] execute : SMOVE %s %v %v", key, destination, member)
	res, err := redis.Int64(r.Conn.Do("SMOVE", arg...))
	if err != nil {
		return
	}
	if res == 1 {
		resBool = true
		return
	}
	return
}

// SPop SPOP key
// 移除并返回集合中的一个随机元素。
func (r *Rds) SPop(key string) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : SPOP %s", key)
	return redis.String(r.Conn.Do("SPOP", key))
}

// SRandMember SRANDMEMBER key [count]
// 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。
// 如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
func (r *Rds) SRandMember(key string, count int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(count)
	InfoFTimes(3, "[Redis Log] execute : SRANDMEMBER %s %v", key, count)
	return redis.Values(r.Conn.Do("SRANDMEMBER", arg...))
}

// SRem SREM key member [member ...]
// 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
func (r *Rds) SRem(key string, member []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range member {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : SREM %s %s", key, strings.Join(Any2Strings(member), " "))
	_, err := r.Conn.Do("SREM", args...)
	return err
}

// SUnion SUNION key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的并集。
func (r *Rds) SUnion(keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}
	for _, v := range keys {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : SUNION %s", strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SUNION", args...))
}

// SUnionStore SUNIONSTORE destination key [key ...]
// 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
func (r *Rds) SUnionStore(key string, keys []string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range keys {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : SUNIONSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(r.Conn.Do("SUNIONSTORE", args...))
}

// TODO 搜索值  SSCAN key cursor [MATCH pattern] [COUNT count]

// ZSet ===============================================================================================================

// ZRange ZRANGE 获取ZSet value 返回集合 有序集成员的列表。
func (r *Rds) ZRange(key string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(0).Add(-1).Add("WITHSCORES")
	InfoFTimes(3, "[Redis Log] execute : ZRANGE %s 0 -1 ", "ZRANGE WITHSCORES", key)
	return redis.Values(r.Conn.Do("ZRANGE", arg...))
}

// ZRangeST ZRANGE key start stop [WITHSCORES]
// 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (r *Rds) ZRangeST(key string, start, stop int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop).Add("WITHSCORES")
	InfoFTimes(3, "[Redis Log] execute : ZRANGE %s %v %v WITHSCORES", key, start, stop)
	return redis.Values(r.Conn.Do("ZRANGE", arg...))
}

// ZRevRange ZREVRANGE key start stop [WITHSCORES]
// 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
func (r *Rds) ZRevRange(key string, start, stop int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop).Add("WITHSCORES")
	InfoFTimes(3, "[Redis Log] execute : ZREVRANGE %s %v %v WITHSCORES", key, start, stop)
	return redis.Values(r.Conn.Do("ZREVRANGE", arg...))
}

// ZAdd ZADD 新创建ZSet 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
func (r *Rds) ZAdd(key string, weight any, field any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key).Add(weight).Add(field)
	InfoFTimes(3, "[Redis Log] execute : ZADD %s %v %v", key, weight, field)
	_, err := r.Conn.Do("ZADD", args...)
	return err
}

// ZCard ZCARD key
// 返回有序集 key 的基数。
func (r *Rds) ZCard(key string) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	InfoFTimes(3, "[Redis Log] execute : ZCARD %s", key)
	return redis.Int64(r.Conn.Do("ZCARD", key))
}

// ZCount ZCOUNT key min max
// 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (r *Rds) ZCount(key string, min, max int64) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(min).Add(max)
	InfoFTimes(3, "[Redis Log] execute : ZCOUNT %s %v %v", key, min, max)
	return redis.Int64(r.Conn.Do("ZCOUNT", arg...))
}

// ZIncrBy ZINCRBY key increment member
// 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
func (r *Rds) ZIncrBy(key, member string, increment int64) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(increment).Add(member)
	InfoFTimes(3, "[Redis Log] execute : ZINCRBY %s %v %v", key, increment, member)
	return redis.String(r.Conn.Do("ZINCRBY", arg...))
}

// ZRangeByScore ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
// 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
// 具有相同 score 值的成员按字典序(lexicographical order)来排列(该属性是有序集提供的，不需要额外的计算)。
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count )，注意当 offset 很大时，
// 定位 offset 的操作可能需要遍历整个有序集，此过程最坏复杂度为 O(N) 时间。
// 可选的 WITHSCORES 参数决定结果集是单单返回有序集的成员，还是将有序集成员及其 score 值一起返回。
// 区间及无限
// min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。
// 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。
func (r *Rds) ZRangeByScore(key string, min, max, offset, count int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(min).Add(max).Add(offset).Add(count)
	InfoFTimes(3, "[Redis Log] execute : ZRANGEBYSCORE %s %v %v %v %v", key, min, max, offset, count)
	return redis.Values(r.Conn.Do("ZRANGEBYSCORE", arg...))
}

// ZRangeByScoreAll 获取所有
func (r *Rds) ZRangeByScoreAll(key string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add("-inf").Add("+inf")
	InfoFTimes(3, "[Redis Log] execute : ZRANGEBYSCORE %s -inf +inf", key)
	return redis.Values(r.Conn.Do("ZRANGEBYSCORE", arg...))
}

// ZRevRangeByScore key max min [WITHSCORES] [LIMIT offset count]
// 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。
func (r *Rds) ZRevRangeByScore(key string, min, max, offset, count int64) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(min).Add(max).Add(offset).Add(count)
	InfoFTimes(3, "[Redis Log] execute : ZREVRANGEBYSCORE %s %v %v %v %v", key, min, max, offset, count)
	return redis.Values(r.Conn.Do("ZREVRANGEBYSCORE", arg...))
}

// ZRevRangeByScoreAll 获取所有
func (r *Rds) ZRevRangeByScoreAll(key string) ([]any, error) {
	if r.Conn == nil {
		return nil, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add("-inf").Add("+inf")
	InfoFTimes(3, "[Redis Log] execute : ZREVRANGEBYSCORE %s -inf +inf", key)
	return redis.Values(r.Conn.Do("ZREVRANGEBYSCORE", arg...))
}

// ZRank ZRANK key member
// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
func (r *Rds) ZRank(key string, member any) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(member)
	InfoFTimes(3, "[Redis Log] execute : ZRANK %s %v", key, member)
	return redis.Int64(r.Conn.Do("ZRANK", arg...))
}

// ZRem ZREM key member [member ...]
// 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
func (r *Rds) ZRem(key string, member []any) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range member {
		args = args.Add(v)
	}
	InfoFTimes(3, "[Redis Log] execute : ZREM %s %s", key, strings.Join(Any2Strings(member), " "))
	_, err := r.Conn.Do("ZREM", args...)
	return err
}

// ZRemRangeByRank ZREMRANGEBYRANK key start stop
// 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
func (r *Rds) ZRemRangeByRank(key string, start, stop int64) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop)
	InfoFTimes(3, "[Redis Log] execute : ZREMRANGEBYRANK %s %v %v", key, start, stop)
	_, err := redis.Int64(r.Conn.Do("ZREMRANGEBYRANK", arg...))
	return err
}

// ZRemRangeByScore ZREMRANGEBYSCORE key min max
// 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
func (r *Rds) ZRemRangeByScore(key string, min, max int64) error {
	if r.Conn == nil {
		return RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(min).Add(max)
	InfoFTimes(3, "[Redis Log] execute : ZREMRANGEBYSCORE %s %v %v", key, min, max)
	_, err := r.Conn.Do("ZREMRANGEBYSCORE", arg...)
	return err
}

// ZRevRank ZREVRANK key member
// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (r *Rds) ZRevRank(key string, member any) (int64, error) {
	if r.Conn == nil {
		return 0, RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(member)
	InfoFTimes(3, "[Redis Log] execute : ZREVRANK %s %v", key, member)
	return redis.Int64(r.Conn.Do("ZREVRANK", arg...))
}

// ZScore ZSCORE key member
// 返回有序集 key 中，成员 member 的 score
func (r *Rds) ZScore(key string, member any) (string, error) {
	if r.Conn == nil {
		return "", RdsNotConnError
	}
	arg := redis.Args{}.Add(key).Add(member)
	InfoFTimes(3, "[Redis Log] execute : ZSCORE %s %v", key, member)
	return redis.String(r.Conn.Do("ZSCORE", arg...))
}

// TODO ZSetZUNIONSTORE ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]
// 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。

// TODO ZSetZINTERSTORE ZINTERSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]
// 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.

// TODO 搜索值  ZSCAN key cursor [MATCH pattern] [COUNT count]
