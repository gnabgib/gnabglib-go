## To Benchmark hex decode vs go native

- `-run=^#` stops an unit tests from executing
- `-count 50` sets iterations to 50
- `-bench .` tells go to benchmark all found methods (not in old docs it says -bench=.)
`go test -bench . -count 50 -run=^#`

## To run tests
`go test -v`

## Test complexity

`gocyclo.exe .\hex_test.go`