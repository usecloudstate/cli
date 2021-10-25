run:
	go run main.go

build:
	./build.sh

setup_semantic_release:
	npm i @semantic-release/commit-analyzer @semantic-release/release-notes-generator @semantic-release/github @semantic-release/changelog @semantic-release/exec @semantic-release/git

update_install_script_version:
	sed -e "s/CLI_VERSION/$(version)/g" install.sh