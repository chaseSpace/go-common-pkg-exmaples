package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
)

/*
Usage of redigo:
https://godoc.org/github.com/gomodule/redigo/redis#pkg-examples


*************常见的err类型
返回这个err表示没有数据
var ErrNil = errors.New("redigo: nil returned")

返回这个err表示连接池满了(Do, Send, Receive, Flush, Err连接池的这几个方法可能返回)
var ErrPoolExhausted = errors.New("redigo: connection pool exhausted")

#############

*************常用的结果转换方法--官方叫helper
#慎用！！！ 如果没有数据返回，err将等于 “nil returned”, 此时如果要判断redis异常就需要多判断一次 err != "nil returned"
func String(reply interface{}, err error) (string, error)  下方有示例
func Strings(reply interface{}, err error) ([]string, error)
func StringMap(result interface{}, err error) (map[string]string, error) HGETALL and CONFIG GET可能会用到
func Bool(reply interface{}, err error) (bool, error)
func Uint64(reply interface{}, err error) (uint64, error)
func ByteSlices(reply interface{}, err error) ([][]byte, error)
func Bytes(reply interface{}, err error) ([]byte, error)
func Float64(reply interface{}, err error) (float64, error)
func Float64s(reply interface{}, err error) ([]float64, error)
func Int(reply interface{}, err error) (int, error)
func Ints(reply interface{}, err error) ([]int, error)
func Int64(reply interface{}, err error) (int64, error)
func IntMap(result interface{}, err error) (map[string]int, error) HGETALL返回的结果可能会用到
func Int64Map(result interface{}, err error) (map[string]int64, error) HGETALL返回的结果可能会用到
func Int64s(reply interface{}, err error) ([]int64, error)
func MultiBulk(reply interface{}, err error) ([]interface{}, error)  不推荐，改用Values代替
func Values(reply interface{}, err error) ([]interface{}, error) 转换数组类型的结果


*************常用的方法
#控制超时
func DoWithTimeout(c conn, timeout time.Duration, cmd string, args ...interface{}) (interface{}, error)
func ReceiveWithTimeout(c conn, timeout time.Duration) (interface{}, error)

#redis命令构造结构体--type Args []interface{} ::比如HMSET多个KV时就会用到
func (args Args) Add(value ...interface{}) Args
func (args Args) AddFlat(v interface{}) Args ::将传入的Map、slice、struct进行扁平化展开

#redis Conn对象
type conn interface {
    	// Close closes the connection.
    Close() error
    	// Err returns a non-nil value when the connection is not usable.
    Err() error
    	// Do sends a command to the server and returns the received reply.
    Do(commandName string, args ...interface{}) (reply interface{}, err error)
    	// Send writes the command to the client's output buffer.
    Send(commandName string, args ...interface{}) error
   		// Flush flushes the output buffer to the Redis server.
    Flush() error
    	// Receive receives a single reply from the Redis server
    Receive() (reply interface{}, err error)
}

#连接redis
func Dial(network, address string, options ...DialOption) (conn, error) ::如redis.Dial("tcp", ":6379") 可通过options控制超时
func DialURL(rawurl string, options ...DialOption) (conn, error) ::如redis.DialURL(os.Getenv("REDIS_URL"))

#dialOptions
func DialClientName(name string) DialOption 目前版本不支持
func DialConnectTimeout(d time.Duration) DialOption  下面有示例
func DialDatabase(db int) DialOption
func DialKeepAlive(d time.Duration) DialOption
func DialPassword(password string) DialOption
func DialReadTimeout(d time.Duration) DialOption
func DialWriteTimeout(d time.Duration) DialOption
	下面几个用的少
func DialTLSConfig(c *tls.Config) DialOption
func DialTLSSkipVerify(skip bool) DialOption
func DialUseTLS(useTLS bool) DialOption
func DialNetDial(dial func(network, addr string) (net.conn, error)) DialOption

#连接池【常用】 下面有示例
type Pool struct {...} 未展开 https://godoc.org/github.com/gomodule/redigo/redis#Pool
func NewPool(newFn func() (conn, error), maxIdle int) *Pool  已废弃
func (p *Pool) ActiveCount() int
func (p *Pool) Close() error
func (p *Pool) Get() conn
func (p *Pool) GetContext(ctx context.Context) (conn, error)
func (p *Pool) IdleCount() int
func (p *Pool) Stats() PoolStats


*************不常用的方法
func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error)
func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error
func ScanStruct(src []interface{}, dest interface{}) error
func NewConn(netConn net.conn, readTimeout, writeTimeout time.Duration) conn 基于已有的net conn返回一个新的 redigo conn
func NewLoggingConn(conn conn, logger *log.Logger, prefix string) conn 返回日志封装过的conn
func NewLoggingConnFilter(conn conn, logger *log.Logger, prefix string, skip func(cmdName string) bool) conn

*************不常用的结构体 PubSubConn
type PubSubConn struct { conn conn }  发布订阅
func (c PubSubConn) Close() error
func (c PubSubConn) PSubscribe(channel ...interface{}) error
func (c PubSubConn) PUnsubscribe(channel ...interface{}) error
func (c PubSubConn) Ping(data string) error
func (c PubSubConn) Receive() interface{}
func (c PubSubConn) ReceiveWithTimeout(timeout time.Duration) interface{}
func (c PubSubConn) Subscribe(channel ...interface{}) error
func (c PubSubConn) Unsubscribe(channel ...interface{}) error


*************不常用的结构体 Scanner https://godoc.org/github.com/gomodule/redigo/redis#Scanner
type Scanner interface {
    RedisScan(src interface{}) error
}
# RedisScan这个方法用来将redis服务器返回的数据copy到目标变量

*************不常用的结构体 Script
type Script struct {} 操作redis脚本 https://godoc.org/github.com/gomodule/redigo/redis#Script

*************不常用的结构体 Subscription
type Subscription struct {}

*/

var err error
var conn redis.Conn
var pool *redis.Pool

var (
	pass = "123"
)

func BaseConn() {
	conn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
		return
	}
	// do authorize
	_, _ = conn.Do("AUTH", pass)
	log.Println("conn is ready...")
}

func TestBase(t *testing.T) {
	BaseConn()
}

func Test_String(t *testing.T) {
	//here is db 0 by default.
	s, err := conn.Do("SET", "name", "lei")
	if err != nil {
		fmt.Println(2, err)
		return
	}
	fmt.Println("set:", s)
	s, err = redis.String(conn.Do("GET", "name"))
	if err != nil {
		fmt.Println(3, err)
		return
	}
	fmt.Println("get:", s)
}

func Test_StringMap(t *testing.T) {
	// 专用于hash散列类型的转换
	var p struct {
		Title  string `redis:"title"`
		Author string `redis:"author"`
		Body   string `redis:"body"`
	}
	p.Title = "Example"
	p.Author = "Gary"
	p.Body = "Hello"

	s, err := conn.Do("HMSET", redis.Args{}.Add("hash_key").AddFlat(&p)...)
	if err != nil {
		fmt.Println("Test_StringMap-1:", err)
		return
	}
	fmt.Println("Test_StringMap-set:", s)

	s_map, _ := redis.StringMap(conn.Do("HGETALL", "hash_key"))
	for k, v := range s_map {
		fmt.Printf("Test_StringMap-k:%v v:%v\n", k, v)
	}
}

func Test_Dial(t *testing.T) {

	// 这个库2.0.0版本不支持传入account，估计固定为root了，最新master支持了clientName，还没发版
	var (
		keepAliveOpt   = redis.DialKeepAlive(3 * time.Minute)
		connTimeoutOpt = redis.DialConnectTimeout(3 * time.Second)
		pwdOpt         = redis.DialPassword("123")
		readTimeoutOpt = redis.DialReadTimeout(3 * time.Second)
		writeOpt       = redis.DialWriteTimeout(3 * time.Second)
	)

	conn, err = redis.Dial("tcp", ":6379", keepAliveOpt,
		connTimeoutOpt,
		pwdOpt,
		readTimeoutOpt,
		writeOpt)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Dial ok")
	// URL Example: redis://user:secret@host:port/db_number?foo=bar&qux=baz
	// rediss://  表示支持ssl的redis服务
	conn, err = redis.DialURL("redis://root:123@127.0.0.1:6379/0")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("DialURL ok")
}

func Test_Pool(t *testing.T) {

	// check health of conn before return to application
	TestOnBorrow := func(c redis.Conn, t time.Time) error {
		// 这个conn在1分钟内被使用过则不用ping检查
		if time.Since(t) < time.Minute {
			return nil
		}
		// 否则通过ping检查可用性
		_, err := c.Do("PING")
		return err
	}
	// create a new pool
	pool = &redis.Pool{
		MaxActive:    8,                 // 最大连接数，比CPU线程数稍大些即可
		MaxIdle:      3,                 // 最大空闲连接数
		IdleTimeout:  240 * time.Second, //超过这个时间则关闭连接，0表示不关闭
		Wait:         false,             // 连接数满无空闲时 是否等待
		TestOnBorrow: TestOnBorrow,
		Dial: func() (redis.Conn, error) {
			pConn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			_, _ = pConn.Do("AUTH", pass)
			return pConn, err
		},
	}

	// get a conn from pool
	operationFunc := func() {
		conn_from_pool := pool.Get()
		defer conn_from_pool.Close()
		s, _ := redis.String(conn_from_pool.Do("GET", "name"))

		fmt.Println("GET name", s)
	}

	// just do it
	operationFunc()
}

func TestExists(t *testing.T) {
	defer pool.Close()
	rdsConn := pool.Get()
	defer rdsConn.Close()
	_, _ = rdsConn.Do("FLUSHALL")
	s, err := rdsConn.Do("EXISTS", "name1") // not exists
	fmt.Println(fmt.Sprintf("%v", s), err)  // 0,nil

	s, err = rdsConn.Do("EXISTS", "name")  // exists
	fmt.Println(fmt.Sprintf("%v", s), err) // 1,nil

	s, err = rdsConn.Do("GET", "name1")
	fmt.Println(s == nil, err) // true,nil

	_, _ = rdsConn.Do("SET", "KeysTest", "xxx")

	s, err = rdsConn.Do("KEYS", "*")
	fmt.Printf("raw keys type:%T err:%v\n", s, err) // []interface, nil

	// keys * 方法可以使用redis.Strings
	s, err = redis.Strings(rdsConn.Do("KEYS", "*"))
	fmt.Println(s, err) // [],nil  若无key，err=nil

	// set k v(str==[]byte) key-value是[]byte也会被转str，都可以用
	_, _ = rdsConn.Do("SET", []byte("byteK"), []byte("xxx"))
	s, err = rdsConn.Do("GET", "byteK")
	fmt.Println(string(s.([]byte)), err) // get 返回的都是[]byte类型
}
