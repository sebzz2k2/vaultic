package utils

const (
	CommandGet = "get"
	CommandSet = "set"
	FILENAME   = "vaultic"
	DELIMITER  = ":"
	NEWLINE    = "\n"
)

var CmdArgs = map[string]int{
	CommandSet: 2,
	CommandGet: 1,
}

var CmdArgsErrors = map[string]string{
	CommandGet: "GET [val]",
	CommandSet: "SET [val] [val]",
}
