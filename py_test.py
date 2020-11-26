# coding=utf-8
# len1 = str(input(':'))
# len2 = str(input(':'))

print (0x85ebca6b)

print("This is a simple add test...")
while 1:
    # raw_input 得到的数据是不包含括号的，输入什么就得到什么，相比于input
    _input = input("please input a str that split by `,` (input `exit` to exit):")
    print("--- recognized input: %s" % str(_input))
    # 需要去除括号, 空格(仅python2需要，python3的input输出只包含输入的内容)
    # s = str(_input).replace("(", "").replace(")", "").replace(" ", "")
    # print("--- processed input: %s" % s)

    ss = str(_input).split(",")

    if len(ss) == 1 and ss[0] == "exit":
        print("exited, bye!")
        break  # or exit(0)
    elif len(ss) != 2:
        print("invalid input, try again!")
        continue

    new_ss = []
    # 判断类型
    for i in ss:
        int_i = 0
        try:
            # print(i)
            # 尝试转换str-->int
            int_i = int(i)
        except:
            print("not a integer, try again!")
            break  # 跳出for循环
        new_ss.append(int_i)
    else:
        print("OK，%s + %s = %s" % (ss[0], ss[1], sum(new_ss)))
