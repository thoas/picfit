package testconv

// import (
// 	"fmt"
// 	"reflect"
// 	"strconv"
// 	"testing"
//
// 	"github.com/cstockton/go-conv/internal/generated"
// 	"github.com/cstockton/go-conv/internal/refutil"
// )
//
// func RunMapTests(t *testing.T, fn func(into, from interface{}) error) {
// 	t.Run("Smoke", func(t *testing.T) {
// 		t.Run("StringIntMapFromStrings", func(t *testing.T) {
// 			into := make(map[string]int64)
// 			err := fn(into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
//
// 			exp := map[string]int64{"0": 12, "1": 345, "2": 6789}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp (%T) --> %[1]v != %v <-- (%[2]T) got", exp, into)
// 			}
//
// 			err = fn(nil, []string{"12", "345", "6789"})
// 			if err == nil {
// 				t.Error("expected non-nil err")
// 			}
// 		})
// 		t.Run("StringIntPtrMapFromStrings", func(t *testing.T) {
// 			i1, i2, i3 := new(int64), new(int64), new(int64)
// 			*i1, *i2, *i3 = 12, 345, 6789
// 			exp := map[string]*int64{"0": i1, "1": i2, "2": i3}
//
// 			into := make(map[string]*int64)
// 			err := fn(into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp (%T) --> %[1]v != %v <-- (%[2]T) got", exp, into)
// 			}
// 			into = make(map[string]*int64)
// 			err = fn(into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp (%T) --> %[1]v != %v <-- (%[2]T) got", exp, into)
// 			}
// 		})
// 	})
//
// 	// tests all supported sources
// 	for _, test := range generated.NewMapTests() {
// 		from, exp := test.From, test.Exp
// 		run := func(into interface{}) {
// 			name := fmt.Sprintf(`From(%T)/Into(%T)`, from, into)
// 			t.Run(name, func(t *testing.T) {
// 				err := fn(into, from)
// 				if err != nil {
// 					t.Error(err)
// 				}
//
// 				if !reflect.DeepEqual(exp, refutil.Indirect(into)) {
// 					t.Logf("from (%T) --> %[1]v", from)
// 					t.Fatalf("\nexp (%T) --> %[1]v\ngot (%[2]T) --> %[2]v", exp, into)
// 				}
// 			})
// 		}
//
// 		// Test the normal type
// 		run(test.Into)
//
// 		typ := reflect.TypeOf(test.Into)
// 		mapVal, mapValPtr := reflect.MakeMap(typ), reflect.New(typ)
// 		mapValPtr.Elem().Set(mapVal)
//
// 		// Ensure pointer to a map works as well.
// 		run(mapValPtr.Interface())
// 	}
// }
//
// // Summary: Not much of a tax here, about 2x as slow.
// //
// // BenchmarkMap/<slice size>/<from> to <to>/Conv:
// //   Measures the most convenient form of conversion using this library.
// //
// // BenchmarkSlice/<slice size>/<from> to <to>/Conv:
// //   Measures using the library only for the conversion, looping for apending.
// //
// // BenchmarkSlice/<slice size>/<from> to <to>/Conv:
// //   Measures not using this library at all, pure Go implementation.
// //
// // BenchmarkMap/Length(1024)/[]string_to_map[int]string/Conv-24    	    1000	   1321364 ns/op
// // BenchmarkMap/Length(1024)/[]string_to_map[int]string/LoopConv-24         	    2000	    896001 ns/op
// // BenchmarkMap/Length(1024)/[]string_to_map[int]string/LoopStdlib-24       	    2000	    652117 ns/op
// // BenchmarkMap/Length(64)/[]string_to_map[int]string/Conv-24               	   20000	     74431 ns/op
// // BenchmarkMap/Length(64)/[]string_to_map[int]string/LoopConv-24           	   20000	     56702 ns/op
// // BenchmarkMap/Length(64)/[]string_to_map[int]string/LoopStdlib-24         	   30000	     44191 ns/op
// // BenchmarkMap/Length(16)/[]string_to_map[int]string/Conv-24               	  100000	     18422 ns/op
// // BenchmarkMap/Length(16)/[]string_to_map[int]string/LoopConv-24           	  100000	     14193 ns/op
// // BenchmarkMap/Length(16)/[]string_to_map[int]string/LoopStdlib-24         	  200000	     10021 ns/op
// // BenchmarkMap/Length(4)/[]string_to_map[int]string/Conv-24                	  300000	      4402 ns/op
// // BenchmarkMap/Length(4)/[]string_to_map[int]string/LoopConv-24            	  500000	      2783 ns/op
// // BenchmarkMap/Length(4)/[]string_to_map[int]string/LoopStdlib-24          	 1000000	      1986 ns/op
// func RunMapBenchmarks(b *testing.B, fn func(into, from interface{}) error) {
// 	for _, num := range []int{1024, 64, 16, 4} {
// 		num := num
//
// 		// slow down is really tolerable, only a factor of 1-3 tops
// 		b.Run(fmt.Sprintf("Length(%d)", num), func(b *testing.B) {
//
// 			b.Run("[]string to map[int]string", func(b *testing.B) {
// 				strs := make([]string, num)
// 				for n := 0; n < num; n++ {
// 					strs[n] = fmt.Sprintf("%v00", n)
// 				}
// 				b.ResetTimer()
//
// 				b.Run("Conv", func(b *testing.B) {
// 					for i := 0; i < b.N; i++ {
// 						into := make(map[string]int64)
// 						err := fn(into, strs)
// 						if err != nil {
// 							b.Error(err)
// 						}
// 						if len(into) != num {
// 							b.Error("bad impl")
// 						}
// 					}
// 				})
// 				b.Run("Stdlib", func(b *testing.B) {
// 					for i := 0; i < b.N; i++ {
// 						into := make(map[string]int64)
//
// 						for seq, s := range strs {
// 							k := fmt.Sprintf("%v", seq)
// 							v, err := strconv.ParseInt(s, 10, 0)
// 							if err != nil {
// 								b.Error(err)
// 							}
// 							into[k] = v
// 						}
// 						if len(into) != num {
// 							b.Error("bad impl")
// 						}
// 					}
// 				})
// 			})
// 		})
// 	}
// }
