def call() {
  stage("build") {
    sh "tools/build.sh ${BRANCH_NAME}"
  }
}

return this
