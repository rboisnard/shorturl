def call() {
  stage("test") {
    sh "tools/test.sh ${BRANCH_NAME}"
  }
}

return this
