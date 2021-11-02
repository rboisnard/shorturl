def call() {
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

        stage("build") {
          sh "pipeline/build.sh ${BRANCH_NAME}"
        }

        stage("test") {
          sh "pipeline/test.sh ${BRANCH_NAME}"
        }

        stage("publish") {
          // TODO: different tags depending on pull-request, normal branch, release branch ...
          String tag = "staging"

          sh "pipeline/publish.sh ${BRANCH_NAME} ${tag} ${config.arch} raspi01:5000"
        }
      }
    }
  }

  parallel(actions)

  node("x86_slave") {
    stage("manifest") {
      unstash "repo"

      // TODO: different tags depending on pull-request, normal branch, release branch ...
      // same as in stage "publish"
      String tag = "staging"

      sh "pipeline/manifest.sh ${tag} raspi01:5000 amd64,arm64"
    }
  }
}

return this
