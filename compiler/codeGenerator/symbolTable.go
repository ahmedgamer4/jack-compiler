package codegenerator

type FieldType string

const (
	Field  FieldType = "field"
	Static FieldType = "static"
	Arg    FieldType = "arg"
	Lcl    FieldType = "local"
)

type Table struct {
	Name  string
	Type  string
	Kind  FieldType
	Index int
}

type SymbolTable struct {
	ClassTable      Table
	SubroutineTable Table
	ClassIdx        int
	SubroutineIdx   int
}

func (s *SymbolTable) Reset() {
}

func (s *SymbolTable) Define(name, typ string, kind FieldType) {
}

func (s *SymbolTable) VarCount() int {
	return 0
}

func (s *SymbolTable) KindOf(name string) FieldType {
	return ""
}

func (s *SymbolTable) TypeOf(name string) string {
	return ""
}

func (s *SymbolTable) IndexOf(name string) string {
	return ""
}
