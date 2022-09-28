/**
 * Check a milestone is set for a pull request
 *
 * @param {object} inputs
 * @param {object} inputs.github - A pre-authenticated octokit/rest.js client with pagination plugins.
 * @param {object} inputs.context - An object containing the context of the workflow run.
 * @param {object} inputs.core - A reference to the '@actions/core' package.
 * @returns {Promise<void>}
 */

module.exports.checkPrMilestone = async ({ github, context, core }) => {
    const pr = context.payload.pull_request;

    if (pr.milestone) {
      core.debug(`This pull request has a milestone: ${pr.milestone.title}`);
    } else {
      core.setFailed("The pull request has no milestone. Set a milestone for the pull request.");
    }
}
