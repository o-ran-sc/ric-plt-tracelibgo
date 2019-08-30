module tracetest

go 1.12

replace tracelibgo => ../

require (
	gerrit.o-ran-sc.org/r/com/golog.git v0.0.1 // indirect
	tracelibgo v0.0.0-00010101000000-000000000000 // indirect
)

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.1
