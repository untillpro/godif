# godif

Go dependency injection for functions (and not only...)

# Usage Example

-  [Interface declaration](iservices/README.md)
-  [Interface implementation](services/declare.go)

# Usage

## Provide key-value

- Declare: `var MyMap map[string]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string]int{}`
    - Use Provide() -> error
  - Provide implementation: `godif.Provide(&MyMap, map[string]int{})`
    - Multiple implementations -> error
    - Non-map or map of incompatible key or value type -> error
  - No implementation -> error
- Provide data: `godif.ProvideKeyValue(&MyMap, "key1", 1)`
  - Multiple values per key -> error
  - Key or value of different types provided -> error

## Provide key-slice

- Declare: `var MyMap map[string][]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string][]int{}`
    - Use `Provide()` -> error
  - Provide implementation: `godif.Provide(&MyMap, map[string][]int{})`
    - Multiple implementations -> error
    - slice of incompatible element type -> error
  - No implementation -> error
- Add initial data if needed: `MyMap["key1"] = append(MyMap["key1"], 42)`
  - Further `godif.ProvideKeyValue()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideKeyValue(&MyMap, "key1", 1)`
  - `godif.ProvideKeyValue(&MyMap, "key1", 2)`
  - `godif.ProvideKeyValue(&MyMap, "key1", []int{3, 4})`

## Provide slice element

- Declare: `var MySlice []string`
- Implement
  - Manually: `MySlice = []string{}`
    - Use `Provide()` -> error
  - Provide implementation: `godif.Provide(&MySlice, []string{})`
    - Multiple implementations -> error
    - Incompatible types -> error
  - No implementation -> error
- Add initial data if needed: `MySlice = append(MySlice, 42)`
  - Further `godif.ProvideSliceElement()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideSliceElement(&MySlice, "str1")`
  - `godif.ProvideSliceElement(&MySlice, []string{"str3", "str4"})`

## Reset all injections
- `godif.Reset()`
- Provided vars will be nilled
- Manually inited vars will be kept
- Data injected into manually inited vars will be kept