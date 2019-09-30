workflow "Code quality" {
  on = "pull_request"
  resolves = [
    "go imports",
<<<<<<< HEAD
    //"go lint",
    //"go vet",
=======
<<<<<<< HEAD
=======
    "go lint",
    "go vet",
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
    "staticcheck",
    "errcheck",
    "go sec",
    "shadow",
  ]
}

<<<<<<< HEAD
action "go imports" {
  uses = "grandcolline/golang-github-actions/imports@v0.2.0"
=======
<<<<<<< HEAD
action "filter: PR ready" {
  uses = "actions/bin/filter@master"
  args = "action 'opened|synchronize'"
}

action "go imports" {
  uses = "grandcolline/golang-github-actions/imports@v0.2.0"
  needs = "filter: PR ready"
=======
action "go imports" {
  uses = "grandcolline/golang-github-actions/imports@v0.2.0"
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
}

action "go lint" {
  uses = "grandcolline/golang-github-actions/lint@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
}

action "go vet" {
  uses = "grandcolline/golang-github-actions/vet@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
}

action "shadow" {
  uses = "grandcolline/golang-github-actions/shadow@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
}

action "staticcheck" {
  uses = "grandcolline/golang-github-actions/staticcheck@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
}

action "errcheck" {
  uses = "grandcolline/golang-github-actions/errcheck@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
  env = {
    IGNORE_DEFER = "true"
  }
}

action "go sec" {
  uses = "grandcolline/golang-github-actions/sec@v0.2.0"
<<<<<<< HEAD
=======
<<<<<<< HEAD
  needs = "filter: PR ready"
=======
>>>>>>> develop
>>>>>>> d6d99a1c13179649e63c4a4bc761510a38681bc9
  secrets = ["GITHUB_TOKEN"]
  env = {
    FLAGS = "-exclude=G104"
  }
}
