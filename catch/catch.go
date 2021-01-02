package catch

// PENDING
const (
	PENDING   int = 0
	FULFILLED int = 1
	REJECTED  int = 2
	FINALLY   int = 3
)

// Exec struct
type Exec struct {
	status int
	ret    interface{}
}

// Exception interface
type Exception interface {
	Then(func(interface{})) Exception
	Catch(func(err interface{})) Exception
	Finally(func()) Exception
}

func (exec *Exec) deferHelper(cb func()) {
	defer func() {
		defer func() {
			r := recover()
			if r != nil && exec.status != FINALLY {
				exec.status = REJECTED
				exec.ret = r

			}
		}()

		if exec.status == PENDING {
			exec.status = FULFILLED
		}
		cb()
	}()
}

// Try func
func Try(cb func() interface{}) Exception {
	exec := &Exec{
		status: PENDING,
	}

	exec.deferHelper(func() {
		exec.ret = cb()
	})

	return exec
}

// Then func
func (exec *Exec) Then(cb func(interface{})) Exception {
	if exec.status == FULFILLED {
		exec.deferHelper(func() {
			cb(exec.ret)
		})
	}
	return exec
}

// Catch func
func (exec *Exec) Catch(cb func(err interface{})) Exception {
	if exec.status == REJECTED {
		exec.deferHelper(func() {
			cb(exec.ret)
		})
	}
	return exec
}

// Finally func
func (exec *Exec) Finally(cb func()) Exception {
	if exec.status != FINALLY {
		exec.status = FINALLY
		exec.deferHelper(cb)
	}

	return exec
}
