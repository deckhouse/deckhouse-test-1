#!/bin/bash
# FAKE SECRETS FOR TESTING ONLY - NOT REAL CREDENTIALS!
# This file is designed to trigger Gitleaks for testing purposes

echo "Testing Gitleaks detection with fake secrets"

# Fake AWS credentials (triggers aws-access-key-id rule)
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

# Fake GitHub token (triggers github-pat rule)  
export GITHUB_TOKEN="ghp_1234567890abcdef1234567890abcdef123456"

# Fake generic API key (triggers generic-api-key rule)
export API_KEY="sk_test_1234567890abcdef1234567890abcdef"

# Fake Slack webhook URL
export SLACK_WEBHOOK="https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"

echo "These are FAKE credentials for testing Gitleaks detection"
echo "This should trigger multiple secret detection rules"