# resdiary-slot-finder
Scraper to find a slot for my pregnant wifes birthday lol

## Environment Variables

The application is configured using environment variables. Below is a list of available variables:

- `DISABLE_PUSHOVER`: (Optional) Set to `true` to disable Pushover notifications. Defaults to `false`.
- `PUSHOVER_API_KEY`: (Required) Your Pushover application API key.
- `PUSHOVER_RECIPIENT`: (Required) Your Pushover user key.
- `RESERVATION_DATE`: (Optional) The date for which to find reservations, in `YYYY-MM-DD` format. Defaults to `2024-12-21`.
- `RESTAURANT_NAMES`: (Optional) A comma-separated list of restaurant identifiers (as used in the resdiary.com URL) to check for availability. Defaults to `ChesilRectory`.
  - Example: `RESTAURANT_NAMES="ChesilRectory,ThePigNearBath,RickSteinSandbanks"`
- `RESTAURANT_COVERS`: (Optional) The number of covers (people) for the reservation. Defaults to `2`.
- `RESERVATION_IGNORE_THRESHOLD_HOUR`: (Optional) Ignore slots that are after this hour (e.g., `21` for 9 PM). Defaults to `21`.
- `RESERVATION_IGNORE_THRESHOLD_MINUTE`: (Optional) Ignore slots that are after this minute when `RESERVATION_IGNORE_THRESHOLD_HOUR` is also set. Defaults to `0`.

## Running the application

You can run the application using Docker or directly with Go.

### Docker
```bash
docker build -t resdiary-slot-finder .
docker run -e PUSHOVER_API_KEY="your_api_key" -e PUSHOVER_RECIPIENT="your_recipient_key" resdiary-slot-finder
```

### Go
```bash
go run main.go
```
Make sure to set the required environment variables in your shell before running.
