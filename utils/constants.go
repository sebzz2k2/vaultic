package utils

const (
	CommandGet    = "GET"
	CommandSet    = "SET"
	CommandDel    = "DEL"
	CommandExists = "EXISTS"
	CommandKeys   = "KEYS"

	FILENAME  = "vaultic"
	DELIMITER = ":"
	NEWLINE   = "\n"
	SUCCESS   = "OK"
)

var CmdArgs = map[string]int{
	CommandSet:    2,
	CommandGet:    1,
	CommandDel:    1,
	CommandExists: 1,
	CommandKeys:   0,
}

var CmdArgsErrors = map[string]string{
	CommandGet: "GET [val]",
	CommandSet: "SET [val] [val]",
}
