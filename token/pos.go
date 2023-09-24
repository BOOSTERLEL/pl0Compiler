package token

import "fmt"

// Pos 类似一个指针, 表示文件中的位置.
type Pos int

// NoPos 类似指针的 nil 值, 表示一个无效的位置.
const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}

type Position struct {
	Filename string // 文件名
	Offset   int    // 偏移量, 从 0 开始
	Line     int    // 行号, 从 1 开始
	Column   int    // 列号, 从 1 开始
}

// Position 装行列号位置
func (p Pos) Position(fileName, src string) Position {
	if !p.IsValid() {
		return Position{
			Filename: fileName,
		}
	}

	var pos = Position{
		Filename: fileName,
		Offset:   int(p) - 1,
		Line:     1,
		Column:   1,
	}

	for _, c := range []byte(src[:pos.Offset]) {
		pos.Column++
		if c == '\n' {
			pos.Column = 1
			pos.Line++
		}
	}
	return pos
}

func (pos Position) IsValid() bool {
	return pos.Line > 0
}

// String 打印以下格式:
//
//	file:row:column    	带文件名的有效位置
//	file:row           	有文件名但没有列的有效位置 (column == 0)
//	line:column         没有文件名的有效位置
//	line                没有文件名且没有列的有效位置 (column == 0)
//	file                文件名无效的位置
//	-                   没有文件名的无效位置
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
