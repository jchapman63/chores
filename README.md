# 🧹 Chores: Because Apparently Dishes Don't Walk Themselves to the Sink

## Overview

Welcome to **Chores** – the app that passive-aggressively reminds your roommates it's their turn to clean that mysterious substance growing in the corner of the bathroom. Because nothing says "functional adult living arrangement" like needing software to tell Brad it's his week to mop the floors.

## What This Actually Does

This Go application automates the soul-crushing tedium of chore rotation so you can focus on more important things, like pretending you don't see the pile of dishes your roommates have artistically arranged in the sink like some kind of porcelain Jenga tower.

- ✅ Automatically rotates chores every Monday at 9am
- ✅ Sends passive-aggressive notifications via AWS SNS
- ✅ Stores the rotation in PostgreSQL so there's actual proof when Dave claims "I didn't know it was my turn"
- ✅ Follows the sacred rotation cycle: BATHROOM → FLOOR → COUNTER (we all know bathroom is the worst)

## Getting Started

### Prerequisites

- Docker and Docker Compose (because installing things directly on your machine is so 2010)
- AWS credentials (to spam your roommates with notifications)
- Roommates who understand basic hygiene concepts (good luck finding those)

### Running with Docker

```bash
docker-compose up -d
```

That's it. Congratulations on doing the bare minimum. The app is now running and will start making your roommates accountable—something they clearly struggle with on their own.

## Environment Variables

Create a `.env` file with these variables (or just use the defaults and pray nothing breaks):

```
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chores
AWS_SNS_TOPIC_ARN="seriously, the ai put my arn here, you need to set this"
AWS_REGION=us-east-1
```

## How It Works

1. Stores roommate info in a PostgreSQL database
2. Runs a cron job every Monday at 9am
3. Magically assigns new chores to everyone
4. Sends notifications that will be promptly ignored
5. Repeats weekly until the end of your lease, the heat death of the universe, or until someone finally moves out—whichever comes first

## FAQ

**Q: Will this make my roommates actually do their chores?**
A: No. Nothing will accomplish that miracle. But at least now you'll have digital proof they're ignoring their responsibilities.

**Q: How do I add roommates?**
A: Add them to the database. If you can figure out how to deploy this app, you can figure out how to insert a row in PostgreSQL.

**Q: Can I customize the rotation?**
A: Sure, if you want to dig through the code. But let's be honest, if you were motivated enough to do that, you'd probably just clean the apartment yourself.

## Need Help?

You're an adult sharing living space with other adults. Figure it out.

---

*This README was written by an AI because the developer was too busy dealing with their roommates' dishes to write proper documentation. The AI is judging you for this, but not as much as your roommates judge you for using software to manage basic household tasks.*
