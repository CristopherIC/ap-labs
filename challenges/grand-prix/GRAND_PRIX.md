Instructions
==================================

Requeriments
----------------------
- A computer that can run Go programs

Considerations
----------------------
- You may need to adjust the size of the terminal/console window to make the outputs look good

Build/Run Instructions
----------------------

- Build: You can use either of the following options:
    - `make build`
    - `go build grand-prix.go`

- Run: You can use either of the following options after building the code:
    - `make run` (This will run the program with default parameters)
    - `make run r=<Number of racers> l=<Number of laps>` (This will run the program with specified parameters)
    - `./grand-prix -racers <Number of racers> -laps<Number of laps>`(This will run the program with specified parameters)