team: engineering-enablement
pipeline: concourse-webhook-server
slack_channel: "#halfpipe-alerts"

triggers:
- type: timer
  cron: "0 5 * * *"

feature_toggles:
- update-pipeline
- docker-decompose

tasks:
- type: docker-compose
  name: build
  save_artifacts:
  - manifest.yml
  - Procfile
  - webserver

- type: deploy-cf
  name: deploy
  api: ((cloudfoundry.api-snpaas))
  space: halfpipe
  deploy_artifact: .
  vars:
    GITHUB_SECRET: ((concourse-webhook-server.github_secret))
    CONCOURSE_ENDPOINT: ((concourse-admin-prod.url))
    CONCOURSE_USERNAME: ((concourse-admin-prod.username))
    CONCOURSE_PASSWORD: ((concourse-admin-prod.password))
    CONCOURSE_DB_HOST: ((halfpipe-concourse-db-prod.host))
    CONCOURSE_DB_USERNAME: ((halfpipe-concourse-db-prod.username_read))
    CONCOURSE_DB_PASSWORD: ((halfpipe-concourse-db-prod.password_read))
