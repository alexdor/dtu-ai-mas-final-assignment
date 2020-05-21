const fs = require("fs");
const { spawn } = require("child_process");
const log = require("single-line-log").stdout;
const argv = require("yargs")
  .option("timeout", {
    alias: "t",
    describe: "timeout for each",
  })
  .option("levels", {
    alias: "l",
    describe: "path to levels folder",
  })
  .option("command", {
    alias: "c",
    describe: "command used to run client",
  })
  .option("ignorePrefix", {
    alias: "i",
    describe: "ignore levels where the filaname start with this prefix",
  })
  .help("h", "Show help")
  .alias("help", "h").argv;

const timeout = argv.timeout || 180;
const levelsDir = argv.levels || "./levels";
const command = argv.command || "java -cp target/classes: Client";
const prefixToIgnore = argv.ignorePrefix;

const results = {
  total: 0,
  solved: 0,
  failed: 0,
  levels: [],
};

function main() {
  const levels = fs
    .readdirSync(levelsDir)
    .filter((lvl) => !prefixToIgnore || !lvl.startsWith(prefixToIgnore));
  results.total = levels.length;

  function runThatLevel(levelIndex) {
    const level = levels[levelIndex];
    const child = spawn("java", [
      "-jar",
      "server.jar",
      "-c",
      command,
      "-l",
      `${levelsDir}/${level}`,
      "-t",
      timeout,
    ]);

    let childOutput = "";
    let solved = null;
    child.stdout.on("data", (data) => {
      childOutput += `${data}`;
    });
    child.on("error", (error) => {
      console.log(error);
    });
    child.on("exit", (code) => {
      clearInterval(timer);

      const isSuccessCode = code === 0;
      if (!isSuccessCode) process.exit();

      solved = childOutput.includes("[server][info] Level solved: Yes.");

      solved ? results.solved++ : results.failed++;

      results.levels.push({
        level,
        status: solved ? "✅" : "❌",
        actions: solved ? childOutput.match("Actions used: (\\d+)")[1] : null,
        time: solved
          ? childOutput.match("Last action time: ([+-]?([0-9]*[.])?[0-9]+)")[1]
          : null,
      });

      log.clear();

      const isFinished = levelIndex === levels.length - 1;
      if (!isFinished) runThatLevel(++levelIndex);
    });

    let count = 0;

    const timer = setInterval(() => {
      count += 0.1;
      logStatus({
        ...results,
        currentlvl: level,
        time: count.toFixed(1),
      });
    }, 100);
  }

  runThatLevel(0);
}

function logStatus(status) {
  let logLine = `🏃‍♂️ Currently running ${status.currentlvl} [${status.time} s]\n`;
  logLine += `✅ Number of solved levels ${status.solved}\n`;
  logLine += `❌ Number of failed levels ${status.failed}\n`;
  logLine += `⏳ Levels left ${status.total - (status.solved + status.failed)}`;

  log(logLine);
}

function printResults() {
  console.log("\n");
  console.log("🧾🧾🧾 The results are in 🧾🧾🧾");
  const { total, solved, failed } = results;
  console.table({ total, solved, failed });
  console.table(results.levels);
}

process.on("SIGINT", () => {
  printResults();
  process.exit(2);
});

process.on("exit", (code) => {
  const isSigIntCode = code === 2;
  if (isSigIntCode) return;

  printResults();
});

main();
