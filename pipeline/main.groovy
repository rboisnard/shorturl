def call() {
  def build = load "pipeline/build.groovy"
  def test = load "pipeline/test.groovy"

  def image_nodes = [:]
  stash name: "repo"

  image_nodes["x86"] = {
    node("x86_slave") {
      unstash "repo"
      build()
      test()
    }
  }

  image_nodes["arm"] = {
    node("arm_slave") {
      unstash "repo"
      build()
      test()
    }
  }

  parallel(image_nodes)
}

return this
