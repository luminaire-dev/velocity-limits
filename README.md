# Velocity Limits Validator 💸

Be kind, this is my first time using Go 🙇🏻‍♀️

## About The Project 📙

In finance, it's common for accounts to have so-called "velocity limits". This program accepts or declines attempts to load funds into customers' accounts in real-time. The program will read load attempts, line by line from `input.txt` and save a single-line JSON payload for each load attempt in `generated-output.txt`.

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

If a load ID is observed more than once for a particular user, all but the first instance will be ignored.

## Running The Project 🏃🏽‍♀️

 1. If you haven't already, install Go: "https://golang.org/doc/instal"
 2. Clone this repo using `git clone https://github.com/luminaire-dev/velocity-limits.git` or `git@github.com:luminaire-dev/velocity-limits.git` 
 3. Open terminal or cmd prompt and cd into the `velocity-limits` directory


Run the project:

```sh
make run
```

Run the unit tests:

```sh
make test
```

Compare the generated output  agaisnt the expected output using `diff`. (No output means the files match ✔️ )

```sh
make comapre
```

## Project Structure

```
.
├── processor                     
│   ├── processor.go               # Business logic - approves or rejects the incoming load based on limits validation
│   ├── processor_test.go          # Unit tests 
│   └── models.go                  # Contains common strucs used by processor.go and processor_test.go
├── main.go                        # Reads input line by line and generates output file
├── main_test.go                   # Unit tests 
├── input.txt                      # Contains all load attemps
├── output.txt                     # Expected output
├── generated_output.txt           # output generated by the program
└── Makefile                       # Makefile
```

## Assumptions

All input data (in input.txt) is valid and formated in porper JSON objects 
