package customerrors

type DuplicateVersionError struct {
}

func (d *DuplicateVersionError) Error() string {
	return "DUPLICATED VERSION ERROR"
}
