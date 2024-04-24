package filter

const (
	DateStr = "string"
	DateInt = "integer"

	OperatorEq         = "="
	OperatorNotEq      = "!="
	OperatorLowerThen  = "<="
	OperatorHigherThen = ">="
	OperatorBetween    = " between "

	ParamModel  = "model"
	ParamRegNum = "reg_num"
	ParamMark   = "mark"
	ParamYear   = "year"

	GroupStart   = "start"
	GroupElement = "el"
	GroupEnd     = "end"
	GroupNil     = "null"
)

type Option struct {
	Limit  int
	Offset int
	Fileds []*Filed
}

type Filed struct {
	Param    string
	Operator string
	Value    string
	DataType string
	Group    string
}

type Options interface {
	GetFileds() []*Filed
	AddFileds(param, operator, value, dataType, slim string)
}

func NewOption() Option {
	return Option{}
}

func (o *Option) GetFileds() []*Filed {
	return o.Fileds
}

func (o *Option) AddFileds(param, operator, value, dataType, group string) {
	o.Fileds = append(o.Fileds, &Filed{
		Param:    param,
		Operator: operator,
		Value:    value,
		DataType: dataType,
		Group:    group,
	})
}
