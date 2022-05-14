# shorturl

Welcome here

In this repository we explore the theme of an URL shortener application.

**The rules are:**
* use redis
* write the application in go
* ~~serve multiple workers with NGINX~~ - *not implemented yet*

**Todo:**
- [x] web app skeleton
- [x] connection with redis
- [x] URL shortener logic
- [x] storage on redis
- [ ] unit tests
- [ ] proper versionning
- [ ] connection with NGINX
- [ ] common pipeline/local tooling
- [ ] multiple workers
- [ ] log/trace in goroutine
- [ ] backup/restore DB
- [ ] non-reg tests (test DB)
- [ ] `/u/` for UI, `/a/` for api
- [ ] storage capacity / deletion
- [ ] redis replication
- [ ] multiple pods
- [ ] monitoring / alerting
- [ ] nice UI

**Technical choices:**
* The available Jenkins CI on my private network provides x86 and arm nodes but has a dependency on *docker* for containers, my development environement has *podman* available.
* My private servers are not powerful enough to run a Kubernetes cluster, if you are not satisfied, you are invited to help fund better hardware ;)