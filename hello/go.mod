module example.com/hello

go 1.15

require (
	example.com/greetings v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20220325121720-054d8573a5d8
)

replace example.com/greetings => ../greetings
