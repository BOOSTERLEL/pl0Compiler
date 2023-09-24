package builtin

import _ "embed"

//go:embed _builtin.ll
var llBuiltin string

// llBuiltin_wasm go:embed _builtin_wasm.ll
var llBuiltin_wasm string

func GetBuiltinLL(goos, goarch string) string {
	switch goos {
	case "wasm":
		return llBuiltin_wasm
	case "darwin":
	case "linux":
	case "windows":
	}
	return llBuiltin
}

const Header = `
declare i32 @pl_0_builtin_exit(i32)
declare i32 @pl_0_builtin_println(i32)
declare i32 @pl_0_builtin_write()

`

const MainMain = `
define i32 @main() {
	call i32() @pl_0_main()
	ret i32 0
}
`
