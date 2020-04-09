from sys import stderr, stdin, stdout


class Parser:
    name = "S2S"
    server_tty = stdin
    walls = set()
    agents_cor = {}
    agents_col = {}
    box_cor = {}
    box_col = {}
    current_line = 0

    def __init__(self):
        self.cases = {
            "colors": self.gotColors,
            "initial": self.gotInitialLine,
            "goal": self.gotGoals,
        }

    def send_message(self, message):
        print(message, flush=True)

    def start(self):
        self.send_message(self.name)
        mode = ""
        while True:
            # for line in :
            line = self.server_tty.readline().rstrip()
            if line == "#end":
                break
            if line.startswith("#"):
                mode = line[1:]
                self.current_line = 0
                continue
            self.cases[mode](line)

    def gotColors(self, line):
        [color, letters] = line.split(":")
        for letter in letters:
            if letter.isdigit():
                self.agents_col[letter] = color
                continue
            self.box_col[letter] = color

    def gotGoals(self, line):
        pass

    # TODO: clean this up
    def gotInitialLine(self, line):
        for i, char in enumerate(line):
            if char == " ":
                continue

            current_cor = (self.current_line, i)
            if char == "+":
                self.walls.add(current_cor)
            elif char.isdigit():
                if not self.agents_cor.get(char, False):
                    self.agents_cor[char] = set()
                self.agents_cor[char].add(current_cor)
            else:
                if not self.box_cor.get(char, False):
                    self.box_cor[char] = set()
                self.box_cor[char].add(current_cor)


if __name__ == "__main__":
    p = Parser()
    p.start()
