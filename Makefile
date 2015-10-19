

build:
	go build ./...
	go vet ./...
	golint ./...
	find . -name '*.go' | xargs gofmt -w -s
	find . -name '*.go' | xargs goimports -w

golang-corporate-oss.md: Makefile
	go run ./github-search/main.go att airbnb aws bitly cloudflare coreos datadog docker ebay elastic etsy facebookgo fastly gilt \
		github golang google hashicorp heroku influxdb microsoft netflix pivotal-golang samsung sendgrid sony soundcloud spotify \
		square stripe uber vimeo yahoo yelp > golang-corporate-oss.md
corp: golang-corporate-oss.md

golang-people-oss.md: Makefile
	go run ./github-search/main.go \
		alecthomas \
		araddon	\
		BurntSushi \
		davecheney \
		fatih \
		mattn \
		miekg \
		mitchellh \
		tinylib \
		tj \
		philhofer \
		robertkrimen \
		robpike \
		ryanuber \
	> golang-people-oss.md

golang-contributors.md: Makefile
	go run ./github-search/main.go \
		0intro 4ad aclements alexcesaro ality atom-symbol bradfitz c4milo campoy crawshaw \
		dsymonds dvyukov jpoirier marete minux mpvl nicolasgarnier rakyll rui314 wathiede \
	> golang-contributors.md

people: golang-people-oss.md

clean:
	rm -f *~
	rm -f people.md corp.md



