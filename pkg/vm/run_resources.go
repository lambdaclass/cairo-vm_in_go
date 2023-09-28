package vm

type RunResources struct {
	nSteps *uint
}

func NewRunResources(nSteps uint) RunResources {
	return RunResources{nSteps: &nSteps}
}

func (r *RunResources) Consumed() bool {
	return r.nSteps != nil && *r.nSteps == 0
}

func (r *RunResources) ConsumeStep() {
	if r.nSteps != nil && *r.nSteps != 0 {
		*r.nSteps--
	}
}

func (r *RunResources) GetNSteps() *uint {
	return r.nSteps
}
