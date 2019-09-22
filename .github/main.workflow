workflow "Code quality" {
  on = "pull_request"
  resolves = [
    "go imports",
    //"go lint",
    //"go vet",
    "staticcheck",
    "errcheck",
    "go sec",
    "shadow",
  ]
}

action "go imports" {
  uses = "grandcolline/golang-github-actions/imports@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

action "go lint" {
  uses = "grandcolline/golang-github-actions/lint@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

action "go vet" {
  uses = "grandcolline/golang-github-actions/vet@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

action "shadow" {
  uses = "grandcolline/golang-github-actions/shadow@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

action "staticcheck" {
  uses = "grandcolline/golang-github-actions/staticcheck@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
}

action "errcheck" {
  uses = "grandcolline/golang-github-actions/errcheck@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
  env = {
    IGNORE_DEFER = "true"
  }
}

action "go sec" {
  uses = "grandcolline/golang-github-actions/sec@v0.2.0"
  secrets = ["GITHUB_TOKEN"]
  env = {
    FLAGS = "-exclude=G104"
  }
}
