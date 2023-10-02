# Tiny URL

---

## Local Development Setup

Ensure you docker service is up and running on your machine 

Run following commands to set up the code base locally:
1. `make setup`
2. `make copy-config`
3. `make build`
4. `make setup-local-infra`

Run tests with coverage locally:
1. `make test-cover-report`

Run application locally:
1. `make start-application`

Run service as a docker image:
1. `docker build -t tiny-url-app .`
2. `docker run -p 8181:8181 --env-file .env tiny-url-app`
