# NEVERAI

This project is the implementation of a Single and Multi Agent system for transportation tasks, this project was developed for the purpose of DTU's course "02285 AI and MAS".

This project uses A* as the main searching algorithm.

## Initial setup

1. Instal [Go](https://golang.org/dl/), [make](https://www.gnu.org/software/make/) and [Java](https://java.com/en/download/) (Java is only used to run the server jar file, none of the code is in Java)
2. `$ git clone git@github.com:alexdor/dtu-ai-and-mas.git`
3. In order to run an example level you can run `make start-gui`
4. Optional: In order to be able to use the runner to print aggregated tables, you need to install [Node.JS](https://nodejs.org/en/download/) and then run `npm i` inside the folder of the project

## Running the levels

The makefile includes multiple options to run a level. The main ones are:

* `make start-gui`: run a level and display the results in the graphical user interface, provided by the server
* `make start`: run one or multiple levels and display the final information on the terminal
* `make runner`: run one or multiple levels and display the information to the terminal (with nice summaries and tables, this command requires)

By default the app is going to use the level at `levels/custom_levels/SASimple.lvl`, in order to specify a different level (or a folder including multiple levels) every command accepts the parameter `level=path_to_file`. For example, to execute the runner on all the levels in the folder `levels/new_levels` you can run `make runner level=levels/new_levels/`
