

build:
	go build ./...
	go vet ./...
	golint ./...
	gofmt -w -s ./...

golang-corporate-oss.md
	go run ./github-search/main.go att airbnb aws bitly cloudflare coreos datadog docker ebay elastic etsy facebookgo fastly gilt \
		github google hashicorp heroku influxdb microsoft netflix samsung sendgrid sony soundcloud spotify \
		square stripe uber vimeo yahoo yelp > golang-corporate-oss.md
corp: corp.md

people.md:
	go run ./github-search/main.go miekg ryanuber araddon mitchellh fatih tj > people.md

people: people.md

clean:
	rm -f *~
	rm -f people.md corp.md



