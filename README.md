# Concourse Webhook Server

<a href="https://concourse.halfpipe.io/teams/engineering-enablement/pipelines/concourse-webhook-server"><img src="http://badger.halfpipe.io/engineering-enablement/concourse-webhook-server" title="build status"></a>

A service that receives webhook events to trigger pipelines.

The idea is to use this service as an additional way of triggering resource checks instead of Concourse polling every minute. Polling this often puts a lot of load on Concourse - across the db, web nodes and workers. With this service resource checks will be faster and we can reduce the frequency of the polling.

* pipeline jobs will be triggered immediately instead of waiting for the next timed check.
* the default resource check interval can be increased (to e.g. 10m), reducing load in Concourse.
* it is isolated and non-intrusive. Nothing relies on it working. If it is down, resources will still check, just more slowly.


### How it works

At the moment it supports GitHub `push` events.

When GitHub sends an event the service..
* reads the repository and branch from the event
* looks in Concourse for any active git resources matching repository and branch
* triggers `check resource`

#### GitHub webhook config

A webhook needs to be created in a GitHub organisation
```
Payload URL:     https://<this-service>/github
Content Type:    json
Secret:          $GITHUB_SECRET
Trigger Events:  Just the push event
```


#### Development

local `./build` and `./run`

or in docker `./build docker` and `./run docker`


#### Logs

<https://kibana.snpaas.eu/s/ee/app/kibana#/discover/2b5edfa0-cbf7-11ea-865f-a1aeec51368c>
