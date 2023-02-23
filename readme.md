# UCSR Analytics

CLI program that lets you download csgo demos from faceit and analyze them for statistics.

## Requirements

This program requires at least go 1.18 to run. You can download the latest version of Go [here](https://golang.org/).

## Usage

### Download

1. Create `demos` folder in root, and folder for each map in `demos` folder.
2. Create `.env` file and put `FACEIT_API_TOKEN` inside.
3. Open `cmd` folder with `cd cmd`.
4. Run `go run . download`.

The program will automatically download your last 100 demos and put them in folders by map.

### Analyze

1. Download demos.
2. Open `cmd` folder with `cd cmd`.
3. Run `go run . analyze`.

The program will parse and analyze your demos for stats and write results to `stats.json`. You can then open `index.html` for data visualization.

## Project structure

```
📦ucsr-analytics
 ┣ 📂cmd
 ┃ ┣ 📜main.go              # entry point for golang script
 ┃ ┣ 📜download.go          # automatic downloading of FACEIT demos
 ┃ ┗ 📜analyze.go           # analyzing FACEIT demos and extracting statistics
 ┣ 📜index.html             # visualization of stats.json
 ┗ 📜stats.json             # file with stats auto generated by golang script
```
