team: engineering-enablement
pipeline: concourse-webhook-server

feature_toggles:
- update-pipeline

tasks:
- type: docker-compose
  name: build
  save_artifacts: [.]

- type: deploy-cf
  name: deploy
  api: ((cloudfoundry.api-snpaas))
  space: halfpipe
  deploy_artifact: .
  vars:
    GITHUB_SECRET: ((concourse-webhook-server.github_secret))
