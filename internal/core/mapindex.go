package core

const (
	leftEvent  = 0
	leftMapper = 2
	leftName   = 3

	rightEvent  = 8
	rightMapper = 10
	rightName   = 11
)

type MapIndexRec struct {
	Event  string
	Mapper string
	Name   string
	Link   *string
}

func ParseMapIndexRecs(records [][]string) []MapIndexRec {
	out := make([]MapIndexRec, 0, len(records)*2)

	for _, row := range records {
		out = appendIfValid(out, row, leftEvent, leftMapper, leftName)
		out = appendIfValid(out, row, rightEvent, rightMapper, rightName)
	}

	return out
}

func appendIfValid(
	out []MapIndexRec,
	row []string,
	eventIdx, mapperIdx, nameIdx int,
) []MapIndexRec {
	event := valueAt(row, eventIdx)
	mapper := valueAt(row, mapperIdx)
	name := valueAt(row, nameIdx)
	if !isData(event, mapper, name) {
		return out
	}

	return append(out, MapIndexRec{
		Event:  event,
		Mapper: mapper,
		Name:   name,
	})
}

func valueAt(row []string, idx int) string {
	if idx >= 0 && idx < len(row) {
		return row[idx]
	}
	return ""
}

func isData(event string, mapper string, name string) bool {
	return event != "" &&
		mapper != "" &&
		name != "" &&
		!(event == "COTD #" && mapper == "Mapper" && name == "Map Name")
}
