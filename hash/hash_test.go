package hash

import (
	"reflect"
	"testing"
)

func TestTokey(t *testing.T) {
	key := Tokey("test1", "test2")

	if key != "e89d66bdfdd4dd26b682cc77e23a86eb" {
		t.Errorf("Tokey fails: %s", key)
	}
}

func TestShard(t *testing.T) {
	tokey := "e89d66bdfdd4dd26b682cc77e23a86eb"

	results := Shard(tokey, 1, 2, false)

	equal := reflect.DeepEqual(results, []string{"e", "8", "e89d66bdfdd4dd26b682cc77e23a86eb"})

	if !equal {
		t.Errorf("Tokey fails: %s", results)
	}

	results = Shard(tokey, 2, 2, false)

	equal = reflect.DeepEqual(results, []string{"e8", "9d", "e89d66bdfdd4dd26b682cc77e23a86eb"})

	if !equal {
		t.Errorf("Tokey fails: %s", results)
	}

	results = Shard(tokey, 2, 3, true)

	equal = reflect.DeepEqual(results, []string{"e8", "9d", "66", "bdfdd4dd26b682cc77e23a86eb"})

	if !equal {
		t.Errorf("Tokey fails: %s", results)
	}
}
