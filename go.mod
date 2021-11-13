module awesome

go 1.17

require (
	ddd v0.0.0
)

replace (
	ddd => ./pkg/ddd
)