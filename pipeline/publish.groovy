def call(String arch) {
  stage("publish") {
    // TODO: different tags depending on pull-request, normal branch, release branch ...
    String tag = "staging"

    sh "tools/publish.sh ${BRANCH_NAME} ${tag} ${arch} raspi01:5000"
  }
}

return this
