// Copyright 2025 Flant JSC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
export GITHUB_TOKEN=<token>
export CONTEXT_REPO_OWNER=deckhouse
export CONTEXT_REPO_NAME=deckhouse-test-1
export BRANCH="chore/save-logs-from-e2e-clusters-on-error"
export WORKFLOW_NAME="Build and test for dev branches"
export JOB_NAME="Build Deckhouse"
*/
const envConfig = {
  repo_owner: process.env.CONTEXT_REPO_OWNER,
  repo_name: process.env.CONTEXT_REPO_NAME,
  github_token: process.env.GITHUB_TOKEN,
  branch: process.env.BRANCH,
  workflow_name: process.env.WORKFLOW_NAME,
  job_name: process.env.JOB_NAME
};

require('../helpers/types');
const Octokit = require('@actions/github');
const core = require('@actions/core');
const github = Octokit.getOctokit(envConfig.github_token);

/** @type {GithubContext} */
const context = {
  repo: {
    owner: envConfig.repo_owner,
    repo: envConfig.repo_name
  }
};

const { waitForJobInWorkflowIsCompletedWithSuccess, isJobInWorkflowCompleted } = require('./validate-job-in-workflow-is-ready')({
  github,
  context,
  core
});

// async function main() {
//   const result = await isJobInWorkflowCompleted(envConfig.branch, envConfig.workflow_name, envConfig.job_name);
// }

// async function main() {
//   const result = await waitForJobInWorkflowIsCompletedWithSuccess(
//     envConfig.branch,
//     envConfig.workflow_name,
//     envConfig.job_name,
//     100
//   );
// }

const githubAction = require('../helpers/github-actions')({
  github,
  context,
  core
});

async function main() {
  const result = await githubAction.GetPullRequestByBranchName('main');
  console.log(result);
  //   const s = {
  //     payload: {
  //       ref: 'refs/heads/feat/test2'
  //     }
  //   };

  //   const d = githu
  // bAction.GetBranchNameFromContext(s);
  //   console.log(d);
  //   githubAction.GetBranchByPullRequestNumber();
}

main();
