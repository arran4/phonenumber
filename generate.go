package phonenumber

//go:generate sh -c "(command -v gosubc >/dev/null 2>&1 && gosubc generate || go run github.com/arran4/go-subcommand/cmd/gosubc@v0.0.17 generate) && sed -i 's/RunCLI(/phonenumber.RunCLI(/' cmd/drawphonecli/root.go && (test -f cmd/drawphonecli/templates/usage.txt || echo 'Usage: {{.Name}} [options]' > cmd/drawphonecli/templates/usage.txt)"
