package differ

type Delta struct {
	StartIndex      int
	EndIndex        int
	Deleted         bool
	UpdatedLiterals []byte
}

type PrettyDelta struct {
	startIndex      int
	endIndex        int
	deleted         bool
	updatedLiterals string
}

// convert byte array to string so that it is easy to read
func PrettifyDelta(deltas map[int]Delta) map[int]PrettyDelta {
	prettyDelta := make(map[int]PrettyDelta)
	for key, value := range deltas {
		prettyDelta[key] = PrettyDelta{
			startIndex:      value.StartIndex,
			endIndex:        value.EndIndex,
			deleted:         value.Deleted,
			updatedLiterals: string(value.UpdatedLiterals),
		}
	}
	return prettyDelta
}
