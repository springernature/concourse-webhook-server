team: engineering-enablement
pipeline: concourse-webhook-server
slack_channel: "#halfpipe-alerts"

feature_toggles:
- update-pipeline

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
    CONCOURSE_DB_HOST: ((concourse-db.host))
    CONCOURSE_DB_USERNAME: ((concourse-db.username_read))
    CONCOURSE_DB_PASSWORD: ((concourse-db.password_read))
