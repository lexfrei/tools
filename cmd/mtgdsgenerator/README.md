# mtgdsgenerator

Magic: The Gathering dataset generator that downloads card images from Scryfall.

## Features

- Downloads bulk card data from Scryfall API
- Supports filtering by language
- Downloads high-resolution card images
- Parallel downloads with configurable workers
- Progress tracking for bulk data download
- Thread-safe concurrent downloads
- Handles double-faced cards

## Installation

```bash
go build -o mtgdsgenerator ./cmd/mtgdsgenerator
```

## Usage

### Basic usage

Download all cards in all languages:

```bash
./mtgdsgenerator
```

### Filter by language

Download only English cards:

```bash
./mtgdsgenerator --lang en
```

Download only Russian cards:

```bash
./mtgdsgenerator --lang ru
```

Supported language codes: `en`, `es`, `fr`, `de`, `it`, `pt`, `ja`, `ko`, `ru`, `zhs`, `zht`, etc.

### Configure parallel downloads

Set number of parallel image downloaders (default: 10):

```bash
./mtgdsgenerator --parallel 20
```

### Specify data type

Choose which Scryfall bulk data to download (default: all_cards):

```bash
./mtgdsgenerator --datatype oracle_cards
```

Available data types:
- `all_cards` - All card printings
- `oracle_cards` - Unique cards (Oracle text)
- `unique_artwork` - Cards with unique artwork
- `default_cards` - Latest printing of each card
- `rulings` - Card rulings

### Combined options

```bash
./mtgdsgenerator --lang en --parallel 20 --datatype oracle_cards
```

## Output

The tool creates:
- `./images/{set}/{lang}/{card_id}.jpg` - Downloaded card images
- `./cards.json` - JSON mapping of card names to image paths

## Example output structure

```
images/
├── mkm/
│   └── en/
│       ├── 00a1234b-5678-90cd-ef12-345678901234.jpg
│       └── 11b2345c-6789-01de-f234-567890123456.jpg
└── lci/
    ├── en/
    │   └── 22c3456d-7890-12ef-3456-789012345678.jpg
    └── ru/
        └── 33d4567e-8901-23f0-4567-890123456789.jpg
cards.json
```

## Performance

- Bulk data download: ~300MB, takes 1-3 minutes depending on connection
- Image downloads: Configurable parallelism (default: 10 workers)
- Progress tracking updates every 500ms

## Error handling

- HTTP timeouts: 30 seconds for images, 5 minutes for bulk data
- Failed downloads are logged but don't stop the process
- Empty or invalid responses are skipped

## Configuration file

You can optionally use a configuration file:

```bash
./mtgdsgenerator --config ~/.mtgdsgenerator.yaml
```

Config file format (YAML):

```yaml
datatype: oracle_cards
lang: en
parallel: 20
```

## Technical details

- Uses Scryfall API (https://scryfall.com/docs/api)
- Implements thread-safe concurrent downloads with sync.Map
- Respects Scryfall rate limits through controlled parallelism
- Downloads only highres/lowres quality images
