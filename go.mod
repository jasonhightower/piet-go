module jasonhightower.com/piet

go 1.20

replace jasonhightower.com/piet/interpreter => ./interpreter

require (
	github.com/jasonhightower/jcr v0.0.0-20230828064230-30af7aac252c
	github.com/sirupsen/logrus v1.9.3
	jasonhightower.com/piet/interpreter v0.0.0-00010101000000-000000000000
)

require (
	github.com/jasonhightower/bytecode v0.0.0-20230827181703-d8c61652ac2a // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
