run:
	go run main.go

build:
	./build.sh

setup_semantic_release:
	npm i @semantic-release/commit-analyzer @semantic-release/release-notes-generator @semantic-release/github @semantic-release/changelog