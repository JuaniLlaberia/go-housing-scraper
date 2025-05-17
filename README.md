# 🏠 Zonaprop Argentina Properties Scraper (Go)

This is an experimental Go project to scrape real estate listings from Zonaprop in Argentina. The tool extracts property data based on given filters, cleans the data, exports it to CSV, and can optionally analyze it using Gemini and create a report.

> ⚠️ **Note**: This is a personal test project. Features and functionality may change frequently or not work 100%. Use at your own risk.

> ⚠️ **Note**: Zonaprop may deny your requests. Try using VPN.

## ✨ Features

- 🔍 Scrape property listings with filters (location, price, type, rooms, bathrooms, etc.)
- 🧼 Clean and format the scraped data
- 📁 Export results to a CSV file
- 🤖 Experimental integration with Gemini for data analysis (Create a report on the scraped properties)

## 📦 Project Setup

Make sure you have Go installed (version 1.20+ recommended).

Clone the repository:

```bash
git clone https://github.com/yourusername/zonaprop-scraper.git
cd zonaprop-scraper
```

Set your Gemini API key:

```bash
.env
GEMINI_API_KEY=<gemini-key>
```

Run the script:

```bash
go run .
```
