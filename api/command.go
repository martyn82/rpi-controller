package api

type ICommand interface {
    String() string
    DeviceName() string
    PropertyName() string
    PropertyValue() string
}
