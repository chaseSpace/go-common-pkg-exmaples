## sync.Pool 源码分析

基于 Go1.12

## 用途与特点解释
- 一个Pool是一组可以单独保存和检索的临时对象
- 任何存储在Pool中的对象都可能会被毫无通知的删除(回收)，如果对象是指针，则可能会变成nil
- `sync.Pool`是goroutine安全的
- Pool的作用是缓存多个已分配但暂时未使用的对象，以便在下次需要时直接从Pool得到，而不是再创建
    在gc(垃圾回收)时会释放压力，适用于构建高性能，线程安全的空闲列表
- 一个推荐的场景是：管理一组在多个goroutine之间共享的临时对象
- `sync.Pool`提供了一种方式让不同goroutine的client均摊分配开销
- 一个很好例子是标准库fmt对Pool的使用，fmt包维护一个动态大小的Pool，缓存的对象是output buffers，
    这个Pool的容量会根据程序中对fmt的使用负载实现动态的扩容/缩容
- Pool对象在第一次使用后不能被Copy(因为内部会创建实体对象)

<br/>


> 注：阅读此文需要对Go运行时调度有一定的了解

## Pool结构
```go
type Pool struct {
	noCopy noCopy // 用于内存对齐

    // 真实类型是slice：[P]poolLocal, P是CPU核心数，这表示每个cpu核心独享一个poolLocal（理解这里对后面很重要）
    // 为什么叫local：与ThreadLocal含义类似，在这里表示所在单个CPU核心独享的对象，也就是poolLocal
	local     unsafe.Pointer
    // local数组的大小 
	localSize uintptr       

    // New创建一个新的对象，用在pool为空的时候
	New func() interface{}
}
```

## pool.Get()
```go
func (p *Pool) Get() interface{} {
	if race.Enabled { // 检查是否开启了竞态检测，推荐阅读：https://golang.org/doc/articles/race_detector.html
		race.Disable() // 若开启则暂时禁用它
	}
    /*
     pin操作内部调用了`runtime_procPin`会将当让goroutine去抢占CPU核心, 在主动让出前，其他goroutine都抢不了这个P
    */ 
	l := p.pin() // 获取当前CPU核心的poolLocal（此核心独享的pool对象），内部逻辑在后面解读
	x := l.private // poolLocal存储的【私有】临时对象复制到x，如果非nil，说明这个CPU核心的pool被Put过，在后面会直接返回这个对象
	l.private = nil // 取走了就要把这个P的private缓存清除
	runtime_procUnpin() // 在这里让出CPU核心控制权
	if x == nil { // 如果poolLocal中的临时对象为nil，就检查poolLocal中的临时对象数组是否有库存，有就从里面取
		l.Lock() // 即将要操作的是shared对象，这是个在多CPU核心之间会【共享】的切片对象，需要加锁
		last := len(l.shared) - 1 // 算出数组最后一个元素的index
		if last >= 0 {  // 若有元素
            // 从slice中取出来
			x = l.shared[last]
			l.shared = l.shared[:last]
		}
		l.Unlock() // shared对象操作完毕，解锁
		if x == nil { // 如果poolLocal完全空
            // 就尝试从其他CPU核心的poolLocal去偷一个对象
            // 为啥叫getSlow？因为这个方法的有一定的时间复杂度，后文会解读
			x = p.getSlow() 
		}
	}
	if race.Enabled { // race原本是开启的
		race.Enable() // 再开启竞态检测
		if x != nil { 
			race.Acquire(poolRaceAddr(x)) // 未知操作？
		}
	}
    // 这种情况就是没有Put过，Pool内都是空的
	if x == nil && p.New != nil {
		x = p.New()
	}
	return x
}
```

## pool.getSlow()
```go
// 由pool.Get内部调用
func (p *Pool) getSlow() (x interface{}) {
	// See the comment in pin regarding ordering of the loads. [源码注释]
	size := atomic.LoadUintptr(&p.localSize) // load-acquire [源码注释] local的大小
	local := p.local                         // load-consume [源码注释] 这个local真实类型是数组 [P]poolLocal
	// 接下来尝试从其他CPU-core上的poolLocal偷一个对象
    // runtime_procPin返回当前goroutine所在的CPU-core id（这个方法同时会占住所在的P）
    // 0<= pid < Number(CPU-cores) 比如8核处理器，pid的范围是 [0,8)
	pid := runtime_procPin()
	runtime_procUnpin() // 取消P的抢占，允许被其他goroutine抢占
    
    // for循环可以计算其他P的poolLocal的内存地址，并通过unsafe.Pointer转换为实体类型
    // 简单来说，这个for循环就是在遍历其他P上的poolLocal对象，尝试从这些对象的shared数组中偷一个item
	for i := 0; i < int(size); i++ {
        // indexLocal方法输入两个参数，第一个是p.local数组的指针地址（也是该数组第一个元素的地址）
		// 第二个是使用pid和size执行取模运算得出的一个数
        // indexLocal计算出其他P的poolLocal的内存地址，并转换为poolLocal对象
		l := indexLocal(local, (pid+i+1)%int(size))
        // 下面要操作计算出的这个poolLocal的shared对象，因为每个poolLocal的shared字段可能被其他P访问和修改，所以每个poolLocal持有一把锁
        // 这是一种分段锁思想，降低了锁的粒度，提高并发性能，Java中的ConcurrentHashMap有用到
		l.Lock() // 上锁
        // 从slice中pop一个元素返出去
		last := len(l.shared) - 1
		if last >= 0 {
			x = l.shared[last]
			l.shared = l.shared[:last]
			l.Unlock() // 最后解锁
			break
		}
		l.Unlock() // 最后解锁
	}
	return x
}

// 这个方法执行指针运算
// 可以这么来看待这个函数的运行过程
// 假设当前CPU是4核的，那就有4个poolLocal简称pl，第一个pl在内存中的地址是addr_1=0x00,
// 已知一个pl占用4个字节(假设的)，可以打印unsafe.Sizeof(poolLocal{})，那么第二个pl的地址就是0x04, 也等于addr_1+(2-1)*4
// 简单来说，poolLocal是数组中的元素，既然是数组，就可以通过一个pl的地址和目标pl的相对位置计算出目标pl的内存地址
// 在这里，参数l是p.local,真实类型是[P]poolLocal, 其内存地址也一定等于首个元素的内存地址
func indexLocal(l unsafe.Pointer, i int) *poolLocal {
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*unsafe.Sizeof(poolLocal{}))
	return (*poolLocal)(lp)
}
```

```go
// pin获取当前CPU核心的poolLocal（此核心独享的pool对象）
func (p *Pool) pin() *poolLocal {
    // 当前g抢占P，其他g只能等待当前g主动让出，在这期间，不会发生gc
	pid := runtime_procPin()
    // 读取localSize，复制一个local
	s := atomic.LoadUintptr(&p.localSize) // load-acquire
	l := p.local                          // load-consume
    // 若当前pid小于localSize，说明这个P在p.local中一定拥有一个poolLocal，计算内存地址找出来即可
    // 在程序运行起来后，这个if条件满足应该是频繁的情况，因为pid的范围是[0,NumberCPU), p.localSize=NumberCPU
    // 所以pid<s
	if uintptr(pid) < s {
		return indexLocal(l, pid)
	}
    // 但是在第一次执行Put/Get时，s=0，因为还没创建p.local,就需要创建p.local
    // 因为有了更多的步骤，所以是slow path
    // （slow path和fast path是CS领域内的一个术语，可自查）
	return p.pinSlow()
}
```

```go
// pinSlow主要做的就是创建local（包含NumberCpu个poolLocal）, 并以指针地址的形式分配给p.local
func (p *Pool) pinSlow() *poolLocal {
	// Retry under the mutex.
	// Can not lock the mutex while pinned.
    // 官方注释只说了在抢占的时候不能上全局锁，没说why
    // 先取消P的抢占
	runtime_procUnpin()
    // 上一个Pool-pkg级的全局锁
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()
    // 上了全局锁后，再抢占 （不知道为何这样操作）
	pid := runtime_procPin()
	// poolCleanup won't be called while we are pinned.
    // 复制一个p.localSize 和 p.local
	s := p.localSize
	l := p.local
    // 再判断一次有没有创建过p.local（因为上面取消了抢占，可能被别的g创建了）
	if uintptr(pid) < s {
		return indexLocal(l, pid)
	}
    // 这里还需要判断local是否nil，为什么？因为导致pid<s的情况不只一种：
    // - 1. 从未创建过p.local（第一次执行Put/Get）
    // - 2. runtime.GOMAXPROCS 在程序运行过程中变化了（一般我们没法通过runtime修改这个，有可能是外部给机器扩容，比较极端）
	if p.local == nil {
		allPools = append(allPools, p)
	}
    // 这一句是官方注释，说如果GOMAXPROCS在gc期间变化了，就重新创建local数组并分配给p.local, 旧的丢弃
	// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal, size)
    // 这里是通过指针的形式把local分配给p.local
	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
	atomic.StoreUintptr(&p.localSize, uintptr(size))         // store-release
	return &local[pid]
}

```

理解了`p.Get`，`p.Put`也就容易理解了

```go
// Put adds x to the pool.
func (p *Pool) Put(x interface{}) {
    // 不允许put nil
	if x == nil {
		return
	}
	if race.Enabled {
		if fastrand()%4 == 0 {
			// Randomly drop x on floor.
			return
		}
		race.ReleaseMerge(poolRaceAddr(x))
		race.Disable()
	}
    // pin做的事
	l := p.pin()
	if l.private == nil {
		l.private = x
		x = nil
	}
	runtime_procUnpin()
	if x != nil {
		l.Lock()
		l.shared = append(l.shared, x)
		l.Unlock()
	}
	if race.Enabled {
		race.Enable()
	}
}
```

go的实现包含垃圾回收(gc)，减少了开发者心智负担，但也增加runtime开销，使用不当也会影响程序性能，即性能要求高的场景不能产生太多的
垃圾，那怎么解决？ 当然就是复用一个已创建的对象了，sync.Pool就是做这样一个事情

#### 使用示例
```go
var p = &sync.Pool{
    New: func() interface{} {
        return 0
    },
}
v1 := p.Get().(int)
p.Put(1)
fmt.Println(v1, p.Get().(int)) // 0， 1
```

需要注意的是Pool不支持设置最大缓存数量和时间的，那万一缓存过多导致太大的内存开销呢？ 看代码：
```go
// 在上面代码的基础上增加一行
var p = &sync.Pool{
    New: func() interface{} {
        return 0
    },
}
v1 := p.Get().(int)
p.Put(1)
runtime.GC() // <----
fmt.Println(v1, p.Get().(int)) // 0， 0
```

sync.Pool缓存对象的期限是两次GC期间，什么时候清空的，怎么清空的？看代码:
```go
func init() {
    // 把清除函数注册到runtime中去，下次gc开始时会调用这个函数
    runtime_registerPoolCleanup(poolCleanup)
}
func poolCleanup() {
    // 【官方注释】
    // This function is called with the world stopped, at the beginning of a garbage collection.
    // It must not allocate and probably should not call any runtime functions.
    // Defensively zero out everything, 2 reasons:
    // 1. To prevent false retention of whole Pools.
    // 2. If GC happens while a goroutine works with l.shared in Put/Get,
    //    it will retain whole Pool. So next cycle memory consumption would be doubled.
    // 遍历所有创建的Pool，挨个重置，重置的方式也很暴力，=nil这种方式这一定会产生很多垃圾对象，不过还好是在gc即将开始的时候清除
    for i, p := range allPools {
        allPools[i] = nil
        for i := 0; i < int(p.localSize); i++ {
            l := indexLocal(p.local, i)
            l.private = nil
            for j := range l.shared {
                l.shared[j] = nil
            }
            l.shared = nil
        }
        p.local = nil
        p.localSize = 0
    }
    allPools = []*Pool{}
}
```
所以，sync.Pool并不会像常见的各种连接池那样长久的保持一些对象，不可以用于类似的用途。

#### 如何做到最小的并发竞争？
在一个全局Pool中缓存一组对象，然后去获取对象，这其中就涉及到多g(goroutine)之间的竞争，如果只是简单用一把大锁来串行化所有g的
操作，那就太影响并发性能了，若你仔细看了代码，你就会发现sync.Pool中应用了缓存+分段锁的设计思想，这里再赘述一下从Pool中获取对象的步骤：
- 首先检查当前g所运行的P(CPU-core)的private字段是否有值，有直接返回，否则继续 【缓存思想】
- 检查P拥有的shared数组对象（需要上一把小锁，即每个P单独的锁），有则pop一个返回，否则继续 【分段锁思想】
- 遍历其他的P，直至找一个对象，然后返回它（slow path）

你还可以参考的中文文章： https://blog.csdn.net/yongjian_lian/article/details/42058893

完。