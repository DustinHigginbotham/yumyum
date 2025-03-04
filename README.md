# Yum Yum
Crafting the perfect recipe backstory can be time-consuming, and let's be honest, not everyone has a heartwarming tale about artisanal rosemary or a life-changing bowl of soup. But don't worry, we've got you covered!

### What is this?
YumYum is a **wholesome, AI-powered food blog post generator**, except... well, something seems *off*. While it starts as a cheerful and nostalgic food blog, the story subtly unravels into **absurdity and psychological horror**. But don't worry, this is all **perfectly normal**. AI, right?

### Features
- AI-generated recipe backstories that start normal... but don't take that for granted.
- Real-time streaming responses for a more immersive descent into madness.
- Completely unaware narrator who is fine. Everything's fine.
- Rate limiting because the stories are *almost* too fun to generate.

### How it Works
1) Enter a recipe name and the ingredients (one per line)
2) Click "Generate" and let the AI craft a rich, emotionally charged food blog post.
3) Watch as the story spirals into an episode of existential horror.

### Tech Stack
- Frontend: Preact + HTM
- Backend: Go (Cloud Run, OpenAI API)
- Database: Redis (Cloud Memorystore) for rate limiting
- Hosting: Google Cloud Platform (Cloud Run, Cloud Storage)
- Streaming: Server-Sent Events (SSE)

## Getting Started

### Running Locally

You'll need:

- Go 1.24
- Node.js + npm
- Redis
- Task (https://taskfile.dev/installation/)

## Development

1) Set up the environmental variables

Copy the `.task-env.sample` file to `.task-env` and then replace the values as needed.

2) Start the API

```bash
task api
```

3) Start the web client
```bash
task web
```

## Known Issues & Bugs
- There is a **slight** issue where the AI-generated stories become possessed by some unknown force. This is expected behavior.
- If you experience excessive existential dread, try rebooting your computer... three times. Or clear your cache.

## Contributing
Pull requests are welcome! If you’d like to contribute, fork the repo and create a branch.

## License
MIT License. Use at your own risk—we take no responsibility for accidental summoning of entities.