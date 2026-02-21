package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for style guide compliance",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var sensitiveKeywords []string

var supportedLoggers = map[string]map[string]bool{
	"log/slog": {
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"DebugContext": true, "InfoContext": true, "WarnContext": true, "ErrorContext": true,
	},
	"go.uber.org/zap": {
		"Debug": true, "Info": true, "Warn": true, "Error": true, "DPanic": true, "Panic": true, "Fatal": true,
		"Debugf": true, "Infof": true, "Warnf": true, "Errorf": true, "DPanicf": true, "Panicf": true, "Fatalf": true,
		"Debugw": true, "Infow": true, "Warnw": true, "Errorw": true, "DPanicw": true, "Panicw": true, "Fatalw": true,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if sensitiveKeywords == nil {
		keywords, err := LoadSensitiveKeywords("../../testdata/.loglinter.yml")
		if err != nil {
			keywords, err = LoadSensitiveKeywords(".loglinter.yml")
			if err != nil {
				sensitiveKeywords = []string{}
			} else {
				sensitiveKeywords = keywords
			}
		} else {
			sensitiveKeywords = keywords
		}
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		msgArg := getLogMessageArg(pass, call)
		if msgArg == nil {
			return
		}

		lit, ok := msgArg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}

		msg := strings.Trim(lit.Value, "\"")
		if msg == "" {
			return
		}

		checkLowerCase(pass, lit, msg)

		checkEnglish(pass, lit, msg)

		checkSpecialChars(pass, lit, msg)

		checkSensitiveData(pass, lit, msg)
	})

	return nil, nil
}

func getLogMessageArg(pass *analysis.Pass, call *ast.CallExpr) ast.Expr {
	obj := typeutil.Callee(pass.TypesInfo, call)
	fn, ok := obj.(*types.Func)
	if !ok || fn.Pkg() == nil {
		return nil
	}

	if methods, ok := supportedLoggers[fn.Pkg().Path()]; ok {
		if methods[fn.Name()] {
			if len(call.Args) > 0 {
				return call.Args[0]
			}
		}
	}
	return nil
}

func checkLowerCase(pass *analysis.Pass, node ast.Node, msg string) {
	firstRune, _ := utf8.DecodeRuneInString(msg)
	if firstRune != utf8.RuneError && unicode.IsLetter(firstRune) && unicode.IsUpper(firstRune) {
		pass.Reportf(node.Pos(), "log message should start with a lowercase letter: %q", msg)
	}
}

func checkEnglish(pass *analysis.Pass, node ast.Node, msg string) {
	for _, r := range msg {
		if r > unicode.MaxASCII && unicode.IsLetter(r) {
			pass.Reportf(node.Pos(), "log message should contain only English words: %q", msg)
			return
		}
	}
}

func checkSpecialChars(pass *analysis.Pass, node ast.Node, msg string) {
	if strings.HasSuffix(msg, ".") || strings.HasSuffix(msg, "!") || strings.HasSuffix(msg, "?") {
		pass.Reportf(node.Pos(), "log message should not end with punctuation: %q", msg)
	}

	for _, r := range msg {
		if (unicode.Is(unicode.So, r) || unicode.Is(unicode.Sk, r)) && r > unicode.MaxASCII {
			pass.Reportf(node.Pos(), "log message should not contain emojis or special symbols: %q", msg)
			return
		}
	}
}

func checkSensitiveData(pass *analysis.Pass, node ast.Node, msg string) {
	lowerMsg := strings.ToLower(msg)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lowerMsg, kw) {
			pass.Reportf(node.Pos(), "log message contains potential sensitive data (%s): %q", kw, msg)
			return
		}
	}
}
