package codegenerator

/**
* A module to handle variables
**/

type FieldType string

const (
	Field  FieldType = "field"
	Static FieldType = "static"
	Arg    FieldType = "arg"
	Lcl    FieldType = "local"
)

type Var struct {
	Type  string
	Kind  FieldType
	Index int
}

type SymbolTable struct {
	ClassSymbolTable      map[string]Var
	SubroutineSymbolTable map[string]Var
	StaticIdx             int
	FieldIdx              int
	ArgIdx                int
	LclIdx                int
}

func (s *SymbolTable) ResetSubroutineTable() {
	s.SubroutineSymbolTable = map[string]Var{}
}

func (s *SymbolTable) Define(name, typ string, kind FieldType) {
	s.ClassSymbolTable = map[string]Var{}
	s.SubroutineSymbolTable = map[string]Var{}

	switch kind {

	case Arg, Lcl:
		if _, ok := s.SubroutineSymbolTable[name]; !ok {
			s.SubroutineSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.VarCount(kind) + 1}
		}
	case Field, Static:
		if _, ok := s.ClassSymbolTable[name]; !ok {
			s.ClassSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.VarCount(kind) + 1}
		}
	default:
		println("var already exists", name)
		// }
	}
}

/**
* Return vars count for every kind
* */
func (s *SymbolTable) VarCount(kind FieldType) int {
	switch kind {
	case Static:
		return s.StaticIdx
	case Field:
		return s.FieldIdx
	case Lcl:
		return s.LclIdx
	case Arg:
		return s.ArgIdx
	default:
		println("Kind does not exist", kind)
		return -1
	}
}

func (s *SymbolTable) KindOf(name string) FieldType {
	if v, ok := s.SubroutineSymbolTable[name]; ok {
		return v.Kind
	} else if v, ok := s.ClassSymbolTable[name]; ok {
		return v.Kind
	} else {
		return ""
	}
}

func (s *SymbolTable) TypeOf(name string) string {
	if v, ok := s.SubroutineSymbolTable[name]; ok {
		return v.Type
	} else if v, ok := s.ClassSymbolTable[name]; ok {
		return v.Type
	} else {
		return ""
	}
}

func (s *SymbolTable) IndexOf(name string) int {
	if v, ok := s.SubroutineSymbolTable[name]; ok {
		return v.Index
	} else if v, ok := s.ClassSymbolTable[name]; ok {
		return v.Index
	} else {
		return -1
	}
}
