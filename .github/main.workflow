workflow "Go lint" {
  on = "pull_request"
  resolves = ["sjkaliski/go-github-actions/lint@v0.4.0"]
}

action "sjkaliski/go-github-actions/lint@v0.4.0" {
  uses = "sjkaliski/go-github-actions/lint@v0.4.0"
  secrets = ["GITHUB_TOKEN"]
}
