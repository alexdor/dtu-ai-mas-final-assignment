const fs = require("fs");
const { spawn } = require("child_process");
const log = require("single-line-log").stdout;
const mdTable = require("markdown-table");
const { context, GitHub } = require("@actions/github");
const core = require("@actions/core");

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
  .option("outputMode", {
    alias: "o",
    describe: "output mode, ",
  })
  .help("h", "Show help")
  .alias("help", "h").argv;

const timeout = argv.timeout || 180;
const levelsDir = argv.levels || "./levels";
const command = argv.command || "java -cp target/classes: Client";
const prefixToIgnore = argv.ignorePrefix;
const shouldOutputContinuously = argv.outputMode === "continuous";

const results = {
  total: 0,
  solved: 0,
  failed: 0,
  levels: [],
};

function parseRegexp(text, regexp) {
  const match = text.match(regexp);
  return match ? match[1] : "Error parsing output";
}

function main() {
  let levels = [levelsDir];
  try {
    levels = fs
      .readdirSync(levelsDir)
      .filter((lvl) => !prefixToIgnore || !lvl.startsWith(prefixToIgnore));
  } catch {}

  results.total = levels.length;

  function runThatLevel(levelIndex) {
    const level = levels[levelIndex];
    const child = spawn("java", [
      "-jar",
      "server.jar",
      "-c",
      command,
      "-l",
      level !== levelsDir ? `${levelsDir}/${level}` : levelsDir,
      "-t",
      timeout,
    ]);

    let childOutput = "";
    let solved = null;

    child.stderr.on("data", (data) => {
      childOutput += `${data}`;
    });
    child.stdout.on("data", (data) => {
      childOutput += `${data}`;
    });
    child.on("error", (error) => {
      console.log(error);
    });
    child.on("exit", (code) => {
      if (shouldOutputContinuously) clearInterval(timer);

      const isSuccessCode = code === 0;

      solved =
        isSuccessCode &&
        childOutput.includes("[server][info] Level solved: Yes.");
      solved ? results.solved++ : results.failed++;
      results.levels.push({
        level,
        status: solved ? "✅" : "❌",
        actions: solved
          ? parseRegexp(childOutput, "Actions used: (\\d+)")
          : null,
        time: solved
          ? parseRegexp(
              childOutput,
              "Last action time: ([+-]?([0-9]*[.])?[0-9]+)"
            )
          : null,
        statesExpanded: parseRegexp(
          childOutput,
          solved
            ? "Goal was found after exploring (\\d+) states"
            : "Explored (\\d+) states"
        ),
        statesGenerated: parseRegexp(childOutput, "Generated (\\d+) states"),
      });

      log.clear();

      if (!shouldOutputContinuously) {
        logStatus(
          {
            ...results,
            currentlvl: level,
          },
          shouldOutputContinuously
        );
      }

      const isFinished = levelIndex === levels.length - 1;
      if (!isFinished) runThatLevel(++levelIndex);
    });

    if (!shouldOutputContinuously) return;

    let count = 0;

    const timer = setInterval(() => {
      count += 0.1;
      logStatus(
        {
          ...results,
          currentlvl: level,
          time: count.toFixed(1),
        },
        shouldOutputContinuously
      );
    }, 100);
  }

  runThatLevel(0);
}

function logStatus(status, isContinuousOutputSet) {
  const runMessage = isContinuousOutputSet ? "Currently running" : "Just ran";
  const timeMessage = isContinuousOutputSet ? `[${status.time} s]` : "";
  let logLine = `🏃 ${runMessage} ${status.currentlvl} ${timeMessage}\n`;
  logLine += `✅ Number of solved levels ${status.solved}\n`;
  logLine += `❌ Number of failed levels ${status.failed}\n`;
  logLine += `⏳ Levels left ${status.total - (status.solved + status.failed)}`;

  log(logLine);
}

const resultsHeader = "🧾🧾🧾 The results are in 🧾🧾🧾";

function printResults() {
  console.log("\n");
  console.log(resultsHeader);
  const { total, solved, failed } = results;
  console.table({ total, solved, failed });
  console.table(results.levels);
}

function getResultsAsMarkdown(actionName) {
  let res = `#### ${resultsHeader} \n\n Results for ${process.env.GITHUB_SHA}, from action ${actionName}\n\n`;
  const { total, solved, failed, levels } = results;
  res += mdTable(
    [
      ["total", "solved", "failed"],
      [total, solved, failed],
    ],
    { align: "c" }
  );
  res += "\n\n";
  res += mdTable([Object.keys(levels[0]), ...levels.map(Object.values)], {
    align: "c",
  });

  return res;
}

function commentResultsOnPr() {
  if (!process.env.CI) return;

  const github_token = process.env.GITHUB_TOKEN;
  const octokit = new GitHub(github_token);

  return octokit.issues.createComment({
    ...context.repo,
    issue_number: (context.payload.pull_request || context.payload).number,
    body: getResultsAsMarkdown(context.workflow),
  });
}
let runCleanup = true;
async function cleanup() {
  if (!runCleanup) {
    return;
  }

  runCleanup = false;
  try {
    printResults();
    await commentResultsOnPr();
  } catch (e) {
    core.setFailed(e);
  }
}

process.on("SIGINT", async () => {
  cleanup();
  process.exit(2);
});

process.on("beforeExit", cleanup);

main();
