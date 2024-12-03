package mapp

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

var ruleReg = regexp.MustCompile(`^@(qual|enum|ignore) `)

type Mapper struct {
	spec    *ast.Field
	imports []Import
}

func (m Mapper) Name() string {
	return m.spec.Names[0].Name
}

func (m Mapper) Comments() []Comment {
	comments := make([]Comment, 0, len(m.spec.Doc.List))
	for _, v := range m.spec.Doc.List {
		comments = append(comments, Comment{spec: v})
	}

	return comments
}

func (m Mapper) Rules() []Rule {
	rules := []Rule{}
	for _, c := range m.Comments() {
		val := c.Value()
		if !strings.HasPrefix(val, "@") {
			continue
		}

		if ruleReg.MatchString(val) {
			buildRuleArg := func(arg string) (RuleArg, string) {
				kvPair := strings.Split(arg, "=")
				if len(kvPair) != 2 {
					panic(fmt.Sprintf("invalid argument format: %s", arg))
				}

				switch kvPair[0] {
				case "-s":
					return RuleArgSource, kvPair[1]
				case "-t":
					return RuleArgTarget, kvPair[1]
				case "-mn":
					return RuleArgMname, kvPair[1]
				case "-mp":
					return RuleArgMpath, kvPair[1]
				default:
					panic(fmt.Sprintf("unknown arg key: %s", kvPair[0]))
				}
			}

			args := strings.Split(val, " ")
			ruleArgs := make(map[RuleArg]string, len(args)-1)
			for _, a := range args[1:] {
				k, v := buildRuleArg(a)
				ruleArgs[k] = v
			}
			rules = append(rules, Rule{
				spec: val,
				args: ruleArgs,
			})
		} else {
			panic(fmt.Sprintf("unsupported rule: %s", val))
		}
	}

	return rules
}

func (m Mapper) Params() []Param {
	fnT, ok := m.spec.Type.(*ast.FuncType)
	if !ok {
		panic("is not a ast.FuncType")
	}

	params := make([]Param, 0, len(fnT.Params.List))
	for _, p := range fnT.Params.List {
		params = append(params, Param{spec: p, imports: m.imports})
	}

	return params
}

func (m Mapper) Results() []Result {
	fnT, ok := m.spec.Type.(*ast.FuncType)
	if !ok {
		panic("is not a ast.FuncType")
	}

	params := make([]Result, 0, len(fnT.Results.List))
	for _, p := range fnT.Results.List {
		params = append(params, Result{spec: p, imports: m.imports})
	}

	return params
}

func (m Mapper) Source() Source {
	return Source{
		spec:  m.Params()[0].spec,
		Param: m.Params()[0],
	}
}

func (m Mapper) Target() Target {
	return Target{
		spec: m.Results()[0].spec,
		r:    m.Results()[0],
	}
}

func (m Mapper) SourceFieldByTarget(targetFullName string) (Field, bool) {
	source := m.Source()

	firstElemEndPos := strings.IndexAny(targetFullName, ".")
	sourceFullName := source.Name() + targetFullName[firstElemEndPos:]
	for _, r := range m.Rules() {
		if r.Type() != RuleTypeQual {
			continue
		}
		tname, ok := r.Arg(RuleArgTarget)
		if !ok {
			continue
		}
		if tname != targetFullName {
			continue
		}
		
		if r.Type() == RuleTypeIgnore {
			panic(fmt.Sprintf("target field: %s must be ignored", targetFullName))
		}

		sname, ok := r.Arg(RuleArgSource)
		if !ok {
			continue
		}
		if strings.Contains(sname, ".") {
			sourceFullName = sname
			break
		}
		lastElemStartPos := strings.LastIndexAny(targetFullName, ".")
		pref := targetFullName[:lastElemStartPos]
		sourceFullName = pref + "." + sname
		break
	}

	return source.FieldByFullName(sourceFullName)
}

func (m Mapper) RulesByFieldFullName(fullName string) []Rule {
	rules := make([]Rule, 0)
	for _, r := range m.Rules() {
		_, ok := r.Arg(RuleArgTarget)
		if !ok {
			continue
		}
		rules = append(rules, r)
	}

	return rules
}
