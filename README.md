# HomeHive Crawler

**HomeHive Crawler** is an advanced real estate data scraping and filtering tool designed to streamline the process of discovering, managing, and analyzing property listings. With a focus on efficiency and scalability, this project serves as a foundation for building data-driven real estate platforms.

## Key Features

- **Real Estate Data Aggregation**: Scrapes property listings from various sources (e.g., Divar, Sheypoor).
- **Powerful Filtering**: Supports dynamic, multi-criteria filtering (e.g., price range, city, neighborhood, bedrooms, etc.).
- **Telegram Bot Integration**: Provides an intuitive interface for users to interact, search, and view listings directly from Telegram.
- **Database Management**: Uses PostgreSQL for storing and managing property data with optimized models.
- **Dynamic Updates**: Detects and updates existing property records to ensure accurate and up-to-date information.
- **Scalability**: Designed with modular architecture for easy feature extension and deployment.

## Tech Stack

- **Backend**: Go (Golang) <span><img src="https://img.shields.io/badge/Golang-1.23-blue" /></span>
- **Database**: PostgreSQL <span><img src="https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white" /></span>
- **Bot Framework**: Telegram Bot API  <span><img src="https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white" /></span>
- **Deployment**: Docker-friendly for seamless integration <span><img src="https://img.shields.io/badge/Docker-2CA5E0?style=flat&logo=docker&logoColor=white" /></span>

## How It Works

1. **Data Crawling**: Gathers property listings in real-time or from provided datasets.
2. **Data Storage**: Saves data in PostgreSQL, ensuring it’s well-structured and ready for queries.
3. **User Interaction**: Allows users to filter properties via Telegram bot commands and view results instantly.
4. **Updates & Notifications**: Keeps listings current and can notify users of matching properties.

## Getting Started

Clone the repository, configure your `.env` file with database credentials and Telegram bot token, and run the project with:

```bash
docker-compose up --build
```
Let **HomeHive Crawler** take the complexity out of property hunting, making it smarter and simpler for everyone! 🚀



## ERD diagram
[HERE](erd_maket.pdf)

## Telegram Bot

User can communicate with our telegram bot to get ads

[Telegram Bot 🔗](https://t.me/quera11_bot)

