package projection

import (
	"testing"
	"reflect"
)


func TestBasicWR(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Commit()

	v,_ := prj.Get("a.b")
	ve := 1
	if v != ve {
		t.Error("was not 1")
	}


	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["b"] = 1
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'b': 1}")
	}

}

func TestBasicWD(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Upsert("a.c", 2)
	prj.Commit()

	prj.Delete("a.b")
	prj.Commit()

	v,_ := prj.Get("a.b")
	if v != nil {
		t.Error("was not nil")
	}

	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["c"] = 2
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'c': 2}")
	}

	prj.Delete("a.c")
	prj.Commit()

	vo2,_ := prj.Get("a")
	if vo2 != nil {
		t.Error("was not nil")
	}
}

func TestSubtreeDelete(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Upsert("a.c.1", 2)
	prj.Upsert("a.c.2", 2)
	prj.Commit()

	prj.Delete("a.c")
	prj.Commit()

	v,_ := prj.Get("a.c")
	if v != nil {
		t.Error("was not nil")
	}

	v2,_ := prj.Get("a.c.1")
	if v2 != nil {
		t.Error("was not nil")
	}

	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["b"] = 1
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'b': 1}")
	}
}

func TestSubtreeUpsert(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Upsert("a.c.1", 2)
	prj.Upsert("a.c.2", 2)
	prj.Commit()

	prj.Upsert("a.c", 2)
	prj.Commit()

	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["b"] = 1
	voe["c"] = 2
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'b': 1, 'c': 2}")
	}
}

func TestSubtreeUpsertRestore(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Upsert("a.c.1", 2)
	prj.Upsert("a.c.2", 2)
	prj.Commit()

	prj.Upsert("a.c", 2)
	prj.Commit()

	prj.Upsert("a.c.2", 2)
	prj.Commit()

	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["b"] = 1
	voesubc := make(map[string]interface{})
	voesubc["2"] = 2
	voe["c"] = voesubc
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'b': 1, 'c': {'2': 2}}")
	}
}

func TestSubtreeUpsertRestoreRollback(t *testing.T) {

	prj := NewProjection()

	prj.Upsert("a.b", 1)
	prj.Upsert("a.c.1", 2)
	prj.Upsert("a.c.2", 2)
	prj.Commit()

	prj.Upsert("a.c", 2)
	prj.Commit()

	prj.Upsert("a.c.2", 2)

	vo,_ := prj.Get("a")
	voe := make(map[string]interface{})
	voe["b"] = 1
	voesubc := make(map[string]interface{})
	voesubc["2"] = 2
	voe["c"] = voesubc
	if !reflect.DeepEqual(vo, voe) {
		t.Error("was not {'b': 1, 'c': {'2': 2}}")
	}

	prj.Rollback()

	vo2,_ := prj.Get("a")
	voe2 := make(map[string]interface{})
	voe2["b"] = 1
	voe2["c"] = 2
	if !reflect.DeepEqual(vo2, voe2) {
		t.Error("was not {'b': 1, 'c': 2}")
	}

}