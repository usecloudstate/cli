#!/usr/bin/env bash

package_name="cloudstate-cli"

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    next_ver=$(npx semantic-release --dryRun | grep -oP 'Published release \K.*? ')

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X main.Version=$next_ver" -o bin/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done