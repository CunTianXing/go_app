序言
Golang遵循“少即是多”的设计哲学，同时又支持闭包(Closure)，那么闭包对于Golang来说肯定有重要的价值。

对于Golang的初学者来说，肯定会有下面的几个疑问：

1、闭包是什么？
2、闭包是怎么产生的？
3、闭包可以解决什么问题？

闭包在函数式编程中广泛使用，所以一提起闭包，读者必然会想起函数式编程，我们先简单回顾一下。

在过去十几年的时间里，面向对象编程大行其道，以至于在大学的教育里，老师也只会教给我们两种编程模型，即面向过程和面向对象。孰不知，在面向对象思想产生之前，函数式编程已经有了数十年的历史。

函数式编程在维基百科中的定义：

In computer science, functional programming is a programming paradigm that treats computation as the evaluation of mathematical functions and avoids state and mutable data.
简单翻译一下，函数式编程是一种编程模型，它将计算机运算看做是数学中函数的计算，并且避免了状态以及变量的概念。

闭包是由函数及其相关的引用环境组合而成的实体，即闭包=函数+引用环境。

这个定义从字面上很难理解，特别对于一直使用命令式语言进行编程的读者，所以本文将结合代码实例进行阐述。

函数是什么
函数是一段可执行代码，编译后就“固化”了，每个函数在内存中只有一份实例，得到函数的入口点便可以执行函数了。在函数式编程语言中，函数是一等公民（First class value：第一类对象，我们不需要像命令式语言中那样借助函数指针，委托操作函数），函数可以作为另一个函数的参数或返回值，可以赋给一个变量。函数可以嵌套定义（嵌套的函数一般为匿名函数），即在一个函数内部可以定义另一个函数，有了嵌套函数这种结构，便会产生闭包。

在面向对象编程中，我们将对象传来传去，而在函数式编程中，我们将函数传来传去。 在函数式编程中，高阶函数是至少满足以下两点的函数：

1、函数可以作为参数被传递
2、函数可以作为返回值输出

匿名函数
匿名函数是指不需要定义函数名的一种函数实现方式，它并不是一个新概念，最早可以回溯到1958年的Lisp语言。但是由于各种原因，C和C++一直都没有对匿名函数给以支持。

匿名函数由一个不带函数名的函数声明和函数体组成，比如：

func(x，y int) int {
    return x + y
}
在Golang中，所有的函数是值类型，即可以作为参数传递，又可以作为返回值传递。

匿名函数可以赋值给一个变量：

f := func() int {
    ...
}
我们可以定义一种函数类型：

type CalcFunc func(x, y int) int
函数可以作为值传递：

func AddFunc(x, y int) int {
    return x + y
}

func SubFunc(x, y int) int {
    return x - y
}

...

func OperationFunc(x, y int, calcFunc CalcFunc) int {
    return calcFunc(x, y)
}

func main() {
    sum := OperationFunc(1, 2, AddFunc)
    difference := OperationFunc(1, 2, SubFunc)
    ...
}
函数可以作为返回值：

// 第一种写法
func add(x, y int) func() int {
    f := func() int {
        return x + y
    }
    return f
}

// 第二种写法
func add(x, y int) func() int {
    return func() int {
        return x + y
    }
}
当函数返回多个匿名函数时建议采用第一种写法：

func calc(x, y int) （func(int), func()) {
    f1 := func(z int) int {
        return (x + y) * z / 2
    }

    f2 := func() int {
        return 2 * (x + y)
    }
    return f1, f2
}
匿名函数的调用有两种方法：

// 通过返回值调用
func main() {
    f1, f2 := calc(2, 3)
    n1 := f1(10)
    n2 := f1(20)
    n3 := f2()
    fmt.Println("n1, n2, n3:", n1, n2, n3)
}

// 在匿名函数定义的同时进行调用：花括号后跟参数列表表示函数调用
func safeHandler() {
    defer func() {
        err := recover()
        if err != nil {
            fmt.Println("some exception has happend:", err)
        }
    }()
    ...
}
闭包的本质
闭包是包含自由变量的代码块，这些变量不在这个代码块内或者任何全局上下文中定义，而是在定义代码块的环境中定义。由于自由变量包含在代码块中，所以只要闭包还被使用，那么这些自由变量以及它们引用的对象就不会被释放，要执行的代码为自由变量提供绑定的计算环境。

闭包的价值在于可以作为函数对象或者匿名函数，对于类型系统而言，这意味着不仅要表示数据还要表示代码。支持闭包的多数语言都将函数作为第一级对象，就是说这些函数可以存储到变量中作为参数传递给其他函数，最重要的是能够被函数动态创建和返回。

Golang中的闭包同样也会引用到函数外的变量，闭包的实现确保只要闭包还被使用，那么被闭包引用的变量会一直存在。从形式上看，匿名函数都是闭包。

我们看一个例子：

func add(n int) func(int) int {
    sum := n
    f := func(x int) int {
        var i int = 2
        sum += i * x
        return sum
    }
    return f
}

func main() {
    f1 := add(10)
    n11 := f1(3)
    n12 := f1(6)
    f2 := add(20)
    n21 := f2(4)
    n22 := f2(8)
}
该例子中函数变量为f，自由变量为sum，同时f为sum提供绑定的计算环境，使得sum和f粘滞在了一起，它们组成的代码块就是闭包。add函数的返回值是一个闭包，而不仅仅是f函数的地址。在该闭包函数中，只有内部的匿名函数f才能访问局部变量i，而无法通过其他途径访问，因此闭包保证了i的安全性。

当我们分别用不同的参数(10, 20)注入add函数而得到不同的闭包函数变量时，得到的结果是隔离的，也就是说每次调用add函数后都将生成并保存一个新的局部变量sum。

按照命令式语言的规则，add函数只是返回了内嵌函数f的地址，但在执行f函数时将会由于在其作用域内找不到sum变量而出错。而在函数式语言中，当内嵌函数体内引用到体外的变量时，将会把定义时涉及到的引用环境和函数体打包成一个整体（闭包）返回。闭包的使用和正常的函数调用没有区别。

现在我们给出引用环境的定义：在程序执行中的某个点所有处于活跃状态的约束所组成的集合，其中的约束指的是一个变量的名字和其所代表的对象之间的联系。

所以我们说“闭包=函数+引用环境”

当每次调用add函数时都将返回一个新的闭包实例，这些实例之间是隔离的，分别包含调用时不同的引用环境现场。不同于函数，闭包在运行时可以有多个实例，不同的引用环境和相同的函数组合可以产生不同的实例。

其实我们可以将闭包函数看成一个类(C++)，一个闭包函数调用就是实例化一个类，闭包的自由变量就是类的成员变量，闭包函数的参数就是类的函数对象的参数。在该例子中，f1和f2可以看作是实例化的两个对象，ni1和ni2(i=1,2)分别可以看作是函数对象的两次调用（参数不同）的返回值。

这让我们想起了一句名言：对象是附有行为的数据，而闭包是附有数据的行为

说明：C++中的函数对象指的是对象具有函数的功能，即类需要重载运算符“()”

闭包的应用
避免程序运行时异常崩溃
Golang中对于一般的错误处理提供了error接口，对于不可预见的错误（异常）处理提供了两个内置函数panic和recover。error接口类似于C/C++中的错误码，panic和recover类似于C++中的try/catch/throw。

当在一个函数执行过程中调用panic()函数时，正常的函数执行流程将立即终止，但函数中之前使用defer关键字延迟执行的语句将正常展开执行，之后该函数将返回到调用函数，并导致逐层向上执行panic流程，直至所属的goroutine中所有正在执行的函数被终止。错误信息将被报告，包括在调用panic()函数时传入的参数，这个过程称为异常处理流程。

recover函数用于终止错误处理流程。一般情况下，recover应该在一个使用defer关键字的函数中执行以有效截取错误处理流程。如果没有在发生异常的goroutine中明确调用恢复过程（调用recover函数），会导致该goroutine所属的进程打印异常信息后直接退出。

对于第三方库的调用，在不清楚是否有panic的情况下，最好在适配层统一加上recover过程，否则会导致当前进程的异常退出，而这并不是我们所期望的。

简单的实现如下：

func thirdPartyAdaptedHandler(...) {
    defer func() {
        err := recover()
        if err != nil {
            fmt.Println("some exception has happend:", err)
        }
    }()
    ...
}
这个例子比较简单，我们再看一个较复杂的例子：使用闭包，让网站的业务逻辑处理程序更安全地运行。

我们定义了一个名为safeHandler的函数，将所有的业务逻辑处理函数（listHandler、viewHandler和uploadHandler）进行一次包装。safeHandler函数有一个参数并且返回一个值，传入的参数和返回值都是一个函数，且都是http.HandlerFunc类型，这种类型的函数有两个参数：http.ResponseWriter和 *http.Request。事实上，我们正是要把业务逻辑处理函数作为参数传入到safeHandler()方法中，这样任何一个错误处理流程向上回溯的时候，我们都能对其进行拦截处理，从而也能避免程序停止运行。

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if e, ok := recover().(error); ok {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            fmt.Println("WARN: panic in %v - %v", fn, e)
            fmt.Println(string(debug.Stack()))
            }
        }()
        fn(w, r)
    }
}
模板方法
笔者前面写了一篇文章《Template Method in Golang》，通过interface和组合的方式实现了模板方法，我们下面用闭包的方式模拟一下类似于模板方法的简单例子。

定义一个转换函数类型：

type Traveser func(ele interface{})
Process函数功能：对切片array进行了traveser处理

func Process(array interface{}, traveser Traveser) error {
    ...
    traveser(array)
    ...
    return nil
}
SortByAscending函数功能：升序排序数据切片中的数据：

func SortByAscending(ele interface{}) {
    ...
}
SortByDescending函数功能：降序排序数据切片中的数据：

func SortByDescending(ele interface{}) {
    ...
}
Process函数调用：

func main() {
    intSlice := make([]int, 0)
    intSlice = append(intSlice, 3, 1, 4, 2)

    Process(intSlice, SortByDescending)
    fmt.Println(intSlice) //[4 3 2 1]
    Process(intSlice, SortByAscending)
    fmt.Println(intSlice) //[1 2 3 4]
}
模板方法模式是定义一个操作中的算法的框架，而将一些步骤延迟到子类中，使得子类可以不改变一个算法的框架就可重新定义该算法的某些特定步骤。在Golang中，模板方法不但可以通过interface和组合的方式实现，而且可以通过闭包的方式实现。

变量的安全性
闭包内定义的局部变量，只有内部的匿名函数才能访问，而无法通过其他途径访问，这就保证了变量的安全性。

回调函数
闭包经常用于回调函数。当IO操作（例如从网络获取数据、文件读写)完成的时候，会对获取的数据进行某些操作，这些操作可以交给函数对象处理。

小结
闭包的概念从字面上很难理解，特别对于一直使用命令式语言进行编程的读者。本文通过Golang代码进行阐述，澄清了闭包初学者的多个困惑，深入分析了闭包的本质，最后分享了闭包在Golang中的四种应用，希望对读者有一定的帮助。
