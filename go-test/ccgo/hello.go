/* Package main - transpiled by c2go

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

// Warning (FieldDecl): 133: Cannot resolve type 'int ( * _Nullable)(void *)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *)'
// Warning (FieldDecl): 134: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 135: Cannot resolve type 'fpos_t ( * _Nullable)(void *, fpos_t, int)' : function pointers are not supported [2] : 'fpos_t ( * _Nullable)(void *, fpos_t, int)'
// Warning (FieldDecl): 136: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 133: Cannot resolve type 'int ( * _Nullable)(void *)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *)'
// Warning (FieldDecl): 134: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 135: Cannot resolve type 'fpos_t ( * _Nullable)(void *, fpos_t, int)' : function pointers are not supported [2] : 'fpos_t ( * _Nullable)(void *, fpos_t, int)'
// Warning (FieldDecl): 136: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'

package main

import "github.com/elliotchance/c2go/noarch"
import "github.com/elliotchance/c2go/darwin"

type __int8_t int8
type __uint8_t uint8
type __int16_t int16
type __uint16_t uint16
type __int32_t int
type __uint32_t uint32
type __int64_t int64
type __uint64_t uint64
type __darwin_intptr_t int32
type __darwin_natural_t uint32
type __darwin_ct_rune_t darwin.CtRuneT
type __mbstate_t struct {
	__mbstate8 []byte
	_mbstateL  int64
}
type __darwin_mbstate_t int64
type __darwin_ptrdiff_t int32
type __darwin_size_t uint32
type __darwin_va_list int64
type __darwin_wchar_t int
type __darwin_rune_t int
type __darwin_wint_t int
type __darwin_clock_t uint32
type __darwin_socklen_t uint32
type __darwin_ssize_t int32
type __darwin_time_t int32
type __darwin_blkcnt_t int64
type __darwin_blksize_t int
type __darwin_dev_t int
type __darwin_fsblkcnt_t uint32
type __darwin_fsfilcnt_t uint32
type __darwin_gid_t uint32
type __darwin_id_t uint32
type __darwin_ino64_t uint64
type __darwin_ino_t uint64
type __darwin_mach_port_name_t uint32
type __darwin_mach_port_t uint32
type __darwin_mode_t uint16
type __darwin_off_t int64
type __darwin_pid_t int
type __darwin_sigset_t uint32
type __darwin_suseconds_t int
type __darwin_uid_t uint32
type __darwin_useconds_t uint32
type __darwin_uuid_t []uint8
type __darwin_uuid_string_t []byte
type __darwin_pthread_handler_rec struct {
	__routine func(interface{})
	__arg     interface {
	}
	__next []__darwin_pthread_handler_rec
}
type _opaque_pthread_attr_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_cond_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_condattr_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_mutex_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_mutexattr_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_once_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_rwlock_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_rwlockattr_t struct {
	__sig    int32
	__opaque []byte
}
type _opaque_pthread_t struct {
	__sig           int32
	__cleanup_stack []__darwin_pthread_handler_rec
	__opaque        []byte
}
type __darwin_pthread_attr_t struct {
	__sig           int32
	__cleanup_stack []__darwin_pthread_handler_rec
	__opaque        []byte
}
type __darwin_pthread_cond_t _opaque_pthread_cond_t
type __darwin_pthread_condattr_t _opaque_pthread_condattr_t
type __darwin_pthread_key_t uint32
type __darwin_pthread_mutex_t _opaque_pthread_mutex_t
type __darwin_pthread_mutexattr_t _opaque_pthread_mutexattr_t
type __darwin_pthread_once_t _opaque_pthread_once_t
type __darwin_pthread_rwlock_t _opaque_pthread_rwlock_t
type __darwin_pthread_rwlockattr_t _opaque_pthread_rwlockattr_t
type __darwin_pthread_t []_opaque_pthread_t
type __darwin_nl_item int
type __darwin_wctrans_t int
type __darwin_wctype_t uint32
type int8_t int8
type int16_t int16
type int32_t int
type int64_t int64
type u_int8_t uint8
type u_int16_t uint16
type u_int32_t uint32
type u_int64_t uint64
type register_t int64
type intptr_t int32
type uintptr_t uint32
type user_addr_t uint64
type user_size_t uint64
type user_ssize_t int64
type user_long_t int64
type user_ulong_t uint64
type user_time_t int64
type user_off_t int64
type syscall_arg_t uint64
type va_list int64
type size_t uint32
type fpos_t int64
type __sbuf struct {
	_base []uint8
	_size int
}
type __sFILEX struct {
}
type __sFILE struct {
	_p       []uint8
	_r       int
	_w       int
	_flags   int16
	_file    int16
	_bf      int64
	_lbfsize int
	_cookie  interface {
	}
	_close interface {
	}
	_read interface {
	}
	_seek interface {
	}
	_write interface {
	}
	_ub      int64
	_extra   []__sFILEX
	_ur      int
	_ubuf    []uint8
	_nbuf    []uint8
	_lb      int64
	_blksize int
	_offset  int
}
type FILE struct {
	_p       []uint8
	_r       int
	_w       int
	_flags   int16
	_file    int16
	_bf      int64
	_lbfsize int
	_cookie  interface {
	}
	_close interface {
	}
	_read interface {
	}
	_seek interface {
	}
	_write interface {
	}
	_ub      int64
	_extra   []__sFILEX
	_ur      int
	_ubuf    []uint8
	_nbuf    []uint8
	_lb      int64
	_blksize int
	_offset  int
}

var// Warning (FieldDecl): 133: Cannot resolve type 'int ( * _Nullable)(void *)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *)'
// Warning (FieldDecl): 134: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 135: Cannot resolve type 'fpos_t ( * _Nullable)(void *, fpos_t, int)' : function pointers are not supported [2] : 'fpos_t ( * _Nullable)(void *, fpos_t, int)'
// Warning (FieldDecl): 136: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 133: Cannot resolve type 'int ( * _Nullable)(void *)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *)'
// Warning (FieldDecl): 134: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
// Warning (FieldDecl): 135: Cannot resolve type 'fpos_t ( * _Nullable)(void *, fpos_t, int)' : function pointers are not supported [2] : 'fpos_t ( * _Nullable)(void *, fpos_t, int)'
// Warning (FieldDecl): 136: Cannot resolve type 'int ( * _Nullable)(void *, char *, int)' : function pointers are not supported [2] : 'int ( * _Nullable)(void *, char *, int)'
__stdinp *noarch.File
var __stdoutp *noarch.File

type off_t int64
type ssize_t int32
// main - transpiled function from file : hello.c , line : 3
func main() {
	noarch.Printf([]byte("Hello World\n\x00"))
	return
}
func init() {
	__stdinp = noarch.Stdin
	__stdoutp = noarch.Stdout
}
