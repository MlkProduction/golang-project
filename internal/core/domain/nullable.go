package domain

/*
Nullable needs to specify:
- field not provided (не было передано (пустая строка))
- field provided: null (передан null)
- field provided: value (передано значение)
*/

type Nullable[T any] struct {
	Value *T
	Set   bool
}
