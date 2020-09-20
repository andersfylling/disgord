#!/bin/bash

# github REST commands
wget https://raw.githubusercontent.com/whiteinge/ok.sh/master/ok.sh -O gith
chmod +x gith

# get last closed milestone (yes this is a race cond)
./gith list_milestones andersfylling/disgord state=closed sort=closed_at direction=desc > closed_milestones.txt
VERSION="$(awk 'NR==1 {print $3}' "closed_milestones.txt")"
if [[ ! ${VERSION:0:1} == "v" ]] ; then
  >&2 echo "ERROR: $(cat closed_milestones.txt)";
  exit 1
fi

# cleanup
rm gith
rm "closed_milestones.txt"

# setup git env
git config user.email "${GITHUB_EMAIL}"
git config user.name "disgord (bot)"
git remote set-url origin https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git
git checkout develop

# update version const
VERSION_FILE="constant/version.go"
echo "package constant" > "${VERSION_FILE}"
echo "const Version = \"${VERSION}\"" >> "${VERSION_FILE}"
go fmt ./...
go mod tidy

# generate CHANGELOG.md
#gem install github_changelog_generator
#github_changelog_generator --release-branch develop -u andersfylling -p disgord

# commit, tag and push
if [[ `git status --porcelain` ]]; then
  git add .
  #git commit -m "gen changelog & set version to ${VERSION}"
  git commit -m "set version to ${VERSION}"
fi
git tag "${VERSION}" -m "Disgord ${VERSION}"
git push origin develop --tags
