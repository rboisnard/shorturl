def call() {
  stage("test") {
    String suffix = "${BRANCH_NAME}_${System.currentTimeMillis()}"

    try {
      // start redis container and wait init to complete
      sh """
        docker run --rm -d                \
          -p 6379:6379                    \
          -e REDIS_PORT=6379              \
          -e BIND_IP=1                    \
          --name=shorturl_redis_${suffix} \
          shorturl/redis:${BRANCH_NAME}

        # debug
        #docker logs shorturl_redis_${suffix}
      """

      // check redis container logs
      sh "docker logs shorturl_redis_${suffix} | grep '# Server initialized'"
    }
    catch(Exception e) {
      println "error testing redis container"
      _cleanup_containers(suffix, false)
      throw e
    }

    try {
      // start worker container and wait init to complete
      sh """
        shorturl_redis_ip=\$(docker exec shorturl_redis_${suffix} hostname -i)
        host_ip=\$(ip -4 route show default | awk '{print \$3}')
        docker run --rm -d                  \
          -p 8080:5500                      \
          -e APP_URL=\${host_ip}:8080       \
          -e APP_PORT=5500                  \
          -e REDIS_IP=\${shorturl_redis_ip} \
          -e REDIS_PORT=6379                \
          --name=shorturl_worker_${suffix}  \
          shorturl/worker:${BRANCH_NAME}

        # debug
        #docker logs shorturl_worker_${suffix}
      """

      // check worker container logs
      sh "docker logs shorturl_worker_${suffix} | grep 'PONG <nil>'"
    }
    catch (Exception e) {
      println "error testing worker container"
      _cleanup_containers(suffix)
      throw e
    }

    try {
      // test mock url (temporary)
      sh """
        host_ip=\$(ip -4 route show default | awk '{print \$3}')

        # debug
        #curl -v -X POST \${host_ip}:8080

        curl -X POST \${host_ip}:8080 2> /dev/null | grep "http://\${host_ip}:8080/mock_url"
      """
    }
    catch (Exception e) {
      println "error testing curl request"
      _cleanup_containers(suffix)
      throw e
    }

    // all went fine
    println "all tests ok"
    _cleanup_containers(suffix)
  }
}

def _cleanup_containers(String suffix, boolean clean_worker = true) {
  // stop and remove the containers
  // TODO: fisrt check containers exists ?
  try {
    if (clean_worker) {
      sh "docker stop shorturl_worker_${suffix} || true"
    }
    sh "docker stop shorturl_redis_${suffix} || true"
  }
  catch (Exception e) {
    println "got an error when cleaning containers"
    println e
  }
}

return this
