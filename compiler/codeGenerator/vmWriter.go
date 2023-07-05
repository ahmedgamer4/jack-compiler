package codegenerator

/**
* Simple module that writes VM commands to the output .vm file
* */

type Segment string
type Command string

type VMWriter struct {
}

func (v VMWriter) writePush(seg Segment) {
}

func (v VMWriter) writePop(seg Segment) {
}

func (v VMWriter) writeArithmetic(command Command) {
}

func (v VMWriter) writeLabel(label string) {
}

func (v VMWriter) writeGoto(label string) {
}

func (v VMWriter) writeIf(label string) {
}

func (v VMWriter) writeCall(name string, nArgs int) {
}

func (v VMWriter) writeFunction(name string, nArgs int) {
}

func (v VMWriter) writeReturn() {
}

/**
* Close the output file / stream
* */
func (v VMWriter) Close() {
}
