# godif

Go dependency injection for functions (and not only...)

# Usage Example

-  [Service implementation](services/impl_test.go)

# Usage

## Provide func implementation

- Declare: `var toInject func()`
- Register to be injected: `godif.Require(&toInject)`
- Provide implementation: `godif.Provide(&toInject, f)`
- Resolve: `godif.ResolveAll()`
  - Incompatible types -> error
  - More than one implementations provided -> error
  - No implementation -> error
  - Something provided from a package but nothing is required for the package -> error (package is not used)
  - Non-assignable var provided on `Require()` or `Provide()` call -> error, further validation is skipped


## Provide key-value

- Declare: `var MyMap map[string]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string]int{}`
  - Provide implementation: `godif.Provide(&MyMap, map[string]int{})`
- Provide data: `godif.ProvideKeyValue(&MyMap, "key1", 1)`
- Resolve: `godif.ResolveAll()`
  - Data provided but not implemented -> error
  - Use `godif.Provide()` if implemented manually -> error 
  - If implementation provided
    - Multiple implementations -> error
    - Non-map or map of incompatible key or value type -> error
  - Multiple values per key -> error
  - Key or value of different types provided -> error
  - Non-assignable var provided or `ProvideKeyValue()` call -> error, further validation is skipped


## Provide key-slice

- Declare: `var MyMap map[string][]int`
- Require skipped, no error
- Implement
  - Manually: `MyMap = map[string][]int{}`
  - Provide implementation: `godif.Provide(&MyMap, map[string][]int{})`
- Add initial data if needed: `MyMap["key1"] = append(MyMap["key1"], 42)`
  - Further `godif.ProvideKeyValue()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideKeyValue(&MyMap, "key1", 1)`
  - `godif.ProvideKeyValue(&MyMap, "key1", 2)`
  - `godif.ProvideKeyValue(&MyMap, "key1", []int{3, 4})`
- Resolve: `godif.ResolveAll()`
  - Data provided but not implemented -> error
  - Use `godif.Provide()` if implemented manually -> error 
  - If implementation provided
    - Multiple implementations -> error
    - Slice of incompatible element type -> error
  - Non-assignable var provided on `ProvideKeyValue()` call -> error, further validation is skipped


## Provide slice element

- Declare: `var MySlice []string`
- Require skipped, no error
- Implement
  - Manually: `MySlice = []string{}`
  - Provide implementation: `godif.Provide(&MySlice, []string{})`
- Add initial data if needed: `MySlice = append(MySlice, 42)`
  - Further `godif.ProvideSliceElement()` calls will append data to the existing slice
- Provide data: 
  - `godif.ProvideSliceElement(&MySlice, "str1")`
  - `godif.ProvideSliceElement(&MySlice, []string{"str3", "str4"})`
- Resolve: `godif.ResolveAll()`
  - Data provided but not implemented -> error
  - Use `godif.Provide()` if implemented manually -> error 
  - Multiple implementations -> error
  - Incompatible types -> error
  - Non-assignable var provided on `ProvideSliceElement()` call -> error, further validation is skipped

## Reset all injections
- `godif.Reset()`
- Provided vars will be nilled
- Manually inited vars will be kept
- Data injected into manually inited vars will be kept