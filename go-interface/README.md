### 要点：
##### 接口是系列方法的集合
##### 单个类型可以实现多个方法
##### 同一个接口可以被多种类型实现
##### 接口声明可内嵌其他接口，导入内嵌接口的所有方法(可导出方法+不可导出方法)，多层内嵌接口也将全部被导入到接口声明中
##### 禁止接口的循环内嵌
##### 接口内方法名必须唯一：自定义方法和内嵌接口包含方法，名称必须唯一
##### 接口变量可以保存所有实现了此接口的所有类型的值：抽象的理论实现
##### 静态类型VS动态类型：接口类型的变量可以被实现了其接口的类型间相互赋值，动态类型
##### 接口类型变量：动态类型、动态值，只有二者均为零值nil时此接口类型的变量才为nil
##### 空接口：可以承载任何类型的变量，也可以说任何类型都实现(满足)了空接口
##### 接口实现：类型定义了包含某接口声明的所有方法(方法名+签名一致)
##### 接口类型值只能访问接口自身定义的方法，原类型的其他变量无法访问：行为的抽象，联想对比Java中子类赋值给父类时多态特性
##### 关键词：接口(interface)，类型(type)，方法(method)，函数(function)，方法签名(signature)，可导出方法(exported method)，满足(satisfy)，实现(implement)，

##### 接口(Interface)使得代码更具灵活性、扩展性，是Golang中多态的实现方式。接口允许指定只有一些行为是需要的，而不需要指定一种特定类型。
##### 行为的定义就是一系列方法的集合

type I interface {
    f1(name string)
    f2(name string) (error, float32)
    f3() int64
}

## 在Golang中，关于接口的两个概念：

#### 1.接口(Interface)：实现此接口必须的系列方法，通过关键词interface定义。
#### 2.接口类型(Interface type)：接口类型变量可以保存实现了任何特定接口的类型值。

#### 接口定义：方法、内嵌其他接口
##### 接口的声明指定了属于此接口的方法，方法定义则通过其名称和签名(输入和返回参数)完成
##### 接口中除了包含方法，也允许内嵌(embedded)其他接口-同一个包内定义或已被导入，此时接口将内嵌接口的所有方法导入到自己的定义中

import "fmt"
type I interface {
     m1()
}
type J interface {
    m2()
    I
    fmt.Stringer
}

#### 接口J包含的方法集为：

##### m1() - 来自内嵌的接口I
##### m2() - 自定义方法
##### String() - 来自fmt.Stringer接口的String()方法(此接口中只有这一个方法)
##### 接口内方法顺序不关紧要，所以接口内可能出现方法和内嵌接口的交错出现的情况。

### 接口将导入内嵌接口的所有方法，包含可导出方法(首字母大写的方法)和非导出方法(首字母小写)。

#### 内嵌接口多层内嵌时，导入所有包含接口的方法
##### 如果接口I嵌入了接口J，而接口又内嵌了其他接口K，则接口K的所有方法也将被加入到接口I的声明中

type I interface {
    J
    i()
}
type J interface {
    K
    j()
}
type K interface {
    k()
}
##### 则接口I包含的方法为：i(), j(), k()

### 禁止接口的循环内嵌
#### 接口的循环内嵌是被禁止的，正在编译阶段将被检测出错误，如下代码将产生error错误：interface type loop involving I

type I interface {
    J
    i()
}
type J interface {
    K
    j()
}
type K interface {
    k()
    I
}

### 接口内的方法名必须唯一
#### 如下代码中自定义方法和内嵌接口包含方法存在命名冲突，则将会抛出编译时错误error：duplicate method i：

type I interface {
    J
    i()
}
type J interface {
    j()
    i(int)
}
#### 接口的组装形式贯穿标准库的各种定义，一个io.ReaderWriter的例子：

type ReadWriter interface {
    Reader
    Writer
}

## 接口类型的变量
#### 接口I类型的变量可以保存任何实现了接口I的值

type I interface {
    method1()
}
type T struct{}
func (T) method1() {}
func main() {
    var i I = T{} // 变量i是接口类型I的变量, T实现了接口I, 则i可以保存类型为T的变量值
    fmt.Println(i)
}

#### 静态类型VS动态类型
#### 变量已有的类型在编译阶段被明确，在声明时指定变量类型，且永不改变，这种情况称为静态类型(static type)或直接简称类型。
#### 接口类型的变量也有一种静态的类型，就是接口自身；他们额外还拥有动态类型(dynamic typ

####  变量i的静态类型是接口I，这是不会改变的。
####  另一方面，动态类型也就是动态变化的，第一次赋值后，i的动态类型为类型T1，然而这并不是固化的，第二次赋值将i的动态类型修改为类型T2。
####  当接口类型的变量值为空值nil时(接口的零值为nil)，则动态类型未被设置。

### 如何获取接口类型变量的动态类型

##### reflect包提供获取动态类型的方法，示例代码:
##### 当变量为零值nil时，reflect包将报错runtime error。

###### 1.fmt.Println(reflect.TypeOf(i).PkgPath(), reflect.TypeOf(i).Name())
###### 2.fmt.Println(reflect.TypeOf(i).String())
##### fmt包通过格式化动词%T也可以获取变量的动态类型：
##### 虽然fmt也是使用reflect来实现的，但当变量i为零值nil时也可以支持。
###### 1.fmt.Printf("%T\n", i)

### 接口类型空值nil

type I interface {
    M()
}
type T struct {}
func (T) M() {} // 类型T实现了接口I
func main() {
    var t *T // 变量t必然是空值nil
    if t == nil {
        fmt.Println("t is nil")  // 输出这里
    } else {
        fmt.Println("t is not nil")
    }
    var i I = t // t是空值，但i呢？
    if i == nil {
        fmt.Println("i is nil")
    } else {
        fmt.Println("i is not nil")
    }
}
t is nil
i is not nil

##### 这个输出结果显然有些吃惊，赋值给变量i的值是nil，但i并不等于nil，接口类型变量包含两个部分：

###### 动态类型(dynamic type)
###### 动态值(dynamic value)

###### 动态值是实际变量实际被赋值的值，在上面的例子中var i I = t，变量i的动态值是nil，但i的动态类型是*T。

###### 通过fmt.Printf("%T\n", i)输出赋值后的变量i的动态类型为*main.T，接口类型变量为nil当且仅当动态类型和动态值均为nil。
###### 这种情况下，即使接口类型的变量保存了一个空值指针(nil pointer)，但这个接口变量并不是nil。
##### 已知的错误是从应该返回接口类型的函数返回未初始化、非接口类型的变量值

type I interface {}
type T struct {}
func F() I { // 函数F应该返回接口类型I，在此例中接口类型的返回值=返回类型为*T，值为nil
    var t *T
    if false { // not reachable but it actually sets value
        t = &T{}
    }
    return t // 这里返回的变量t是空值nil
}
func main() {
    fmt.Printf("F() = %v\n", F()) // 返回参数的动态值为nil
    fmt.Printf("F() is nil: %v\n", F() == nil) // 返回参数为接口类型，此接口类型的值并不为nil
    fmt.Printf("type of F(): %T", F()) // 返回参数类型为类型T
}

F() = <nil>
F() is nil: false
type of F(): *main.T

##### 函数返回的接口类型值的动态类型为*main.T，所以它不等于空值nil

### 空的接口
##### 接口的方法集可以完全为空


type I interface {}
type T struct {}
func (T) M() {}
func main() {
    var i I = T{}
    _ = i
}

##### 空接口自动被任何类型实现，所以任何类型都可以赋值给这种空接口类型的变量。
##### 空接口的动态或静态类型的行为和非空接口一致。
##### fmt.Println函数中可变参数中空接口广泛使用。

### 实现接口
#### 实现方法所有方法的任何类型都自动满足(实现)这个接口，不需要想Java中那样显示声明类型实现了哪个接口。
##### Go编译器自动检测类型对接口的实现，这是Golang语言级的强大特性。
import (
    "fmt"
    "regexp"
)
type I interface {
    Find(b []byte) []byte // 接口I包含方法Find，而Regexp实现包含此方法实现(方法名+签名)，则Regexp实现了接口I
}
func f(i I) {
    fmt.Printf("%s\n", i.Find([]byte("abc")))
}
func main() {
    var re = regexp.MustCompile(`b`) // 返回类型为*Regexp
    f(re)
}
##### 这里我们定义了一个接口I，在没有修改内置的regexp模块的情况下，使得：regexp.Regexp类型实现了接口I。

###### 一个类型可以实现多个接口，一个接口可以被多个类型实现
###### 一个接口实现某接口，赋值给接口类型后，只能访问接口自身定义的方法

### 接口类型行为的抽象
#### 接口类型值只能访问此接口类型的方法，其隐藏原类型中包含的其他值，比如结构体、数组、scalar等等。
type I interface {
    M1()
}
type T int64
func (T) M1() {}
func (T) M2() {}
func main() {
    var i I = T(10) // 接口类型值i只能访问M1()方法
    i.M1()
    i.M2() // i.M2 undefined (type I has no field or method M2)
}
