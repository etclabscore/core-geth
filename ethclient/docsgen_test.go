package ethclient

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"text/template"

	meta_schema "github.com/open-rpc/meta-schema"
)

// The test in this file is only for use as a development tool.

// TestRPCDiscover_BuildStatic puts the OpenRPC document in build/static/openrpc.json.
// This is intended to be run as a documentation development tool (as opposed to an actual _test_).
// NOTE that Go maps don't guarantee order, so the diff between runs can be noisy.
func TestRPCDiscover_BuildStatic(t *testing.T) {
	if os.Getenv("COREGETH_GEN_OPENRPC_DOCS") == "" {
		return
	}
	err := os.MkdirAll("../build/static", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	backend, _ := newTestBackend(t)
	client := backend.Attach()
	defer backend.Close()
	defer client.Close()

	// Current workaround for https://github.com/open-rpc/meta-schema/issues/356.
	res, err := backend.InprocDiscovery_DEVELOPMENTONLY().Discover()
	if err != nil {
		t.Fatal(err)
	}

	// Should do it this way.
	// res := &meta_schema.OpenrpcDocument{}
	// err = client.Call(res, "rpc.discover")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("../build/static/openrpc.json", data, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if res.Methods == nil {
		return
	}

	tpl := template.New("openrpc_doc")

	trimNameSpecialChars := func(s string) string {
		remove := []string{".", "*"}
		for _, r := range remove {
			s = strings.ReplaceAll(s, r, "")
		}
		return s
	}

	contentDescriptorGenericName := func(cd *meta_schema.ContentDescriptorObject) (name string) {
		defer func() {
			name = fmt.Sprintf("<%s>", name) // make it generic-looking
		}()
		if cd.Description == nil {
			return string(*cd.Name)
		}
		if cd.Name == nil {
			return string(*cd.Description)
		}
		tr := trimNameSpecialChars(string(*cd.Description))
		if string(*cd.Name) == tr {
			return string(*cd.Description)
		}
		return string(*cd.Name)
	}

	tpl.Funcs(template.FuncMap{
		// tomap gives a plain map from a JSON representation of a given value.
		// This is useful because the meta_schema data types, being generated, and conforming to a pretty
		// complex data type in the first place, are not super fun to interact with directly.
		"tomap": func(any interface{}) map[string]interface{} {
			out := make(map[string]interface{})
			data, _ := json.Marshal(any)
			json.Unmarshal(data, &out)
			return out
		},
		// asjson returns indented JSON.
		"asjson": func(any interface{}, prefix, indent int) string {
			by, _ := json.MarshalIndent(any, strings.Repeat("    ", prefix), strings.Repeat("    ", indent))
			return string(by)
		},
		// bulletJSON handles transforming a JSON JSON schema into bullet points, which I think are more legible.
		"bulletJSON": printBullet,
		"sum":        func(a, b int) int { return a + b },
		// trimNameSpecialChars removes characters that the app-specific content descriptor naming
		// method will also remove, eg '*hexutil.Uint64' -> 'hexutilUint64'.
		// These "special" characters were removed because of concerns about by-name arguments
		// and the use of titles for keys.
		"trimNameSpecialChars": trimNameSpecialChars,
		// contentDescriptorTitle returns the name or description, in that order.
		"contentDescriptorGenericName": contentDescriptorGenericName,
		// isSubscribeMethod checks whether the method is a `*_subscribe` method
		"isSubscribeMethod": func(m *meta_schema.MethodObject) bool {
			return *m.Result.ContentDescriptorObject.Name == "subscriptionID"
		},
		// isSubscriptionableMethod checks whether the method returns a subscription
		"isSubscriptionableMethod": func(m *meta_schema.MethodObject) bool {
			return *m.Result.ContentDescriptorObject.Description == "*rpc.Subscription"
		},
		// isFilterMethod checks whether the method is being used as a filter.
		"isFilterMethod": func(m *meta_schema.MethodObject) bool {
			// TODO: later, check for `*filters.PublicFilterAPI` only
			return *m.Result.ContentDescriptorObject.Name == "rpcID"
		},
		// methodFormatJSConsole is a pretty-printer that returns the JS console use example for a method.
		"methodFormatJSConsole": func(m *meta_schema.MethodObject) string {
			name := string(*m.Name)
			formattedName := strings.Replace(name, "_", ".", 1)
			getParamName := func(cd *meta_schema.ContentDescriptorObject) string {
				if cd.Name != nil {
					return string(*cd.Name)
				}
				return string(*cd.Description)
			}
			paramNames := func() (paramNames []string) {
				if m.Params == nil {
					return nil
				}
				for _, n := range *m.Params {
					if n.ContentDescriptorObject == nil {
						continue // Should never happen in our implementation; never uses refs.
					}
					paramNames = append(paramNames, getParamName(n.ContentDescriptorObject))
				}
				return
			}()
			return fmt.Sprintf("%s(%s);", formattedName, strings.Join(paramNames, ","))
		},
		// methodFormatCURL is a pretty printer that returns the 'curl' method invocation example string.
		"methodFormatCURL": func(m *meta_schema.MethodObject) string {
			paramNames := ""
			if m.Params != nil {
				out := []string{}
				for _, p := range *m.Params {
					out = append(out, contentDescriptorGenericName(p.ContentDescriptorObject))
				}
				paramNames = strings.Join(out, ", ")
			}

			return fmt.Sprintf(`curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "%s", "params": [%s]}'`, *m.Name, paramNames)
		},
		// methodFormatWS is a pretty printer that returns the websocket's method invocation example string.
		"methodFormatWS": func(m *meta_schema.MethodObject) string {
			paramNames := ""
			if m.Params != nil {
				out := []string{}
				for _, p := range *m.Params {
					out = append(out, contentDescriptorGenericName(p.ContentDescriptorObject))
				}
				paramNames = strings.Join(out, ", ")
			}

			return fmt.Sprintf(`wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "%s", "params": [%s]}'`, *m.Name, paramNames)
		},
		// methodFormatWSSubscribe is a pretty printer that returns the websocket's method invocation example string using the package `*_subscribe` method.
		"methodFormatWSSubscribe": func(m *meta_schema.MethodObject) string {
			methodName := string(*m.Name)
			outParams := []string{}

			methodParts := strings.Split(methodName, "_")
			if len(methodParts) == 2 {
				methodName = methodParts[0] + "_subscribe"
				outParams = append(outParams, fmt.Sprintf("\"%s\"", methodParts[1]))
			}

			paramNames := ""
			if m.Params != nil {
				for _, p := range *m.Params {
					outParams = append(outParams, contentDescriptorGenericName(p.ContentDescriptorObject))
				}
				paramNames = strings.Join(outParams, ", ")
			}

			return fmt.Sprintf(`wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "%s", "params": [%s]}'`, methodName, paramNames)
		},
	})

	tpl, err = tpl.Parse(`
{{ define "schemaTpl" }}
	` + "```" + `
{{ bulletJSON . 1 }}
	` + "```" + `

	<details class="cite"><summary>View Raw</summary>
	` + "```" + `
    {{ asjson . 1 1 }}
	` + "```" + `
	</details>
{{ end }}

{{ define "contentDescTpl" -}}
{{ $nameyDescription := trimNameSpecialChars .description }}
{{ if eq .name $nameyDescription }}
<code>{{ .description }}</code> {{ if .summary }}_{{ .summary }}_{{- end }}
{{- else -}}
{{ .name }} <code>{{ .description }}</code> {{ if .summary }}_{{ .summary }}_{{- end }}
{{- end }}

  + Required: {{ if .required }}✓ Yes{{ else }}No{{- end}}
{{ if .deprecated }}  + Deprecated: :warning: Yes{{- end}}
{{ if (or (gt (len .schema) 1) .schema.properties) }}
=== "Schema"

	` + "```" + ` Schema
	{{ bulletJSON .schema 1 }}
	` + "```" + `

=== "Raw"

	` + "```" + ` Raw
	{{ asjson .schema 1 1 }}
	` + "```" + `
{{ end }}
{{ end }}

{{ define "methodTpl" }}
{{ $methodmap := tomap . }}
### {{ .Name }}

{{ .Summary }}

#### Params ({{ .Params | len }})
{{ if gt (.Params | len) 0 }}
{{ if eq $methodmap.paramStructure "by-position" }}Parameters must be given _by position_.{{ else if eq $methodmap.paramStructure "by-name" }}Parameters must be given _by name_.{{ end }}
{{ range $index, $param := .Params }}
{{ $parammap := . | tomap }}
__{{ sum $index 1 }}:__ {{ template "contentDescTpl" $parammap }}
{{ end }}
{{ else }}
_None_
{{- end}}

#### Result

{{ if .Result -}}
{{ $result := .Result | tomap }}
{{- if ne $result.name "Null" }}
{{ template "contentDescTpl" $result }}
{{- else -}}
_None_
{{- end }}
{{- end }}

#### Client Method Invocation Examples

{{ if and (not (isSubscribeMethod .)) (and (not (isSubscriptionableMethod .)) (not (isFilterMethod .))) }}
=== "Shell HTTP"

	` + "```" + ` shell
	{{ methodFormatCURL . }}
	` + "```" + `
{{ end }}

{{ $shellWSExample := methodFormatWS . }}
{{ if isSubscriptionableMethod . }}
{{ $shellWSExample = methodFormatWSSubscribe . }}
{{ end }}

=== "Shell WebSocket"

	` + "```" + ` shell
	{{ $shellWSExample }}
	` + "```" + `

{{ if and (not (isSubscribeMethod .)) (and (not (isSubscriptionableMethod .)) (not (isFilterMethod .))) }}
=== "Javascript Console"

	` + "```" + ` js
	{{ methodFormatJSConsole . }}
	` + "```" + `
{{ end }}

{{ $docs := .ExternalDocs | tomap }}
<details><summary>Source code</summary>
<p>
{{ .Description }}
<a href="{{ $docs.url }}" target="_">View on GitHub →</a>
</p>
</details>

---
{{- end }}

| Entity | Version |
| --- | --- |
| Source | <code>{{ .Info.Version }}</code> |
| OpenRPC | <code>{{ .Openrpc }}</code> |

---

{{ range .Methods }}
{{ template "methodTpl" . }}
{{ end }}
`)
	if err != nil {
		t.Fatal(err)
	}

	moduleMethods := func() (grouped map[string][]meta_schema.MethodObject) {
		if res.Methods == nil {
			return
		}
		grouped = make(map[string][]meta_schema.MethodObject)
		for _, m := range *res.Methods {
			moduleName := strings.Split(string(*m.Name), "_")[0]
			group, ok := grouped[moduleName]
			if !ok {
				group = []meta_schema.MethodObject{}
			}
			group = append(group, m)
			grouped[moduleName] = group
		}
		return
	}()

	_ = os.MkdirAll("../docs/JSON-RPC-API/modules", os.ModePerm)

	for module, group := range moduleMethods {
		fname := fmt.Sprintf("../docs/JSON-RPC-API/modules/%s.md", module)
		fi, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		fi.Truncate(0)

		nDoc := &meta_schema.OpenrpcDocument{}
		*nDoc = *res
		nDoc.Methods = (*meta_schema.Methods)(&group) //nolint:gosec
		err = tpl.Execute(fi, nDoc)
		if err != nil {
			t.Fatal(err)
		}
		fi.Sync()
		fi.Close()
	}
}

// printBullet handles transforming a JSON JSON schema into bullet points, which I think are more legible.
func printBullet(any interface{}, depth int) (out string) {
	defer func() {
		out += "\n"
	}()
	switch typ := any.(type) {
	case map[string]interface{}:
		out += "\n"
		ordered := []string{}
		for k := range typ {
			ordered = append(ordered, k)
		}
		sort.Slice(ordered, func(i, j int) bool {
			return ordered[i] < ordered[j]
		})
		for _, k := range ordered {
			v := typ[k]
			out += fmt.Sprintf("%s- %s: %s", strings.Repeat("\t", depth), k, printBullet(v, depth+1))
		}
	case []interface{}:

		// Don't know why this isn't working. Doesn't fire.
		// if c, ok := any.([]string); ok {
		// 	return strings.Join(c, ",")
		// }

		stringSet := []string{}
		for _, vv := range typ {
			if s, ok := vv.(string); ok {
				stringSet = append(stringSet, s)
			}
		}
		if len(stringSet) == len(typ) {
			return strings.Join(stringSet, ", ")
		}

		out += "\n"
		for _, vv := range typ {
			out += printBullet(vv, depth+1)
		}
	default:
		return fmt.Sprintf("`%v`", any)
	}
	return
}
