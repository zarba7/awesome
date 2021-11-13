module domain_role

go 1.17


require (
	ddd v0.0.0
	internal v0.0.0
)
replace (
	internal => ../../internal
	ddd => ../../pkg/ddd
)