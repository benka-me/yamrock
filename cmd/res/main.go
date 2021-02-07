package res

// config has:
// - name: "config"
// - type: "object"
// - interface with one implementation by children:  {NAME}() {TYPE}
// - struct with one var by children (Obj{NAME})
// - local variable which contain the value (var{NAME}) of children
// For each children
// 		- data inside object (Var{DATA_NAME}) type/ Obj{NAME}
//      - one method as func {NAME}() {TYPE} {}

// First Level Object
func GetConfig() Config {
	return &varConfig
}

type Config interface {
	Name() string
	Size() int64
	Deps() []string
	Sub() Sub
}
type ObjConfig struct {
	VarName string   `yaml:"name"`
	VarSize int64    `yaml:"size"`
	VarSub  *ObjSub  `yaml:"sub"`
	VarDeps []string `yaml:"deps"`
}

var varConfig = ObjConfig{
	VarName: "hello-world",
	VarSize: 400,
	VarSub:  varSub,
	VarDeps: []string{
		"bonjour",
		"salam",
		"hola",
		"konichua",
		"bom dia",
	},
}

func (m *ObjConfig) Name() string {
	return m.VarName
}
func (m *ObjConfig) Size() int64 {
	return m.VarSize
}
func (m *ObjConfig) Sub() Sub {
	return m.VarSub
}
func (m *ObjConfig) Deps() []string {
	return m.VarDeps
}

// Second Level Obj
type Sub interface {
	Name() string
	Size() int64
}
type ObjSub struct {
	VarName string `yaml:"name"`
	VarSize int64  `yaml:"size"`
}

var varSub = &ObjSub{
	VarName: "byebye",
	VarSize: 800,
}

func (m *ObjSub) Name() string {
	return m.VarName
}
func (m *ObjSub) Size() int64 {
	return m.VarSize
}
