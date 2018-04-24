int main() {
    extern int sum(int a, int b);
    sum(1,2)
    return 0;
}


// c thread: sum
// ctx = cgo_runtime_init_done()
// runtime/cgo/crosscall2
// goroutine: _cgoexp_xxx_sum
// goroutine: _cgowrap_xxx_sim
// goroutine: sum
// cgo_release_context(ctx)
