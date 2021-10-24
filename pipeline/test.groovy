def call() {
  stage("test") {
    docker.image("shorturl/worker:staging").inside() {
      sh """
        echo "not implemented yet"
      """
    }
  }
}

return this
