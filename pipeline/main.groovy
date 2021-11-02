def call() {
  def build = load "pipeline/build.groovy"
  def test = load "pipeline/test.groovy"
  def publish = load "pipeline/publish.groovy"
  def manifest = load "pipeline/manifest.groovy"

  def actions = [:]
  stash name: "repo"

  archs = [
    "x86": ["node": "x86_slave", "arch": "amd64"],
    "arm": ["node": "arm_slave", "arch": "arm64"]
  ]

  archs.each{ arch, config ->
    actions[arch] = {
      node(config.node) {
        unstash "repo"
        build()
        test()
        publish(config.arch)
      }
    }
  }

  parallel(actions)

  node("x86_slave") {
    // TODO: different tags depending on pull-request, normal branch, release branch ...
    // same as in publish()
    manifest("staging")
  }
}

return this
