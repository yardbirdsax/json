package json

type jsonSyntax int

const (
  OBJSTART jsonSyntax = iota // {
  OBJEND // }
  QUOTE // "
  COLON // :
  COMMA // ,
)

var jsonSyntaxes = []rune{
  OBJSTART: '{',
  OBJEND: '}',
  QUOTE: '"',
  COLON: ':',
  COMMA: ',',
}

func (t jsonSyntax) String() string{
  return string(jsonSyntaxes[t])
}

func (t jsonSyntax) Rune() rune {
  return jsonSyntaxes[t]
}