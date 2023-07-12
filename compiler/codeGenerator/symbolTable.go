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

func NewSTable() *SymbolTable {
	s := &SymbolTable{}
	s.ClassSymbolTable = map[string]Var{}
	s.SubroutineSymbolTable = map[string]Var{}
	s.FieldIdx = -1
	s.LclIdx = -1
	s.ArgIdx = -1
	s.StaticIdx = -1
	return s
}

func (s *SymbolTable) ResetSubroutineTable() {
	s.SubroutineSymbolTable = map[string]Var{}
}

func (s *SymbolTable) Define(name, typ string, kind FieldType) {
	switch kind {

	case Arg, Lcl:
		if _, ok := s.SubroutineSymbolTable[name]; !ok {
			s.SubroutineSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.getNextIdx(kind)}
		}
	case Field, Static:
		if _, ok := s.ClassSymbolTable[name]; !ok {
			s.ClassSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.getNextIdx(kind)}
		}
	default:
		println("var already exists", name)
	}
}

func (s *SymbolTable) getNextIdx(kind FieldType) int {
	var idxK *int

	switch kind {
	case Static:
		idxK = &s.StaticIdx
	case Field:
		idxK = &s.FieldIdx
	case Lcl:
		idxK = &s.LclIdx
	case Arg:
		idxK = &s.ArgIdx
	default:
		println("Kind does not exist", kind)
		return -1
	}

	*idxK++
	return *idxK
}

/**
* Return vars count for the passed kind
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
		println(name)
		// panic("empty")
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
