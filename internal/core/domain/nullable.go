package domain

type NullableString struct {
	Value *string
	Set bool
}


type Nullable[T any] struct {
	Value *T
	Set bool
}