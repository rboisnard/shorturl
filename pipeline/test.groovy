def call() {
  stage("test") {
    sh """
      docker run --rm -d -p 5500:5500 -e PORT=5500 --name=shorturl_worker_${BRANCH_NAME} raspi01:5000/shorturl/worker:staging
      shorturl_worker_ip=\$(docker exec shorturl_worker_${BRANCH_NAME} hostname -i)
      echo
      echo "### logs ###"
      docker logs shorturl_worker_${BRANCH_NAME}
      echo
      echo "### curl ###"
      curl -v \${shorturl_worker_ip}:5500
      docker stop shorturl_worker_${BRANCH_NAME}
    """
  }
}

return this
