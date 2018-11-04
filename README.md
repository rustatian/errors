## Errors
This is a project to work with errors in a nice way. Motivated by https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html
Errors fully compatible with the std errors

##Sample of a usage
1. Clone this repo into your project
2. Create the folder `errors` and put this code into that folder
3. Define your own error kinds and update tests a little bit
4. Profit ++

You can use this with different build tags also:
1. `debug` is used to get the call stack in case of error with operations (if they defined)
2. `w/o tags` is used to see errors and operations only (`Op` data type), w/o stack trace

Examples:
1. Using with operations:
```Go
func foo() error {
    const Op = "foo"
    err := someOperation()
    return E(Op, err)
}
```

2. Using with kinds of errors:
```Go
func foo() error {
	err := someFunc() // <- select from database
}

func someFunc() (int, error) {
	const Op = "GET"
	// GET from database...
	err := getFromDatabase // <- just to simplify
	if err != nil {
		return E(Op, err, DataBaseFail) // <- DataBaseFail is error kind defined in errors.go
	}
	return 
}

// As the result you will see error with error code and operations. 
// So, you will find the place where this error occurred easily and fast 

```

Valery Piashchynski, 
SpiralScout.com
