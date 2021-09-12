package misc

// Status mimic Mongoose's Update, Delete status
type Status struct {
	IsOk bool `json:"ok"`
}

func (ds *Status) SetStatus(isOk bool) {
	ds.IsOk = isOk
}
