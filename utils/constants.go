package utils

const (
	CommandGet = "get"
	CommandSet = "set"
)

var validCommands = []string{CommandGet, CommandSet}

var cmdArgs = map[string]int{
	CommandSet: 2,
	CommandGet: 1,
}

var cmdArgsErrors = map[string]string{
	CommandGet: "GET [val]",
	CommandSet: "SET [val] [val]",
}
