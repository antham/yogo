compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"
fmt:
	gofmt -l -s -w .
version:
	git stash -u
	sed -i "s/[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+/$(v)/g" main.go
	git add -A
	git commit -m "feat(version) : "$(v)
	git tag v$(v) master
