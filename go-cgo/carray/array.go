package main

//在go中使用C数组最简单的方法是通过数组将其转换为切片：
//cvFacetSlice := (*[1 << 30]C.cv_face_t)(unsafe.Pointer(face))[:size:size]

// cArray := C.malloc(C.size_t(c_count) * C.size_t(unsafe.Sizeof(uintptr(0))))

// // convert the C array to a Go Array so we can index it
// a := (*[1<<30 - 1]*C.char)(cArray)
// for index, value := range strs {
//     a[index] = C.CString(value)
// }

// team := C.get_team()
// C.free(unsafe.Pointer(team))
// teamSlice := (*[1 << 30]C.team)(unsafe.Pointer(team))[:teamSize:teamSize]

// err := C.f(&c_count, (***C.char)(unsafe.Pointer(&cArray)))
// 1<<30-1, or 1073741823.//最大大小的数组实际上并没有分配，但是Go需要的是常量大小的数组，而且1 << 30将会足够大。 该数组立即转换为一个切片，正确设置长度和容量。
