.PHONY:
dev:
	go run main.go check --host planetary-quantum.com

build:
	goreleaser \
		--snapshot \
		--skip-publish \
		--rm-dist
