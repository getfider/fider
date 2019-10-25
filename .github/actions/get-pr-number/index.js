const core = require('@actions/core');
const github = require('@actions/github');

try {
  core.debug(github.context.payload.pull_request);
} catch (error) {
  core.setFailed(error.message);
}