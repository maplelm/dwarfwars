package game

type OOB struct {
	Details string
}

func (b OOB) Error() string {
	if b.Details == "" {
		return "Out Of Bounds"
	}
	return "Out Of Bounds: " + b.Details
}
