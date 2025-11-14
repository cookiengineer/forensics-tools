package blocklists

var Blocklists map[string]map[string]bool

func init() {
	Blocklists = make(map[string]map[string]bool)
}
