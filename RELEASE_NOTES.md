# Initial Release - BMKG Weather Alert Bot (v1.0.0)

This is the initial release of the **BMKG Weather Alert Bot**, a lightweight Telegram bot application written in Go designed to monitor and distribute real-time weather alerts from the Indonesian Agency for Meteorology, Climatology, and Geophysics (BMKG).

This version operates as a **single-target broadcaster**, making it suitable for personal use or dedicated weather monitoring groups.

### Key Features

*   **Real-time Monitoring**: Automatically polls the BMKG RSS feed at configurable intervals (default: every 3 minutes) to detect new weather warnings.
*   **Visual Alert Support**: Capable of retrieving and distributing official infographic maps directly from the CAP (Common Alerting Protocol) data source when available.
*   **Province Filtering**: Includes a configuration option via environment variables to filter notifications based on specific provinces.
*   **Reliable Delivery**: Implements a direct image upload mechanism to ensure successful delivery of infographics, bypassing potential external link issues in Telegram.
*   **Lightweight Persistence**: Utilizes a file-based storage system (`history.json`) to track processed alerts and prevent duplicate notifications without requiring a heavy database engine.
*   **Docker Integration**: Fully containerized with a multi-stage Dockerfile and Docker Compose configuration for easy deployment.

### Current Limitations

*   **Single Broadcast Target**: Currently supports broadcasting only to a single destination (`TARGET_CHAT_ID`) defined in the configuration.
*   **No Interactive Subscription**: Does not yet support multi-user management or dynamic subscription commands via chat interface.

### Deployment

1.  Configure the `.env` file with your Telegram Bot Token and Target Chat ID.
2.  Deploy using Docker Compose:
    ```bash
    docker-compose up -d
    ```
