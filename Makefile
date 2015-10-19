

build:
	go build ./...
	go vet ./...
	golint ./...
	find . -name '*.go' | xargs gofmt -w -s
	find . -name '*.go' | xargs goimports -w

golang-corporate-oss.md: Makefile
	go run ./github-search/main.go att airbnb aws bitly cloudflare coreos datadog docker ebay elastic etsy facebookgo fastly gilt \
		github google hashicorp heroku influxdb microsoft netflix samsung sendgrid sony soundcloud spotify \
		square stripe uber vimeo yahoo yelp > golang-corporate-oss.md
corp: golang-corporate-oss.md

golang-people-oss.md: Makefile
	go run ./github-search/main.go miekg ryanuber araddon mitchellh fatih tj BurntSushi philhofer tinylib alecthomas > golang-people-oss.md

people: golang-people-oss.md

clean:
	rm -f *~
	rm -f people.md corp.md



