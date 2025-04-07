package internal

type SupportedVarType string

const (
	StringType SupportedVarType = "string"
	IntType    SupportedVarType = "int"
	FloatType  SupportedVarType = "float"
	BoolType   SupportedVarType = "bool"
	AnyType    SupportedVarType = "any"
)

type VarEntry struct {
	Name         string
	Required     bool
	Type         SupportedVarType
	DefaultValue any
	Pattern      string
}

type Config struct {
	Vars []VarEntry
}

/*
This function reads .env.yaml file and returns config.
*/
func GetConfig() (Config, error) {
	// TODO
	return Config{}, nil
}
