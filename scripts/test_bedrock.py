import boto3, json, os

# Load env vars
from pathlib import Path
env_file = Path(__file__).parent / ".env"
if env_file.exists():
    for line in env_file.read_text().splitlines():
        line = line.strip()
        if line and not line.startswith("#") and "=" in line:
            key, val = line.split("=", 1)
            os.environ.setdefault(key.strip(), val.strip())

client = boto3.client('bedrock-runtime', region_name='eu-west-2')
body = json.dumps({
    'anthropic_version': 'bedrock-2023-05-31',
    'max_tokens': 100,
    'messages': [{'role': 'user', 'content': [{'type': 'text', 'text': 'Reply with just the word YES'}]}]
})
resp = client.invoke_model(modelId='anthropic.claude-3-haiku-20240307-v1:0', body=body, contentType='application/json')
result = json.loads(resp['body'].read())
print('SUCCESS:', result['content'][0]['text'])
