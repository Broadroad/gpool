package connpool

import "github.com/gpool/errors"

var (
	// NilERROR means the ConnPool is nil
	NilERROR = errors.ConnPoolError{"GPool is nil", 10000}
	// FullERROR means the ConnPool is full
	FullERROR = errors.ConnPoolError{"GPool is full", 20000}
	// FillERROR means the ConnPool can't fill
	FillERROR = errors.ConnPoolError{"GPool can't fill", 30000}
	// ParameterERROR means parameter is illgle
	ParameterERROR = errors.ConnPoolError{"Parameter is illagle", 40000}
)
