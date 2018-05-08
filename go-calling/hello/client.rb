require 'ffi'

# Module that represents shared lib
module Hello
  extend FFI::Library

  ffi_lib './hello.so'

  # define class GoSlice to map to:
  # C type struct { void *data; GoInt len; GoInt cap; }
  class GoSlice < FFI::Struct
    layout :data,  :pointer,
           :len,   :long_long,
           :cap,   :long_long
  end

  # define class GoString to map:
  # C type struct { const char *p; GoInt n; }
  class GoString < FFI::Struct
    layout :p,     :pointer,
           :len,   :long_long
  end

  # foreign function definitions
  attach_function :Add, [:long_long, :long_long], :long_long
  attach_function :Cosine, [:double], :double
  attach_function :Sort, [GoSlice.by_value], :void
  attach_function :Log, [GoString.by_value], :int
end

# Call Add
print "hello.Add(12, 99) = ",  hello.Add(12, 99), "\n"

# Call Cosine
print "hello.Cosine(1) = ", hello.Cosine(1), "\n"

# call Sort
nums = [92,101,3,44,7]
ptr = FFI::MemoryPointer.new :long_long, nums.size
ptr.write_array_of_long_long  nums
slice = hello::GoSlice.new
slice[:data] = ptr
slice[:len] = nums.size
slice[:cap] = nums.size
hello.Sort(slice)
sorted = slice[:data].read_array_of_long_long nums.size
print "hello.Sort(", nums, ") = ", sorted, "\n"

# Call Log
msg = "Hello Ruby!"
gostr = hello::GoString.new
gostr[:p] = FFI::MemoryPointer.from_string(msg)
gostr[:len] = msg.size
print "logid ", hello.Log(gostr), "\n"
