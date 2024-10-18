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

### URI Params

#### by struct
```
params := struct{
    Foo string `uri:"foo"`
    Bar string `uri:"bar"`
}{
    Foo: "abc",
    Bar: "xyz",
}

builder := GinTestContext.NewBuilder()
builder.SetURIParams(params)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

#### by map
```
params := map[string]string
params["foo"] = "abc"
params["bar"] = "xyz"

builder := GinTestContext.NewBuilder()
builder.SetURIParams(paramsß)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

### Queries

#### by struct
```
queries := struct{
    Foo string `form:"foo"`
    Bar string `form:"bar"`
}{
    Foo: "abc",
    Bar: "xyz",
}

builder := GinTestContext.NewBuilder()
builder.SetQueries(queries)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

#### by map
```
queries := map[string]string
queries["foo"] = "abc"
queries["bar"] = "xyz"

builder := GinTestContext.NewBuilder()
builder.SetQueries(queries)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

### Body

#### by struct
```
body := struct{
    Foo string `json:"foo"`
    Bar string `json:"bar"`
}{
    Foo: "abc",
    Bar: "xyz",
}

builder := GinTestContext.NewBuilder()
builder.SetBody(body)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```

#### by map
```
body := map[string]string
body["foo"] = "abc"
body["bar"] = "xyz"

builder := GinTestContext.NewBuilder()
builder.SetBody(body)

context, err := builder.GetContext()
require.NoError(t, err)

yourHandler(context)
```