package codegenerator

/**
* A module to handle variables
**/

type FieldType string

const (
	Field  FieldType = "field"
	Static FieldType = "static"
	Arg    FieldType = "argument"
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
	s.ArgIdx = -1
	s.LclIdx = -1
	s.SubroutineSymbolTable = map[string]Var{}
}

func (s *SymbolTable) ResetClassTable() {
	s.StaticIdx = -1
	s.FieldIdx = -1
	s.ClassSymbolTable = map[string]Var{}
}

func (s *SymbolTable) Define(name, typ string, kind FieldType) {
	switch kind {

	case Arg, Lcl:
		if _, ok := s.SubroutineSymbolTable[name]; !ok {
			s.SubroutineSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.getNextIdx(kind)}
		} else {
			println("var already exists", name)
		}
	case Field, Static:
		if _, ok := s.ClassSymbolTable[name]; !ok {
			s.ClassSymbolTable[name] = Var{Type: typ, Kind: kind, Index: s.getNextIdx(kind)}
		} else {
			println("var already exists", name)
		}
	default:
		break
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
		return s.StaticIdx + 1
	case Field:
		return s.FieldIdx + 1
	case Lcl:
		return s.LclIdx + 1
	case Arg:
		return s.ArgIdx + 1
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
