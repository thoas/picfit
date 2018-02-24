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
// func RunSliceTests(t *testing.T, fn func(into, from interface{}) error) {
// 	t.Run("Smoke", func(t *testing.T) {
// 		t.Run("IntSliceFromStrings", func(t *testing.T) {
// 			var into []int
// 			exp := []int{12, 345, 6789}
//
// 			err := fn(&into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp (%T) --> %[1]v != %v <-- (%[2]T) got", exp, into)
// 			}
//
// 			err = fn(nil, []string{"12", "345", "6789"})
// 			if err == nil {
// 				t.Error("expected non-nil err")
// 			}
// 		})
// 		t.Run("IntPtrSliceFromStrings", func(t *testing.T) {
// 			var into []*int
// 			i1, i2, i3 := new(int), new(int), new(int)
// 			*i1, *i2, *i3 = 12, 345, 6789
// 			exp := []*int{i1, i2, i3}
//
// 			err := fn(&into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp --> (%T) %#[1]v != %T %#[2]v <-- got", exp, into)
// 			}
//
// 			into = []*int{}
// 			err = fn(&into, []string{"12", "345", "6789"})
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if !reflect.DeepEqual(exp, into) {
// 				t.Fatalf("exp --> (%T) %#[1]v != %T %#[2]v <-- got", exp, into)
// 			}
// 		})
// 	})
//
// 	// tests all supported sources
// 	for _, test := range generated.NewSliceTests() {
// 		into, from, exp := test.Into, test.From, test.Exp
//
// 		name := fmt.Sprintf(`From(%T)/Into(%T)`, from, into)
// 		t.Run(name, func(t *testing.T) {
// 			err := fn(into, from)
// 			if err != nil {
// 				t.Error(err)
// 			}
//
// 			if !reflect.DeepEqual(exp, refutil.Indirect(into)) {
// 				t.Logf("from (%T) --> %[1]v", from)
// 				t.Fatalf("\nexp (%T) --> %[1]v\ngot (%[2]T) --> %[2]v", exp, into)
// 			}
// 		})
// 	}
// }
//
// // Summary:
// //
// // BenchmarkSlice/<slice size>/<from> to <to>/Conv:
// //   Measures the most convenient form of conversion using this library.
// //
// // BenchmarkSlice/<slice size>/<from> to <to>/Conv:
// //   Measures using the library only for the conversion, looping for apending.
// //
// // BenchmarkSlice/<slice size>/<from> to <to>/Conv:
// //   Measures not using this library at all, pure Go implementation.
// //
// func RunSliceBenchmarks(b *testing.B, fn func(into, from interface{}) error) {
// 	for _, num := range []int{1024, 64, 16, 4} {
// 		num := num
//
// 		// slow down is really tolerable, only a factor of 1-3 tops
// 		b.Run(fmt.Sprintf("Length(%d)", num), func(b *testing.B) {
//
// 			b.Run("[]string to []int64", func(b *testing.B) {
// 				strs := make([]string, num)
// 				for n := 0; n < num; n++ {
// 					strs[n] = fmt.Sprintf("%v00", n)
// 				}
// 				b.ResetTimer()
//
// 				b.Run("Conv", func(b *testing.B) {
// 					for i := 0; i < b.N; i++ {
// 						var into []int64
// 						err := fn(&into, strs)
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
// 						var into []int64
//
// 						for _, s := range strs {
// 							v, err := strconv.ParseInt(s, 10, 0)
// 							if err != nil {
// 								b.Error(err)
// 							}
// 							into = append(into, v)
// 						}
// 						if len(into) != num {
// 							b.Error("bad impl")
// 						}
// 					}
// 				})
// 			})
//
// 			b.Run("[]string to []*int64", func(b *testing.B) {
// 				strs := make([]string, num)
// 				for n := 0; n < num; n++ {
// 					strs[n] = fmt.Sprintf("%v00", n)
// 				}
// 				b.ResetTimer()
//
// 				b.Run("Library", func(b *testing.B) {
// 					for i := 0; i < b.N; i++ {
// 						var into []*int64
// 						err := fn(&into, strs)
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
// 						into := new([]*int64)
//
// 						for _, s := range strs {
// 							v, err := strconv.ParseInt(s, 10, 0)
// 							if err != nil {
// 								b.Error(err)
// 							}
// 							*into = append(*into, &v)
// 						}
// 						if len(*into) != num {
// 							b.Error("bad impl")
// 						}
// 					}
// 				})
// 			})
// 		})
// 	}
// }
