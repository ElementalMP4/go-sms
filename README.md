# Go SMS

Have you ever wanted to text a large language model? Probably not. But here's a mechanism to allow you do so so if you wish.

## Config

```jsonc
{
    "apiBase": "http://192.168.9.1", // Base URL for the Huawei E3372
    "ignore": [
        "1234567890" // List of numbers to ignore
    ],
    "ollamaBase": "http://localhost:11434", // Ollama base URL
    "model": "gemma3:12b", // Model to use
    "systemPrompt": "You are an agent employed to respond to SMS messages." // System prompt
}
```