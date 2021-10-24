@Library("jenkins-ci-lib") _

node {
  stage("checkout") {
    local_checkout()
  }

  def local_pipeline = load "pipeline/main.groovy"
  local_pipeline()
}
