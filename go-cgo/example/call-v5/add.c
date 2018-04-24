#include "add.h"
#include "_cgo_export.h"

int c_add(int a, int b) {
    return GoAdd(a,b);
}

//在C文件中使用_cgo_export.h头文件
//C文件必须在同一个包，否则会找不到头文件
