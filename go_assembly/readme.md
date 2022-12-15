# 汇编

汇编语言是机器语言利用助记符的表示，表达能力上等价。指令集是对CPU架构硬件的抽象，汇编语言是对指令集的一种描述。
学习汇编语言的过程和理解CPU硬件架构是分不开的，包括但不限于寄存器组，寻址方式，指令字长及格式，支持的指令操作等。

所以，汇编语言的种类一般是根据`CPU架构`区分的，CPU架构一般有两大类，`x86`和`arm`；除此之外，还要根据具体CPU架构下的
指令集进行区分

## 指令集
它是CPU中用来计算和控制计算机操作系统的一套指令的集合，每种CPU在设计时就规定了一系列与其他硬件相配合的指令集。
>指令集的先进与否，密切关系到CPU的性能发挥

通俗的理解，指令集就是CPU直接接受的语言

### 指令集分类
先从大类区分，一般分为精简指令集(CISC)和复杂指令集(RISC) 。

CISC的特点：
- 指令数目少，每条指令都采用标准字长，执行时间短
- 指令顺序执行，优点是控制简单，但计算机各部分的利用率不高，执行速度慢
- 适用场景：实现较为复杂的功能

RISC的特点：
- 是指针CISC指令集的一些常用指令的优化设计，放弃了一些复杂指令，通过组合指令来完成复杂操作
- 适用场景：实现不那么复杂的功能，且低功耗

>一些历史背景：最开始Intel X86的第一个CPU定义了一套指令集，这是最开始的指令集。后来一些公司发现很多指令并不常用，
> 所以决定再设计一套简单高效的指令集即RISC，把原来的叫做CISC。

典型的CISC指令集：Intel的x86指令集，AMD的x86-64指令集
典型的RISC指令集：ARM、MIPS等

不管是CISC还是RISC类指令集，都在持续发展，所以它们各自都会派生出更多指令集

## 汇编语法

Intel和AT&T两种，可以理解为两种风格，符号系统有点差别

## 8086汇编又是什么

8086和x86是Intel的两种处理器系列，其中8086一般指Intel于1978年所设计的16位微处理器芯片;
x86是80x86系列，80286开始有了保护模式，80386（1985）是32位CPU，以及后续还有80486，80586等。

同厂商下的不同系列CPU可能有些不同，所以可能又单独使用一种以该系列CPU命名的汇编语言，但主要语法是差不多的。

## 如何选择哪种汇编语言学习？

先确定你学后用来应用的指令集架构是哪种（x86还是arm），再确定想要学哪种风格（Intel还是AT&T），不同风格的汇编语言
的资料丰富程度不同（x86较多），自行搜索