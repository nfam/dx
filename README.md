# Syntax

```yaml
exp:
    ( "//(any)":     string                                  )?
    ( "has":         string | [string] | [string | [string]] )?
    ( "slice":       slice | [slice] | [slice | [slice]]     )?
    (
        ( "next":    exp | [exp]   )|
        ( "value":   any           )|
        ( "list":    list          )|
        ( "dict":    dict          )|
        ( "convert": string | time )
    )?
slice:
    ( "//(any)":     string                                  )?
    ( "after":       string | [string] | [string | [string]] )?
    ( "before":      string | [string] | [string | [string]] )?
    ( "beforeLast":  string | [string] | [string | [string]] )?
    ( "afterLast":   string | [string] | [string | [string]] )?
    ( "join":       [string | join | [join]]                 )?
    ( "func":        string | [string]                       )?
join:
    ( "//(any)":     string                                  )?
    ( "has":         string | [string] | [string | [string]] )?
    ( "slice":       slice | [slice] | [slice | [slice]]     )?
    ( "value":       string                                  )?
list:
    ( "//(any)":     string                                  )?
    ( "sep":         string | [string]                       )?
    ( "item":        exp | [exp]                             )?
dict:
    ( $name:         exp | [exp | null]                      )*
time:
      "name":        "timeUnix" | "timeUnixMilli"
      "layout":      string | [string]
    ( "location":    number | string                         )?
```

## has string
```
has:(value)
not:(value)
(value)
```

## Built-in func
```
trim
trimPrefix:(prefix)
trimSuffix:(suffix)
trimLeft:(cutset)
trimRight:(cutset)
```

## Built-in convert
```
atoi
atof
```
