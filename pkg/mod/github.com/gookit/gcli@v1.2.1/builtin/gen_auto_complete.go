package builtin

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli"
	"github.com/gookit/gcli/interact"
	"github.com/gookit/gcli/utils"
	"github.com/gookit/goutil/cliutil"
)

// current supported shell
const (
	ZshShell  = "zsh"
	BashShell = "bash"
)

//
var genOpts = &struct {
	shell   string
	binName string
	output  string
}{}

var shellTpls = map[string]string{
	"zsh":  zshCompleteScriptTpl,
	"bash": bashCompleteScriptTpl,
}

// GenAutoCompleteScript create command
func GenAutoCompleteScript() *gcli.Command {
	c := gcli.Command{
		Func:    doGen,
		Name:    "gen:ac",
		Aliases: []string{"genac", "gen-ac"},
		// des
		UseFor: "generate auto complete scripts for current application",
	}

	shell := cliutil.CurrentShell(true)
	if shell == "" {
		shell = "bash"
	}

	c.StrOpt(
		&genOpts.shell,
		"shell",
		"s",
		shell,
		"the shell env name for want generated, allow: zsh,bash",
	).StrOpt(
		&genOpts.binName,
		"bin-name",
		"b",
		"",
		"your packaged application bin file name.",
	).StrOpt(
		&genOpts.output,
		"output",
		"o",
		"auto-completion.{shell}",
		"output shell auto completion script file name.",
	)

	return &c
}

func doGen(c *gcli.Command, _ []string) int {
	if len(genOpts.binName) == 0 {
		genOpts.binName = c.BinName()
	}

	genOpts.binName = strings.TrimSuffix(strings.Trim(genOpts.binName, "./"), ".exe")

	if len(genOpts.output) == 0 {
		genOpts.output = genOpts.binName + "." + genOpts.shell
	} else {
		genOpts.output = strings.Replace(genOpts.output, "{shell}", genOpts.shell, 1)

		// check suffix
		if !strings.Contains(genOpts.output, ".") {
			genOpts.output += "." + genOpts.shell
		}
	}

	color.Info.Tips("\n  %+v\n", genOpts)
	data := map[string]interface{}{
		"Shell":    genOpts.shell,
		"BinName":  genOpts.binName,
		"FileName": genOpts.output,
	}

	if genOpts.shell == BashShell {
		data = buildForBashShell(data)
	} else if genOpts.shell == ZshShell {
		data = buildForZshShell(data)
	} else {
		color.Error.Tips("--shell option only allow: zsh,bash")
		return gcli.ERR
	}

	str := utils.RenderTemplate(shellTpls[genOpts.shell], &data)

	color.Info.Println("Now, will write content to file ", genOpts.output)
	color.Normal.Print("Continue?")

	if !interact.AnswerIsYes(true) {
		color.Info.Print("\nBye :)\n")
		return 0
	}

	// ?????????????????????????????????????????????????????????
	err := ioutil.WriteFile(genOpts.output, []byte(str), 0664)
	if err != nil {
		return c.Errorf("Write file error: %s", err.Error())
	}

	color.Success.Println("\nOK, auto-complete file generate successful")
	return 0
}

var bashCompleteScriptTpl = `#!/usr/bin/env {{.Shell}}

# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# usage: source {{.FileName}}
# run 'complete' to see registered complete function.


_complete_for_{{.BinName}} () {
    local cur prev
    _get_comp_words_by_ref -n = cur prev

    COMPREPLY=()
    commands="{{join .CmdNames " "}} help"

    case "$prev" in{{range $k,$v := .NameOpts}}
        {{$k}})
            COMPREPLY=($(compgen -W "{{$v}}" -- "$cur"))
            return 0
            ;;{{end}}
        help)
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            return 0
            ;;
    esac

    COMPREPLY=($(compgen -W "$commands" -- "$cur"))

} &&
# complete -F {auto_complete_func} {bin_filename}
# complete -F _complete_for_{{.BinName}} -A file {{.BinName}} {{.BinName}}.exe
complete -F _complete_for_{{.BinName}} {{.BinName}} {{.BinName}}.exe
`

func buildForBashShell(data map[string]interface{}) map[string]interface{} {
	var cNames []string

	// {cmd name: opts}
	nameOpts := make(map[string]string)

	for n, c := range gcli.AllCommands() {
		// skip self
		if n == "genac" || n == "gen-ac" {
			continue
		}

		ops := c.OptNames()
		if len(ops) == 0 {
			continue
		}

		ns := c.Aliases
		key := n

		if len(ns) > 0 {
			ns = append(ns, n)
			key = strings.Join(ns, "|")
			cNames = append(cNames, ns...)
		} else {
			cNames = append(cNames, n)
		}

		var opList []string
		for op, st := range ops {
			if st != "" {
				opList = append(opList, "-"+st)
			}

			pfx := "--"
			if len(op) == 1 {
				pfx = "-"
			}

			opList = append(opList, pfx+op)
		}

		nameOpts[key] = strings.Join(opList, " ")
	}

	data["CmdNames"] = cNames
	data["NameOpts"] = nameOpts

	return data
}

var zshCompleteScriptTpl = `#compdef {{.BinName}}
# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# usage: source {{.FileName}}

_complete_for_{{.BinName}} () {
    typeset -a commands
    commands+=({{range $k,$v := .NameDes}}
        '{{$k}}[{{$v}}]'{{end}}
        'help[Display help information]'
    )

    if (( CURRENT == 2 )); then
        # explain commands
        _values 'cliapp commands' ${commands[@]}
        return
    fi

    case ${words[2]} in{{range $k,$vs := .NameOpts}}
    {{$k}})
        _values 'command options' \{{range $vs}}
            {{.}}{{end}}
        ;;{{end}}
    help)
        _values "${commands[@]}"
        ;;
    *)
        # use files by default
        _files
        ;;
    esac
}

compdef _complete_for_{{.BinName}} {{.BinName}}
compdef _complete_for_{{.BinName}} {{.BinName}}.exe
`

func buildForZshShell(data map[string]interface{}) map[string]interface{} {
	type opInfos []string

	// {cmd name: cmd des}. in zsh eg: 'build[compile packages and dependencies]'
	nameDes := make(map[string]string)
	// {cmd name: {opt: opt des}}.
	// in zsh eg:
	// '-x[description]:message:action'
	// {-h,--help}'[Show usage message]' // multi name
	nameOpts := make(map[string]opInfos)

	for n, c := range gcli.AllCommands() {
		// skip self
		if n == "genac" || n == "gen-ac" {
			continue
		}
		nameDes[c.Name] = fmtDes(c.UseFor) + "(alias " + c.AliasesString() + ")"

		ops := c.OptNames()
		oplen := len(ops)
		if oplen == 0 {
			continue
		}

		ns := c.Aliases
		key := n

		if len(ns) > 0 {
			ns = append(ns, n)
			key = strings.Join(ns, "|")
		}

		sfx := " \\"
		var i int
		var opis []string
		for op, st := range ops {
			i++
			pfx := "--"
			opDes := fmtDes(c.Flags.Lookup(op).Usage)

			if len(op) == 1 {
				pfx = "-"
			}

			opKey := pfx + op
			desTpl := "'%s[%s]'%s"

			if st != "" {
				desTpl = "%s'[%s]'%s"
				opKey = fmt.Sprintf("{-%s,%s}", st, pfx+op)
			}

			// latest item
			if oplen == i {
				sfx = ""
			}

			opis = append(opis, fmt.Sprintf(desTpl, opKey, opDes, sfx))
		}

		nameOpts[key] = opis
	}

	data["NameDes"] = nameDes
	data["NameOpts"] = nameOpts

	return data
}

func fmtDes(str string) string {
	str = color.ClearTag(str)
	return strings.NewReplacer("`", "", "[", "", "]", "").Replace(str)
}
