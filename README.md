# Velocity Limits Validator 💸

## About The Project 📙

In finance, it's common for accounts to have so-called "velocity limits". This program accepts or declines attempts to load funds into customers' accounts in real-time. The program will read load attempts (line by line) from `input.txt` and save single-line JSON payloads for each load attempt in `generated-output.txt`.

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

If a load ID is observed more than once for a particular user, all but the first instance will be ignored.

## Running The Project 🏃🏽‍♀️

 1. If you haven't already, install Go: "https://golang.org/doc/instal"
 2. Clone this repo using `git clone https://github.com/luminaire-dev/velocity-limits.git` or `git@github.com:luminaire-dev/velocity-limits.git` 
 3. Open terminal or cmd prompt and cd into the `velocity-limits` directory
 4. run `go run .`
 5. view generated_output.txt for results.
