#!/usr/bin/env python3
"""
Configuration file with hardcoded secrets for gitleaks testing
DO NOT USE IN PRODUCTION - THESE ARE FAKE SECRETS FOR TESTING
"""

import os
import requests

# Hardcoded secrets that gitleaks should detect
class Config:
    # AWS credentials
    AWS_ACCESS_KEY = "AKIAIOSFODNN7EXAMPLE"
    AWS_SECRET_KEY = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    
    # Database password
    DB_PASSWORD = "MyS3cr3tP@ssw0rd123!"
    
    # API tokens
    GITHUB_API_TOKEN = "ghp_abcdef1234567890abcdef1234567890abcdef12"
    SLACK_BOT_TOKEN = "xoxb-1234567890-1234567890-abcdefghijklmnop12345678"
    
    # Encryption key
    ENCRYPTION_KEY = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"

def connect_to_database():
    """Connect to database using hardcoded credentials"""
    connection_string = f"postgresql://admin:{Config.DB_PASSWORD}@localhost:5432/testdb"
    return connection_string

def make_api_request():
    """Make API request with hardcoded token"""
    headers = {
        "Authorization": f"Bearer {Config.GITHUB_API_TOKEN}",
        "User-Agent": "TestApp/1.0"
    }
    return headers

# Another way secrets might leak
api_key = "sk-1234567890abcdef1234567890abcdef12345678"
webhook_url = "https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz1234567890"