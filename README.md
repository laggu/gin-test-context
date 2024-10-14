# GinTestContext
Test your Gin handler easily!

## Install
```
go get -u github.com/laggu/GinTestContext
```

## Features
GinTestContext makes gin.Context which is set with various variables you need

### What you can do
* setting headers
* setting uri params
* setting queries
* setting body

## Examples

### Headers

#### by struct
```
header := struct{
    Foo string `header:"foo"`
    Bar string `header:"bar"`
}{
    Foo: "abc",
    Bar: "xyz",
}

builder := GinTestContext.NewBuilder()
builder.SetHeaders(header)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

#### by map
```
header := map[string]string
header["foo"] = "abc"
header["bar"] = "xyz"

builder := GinTestContext.NewBuilder()
builder.SetHeaders(header)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```
