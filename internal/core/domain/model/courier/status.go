package courier

const (
	StatusEmpty Status = ""
	StatusFree  Status = "free"
	StatusBusy  Status = "busy"
)

type Status string

func (s Status) Equals(other Status) bool {
	return s == other
}

func (s Status) IsEmpty() bool {
	return s == StatusEmpty
}

func (s Status) String() string {
	return string(s)
}
