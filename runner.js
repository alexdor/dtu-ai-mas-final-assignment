const fs = require("fs");
const { spawn } = require("child_process");
const log = require("single-line-log").stdout;
const mdTable = require("markdown-table");
const github = require("@actions/github");
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
    child.stdout.on("data", (data) => {
      childOutput += `${data}`;
    });
    child.on("error", (error) => {
      console.log(error);
    });
    child.on("exit", (code) => {
      if (shouldOutputContinuously) clearInterval(timer);

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
    { align: ["c", "c", "c"] }
  );

  res += "\n\n<details>\n<summary> Levels </summary>\n";
  res += mdTable([Object.keys(levels[0]), ...levels.map(Object.values)], {
    align: ["c", "c", "c", "c"],
  });
  res += "\n</details>";

  return res;
}

function commentResultsOnPr() {
  if (!process.env.CI) return;

  try {
    const github_token = process.env.GITHUB_TOKEN;
    const { context } = github;
    const octokit = new github.GitHub(github_token);

    const commentParams = {
      ...context.repo,
      issue_number: (
        context.payload.issue ||
        context.payload.pull_request ||
        context.payload
      ).number,
      body: getResultsAsMarkdown(context.action),
    };
    console.log("params", commentParams);
    return octokit.issues
      .createComment(commentParams)
      .then((res) => console.log(res.status, res.data))
      .catch((e) => core.setFailed(e));
  } catch (e) {
    core.setFailed(e);
  }
}

process.on("SIGINT", () => {
  printResults();
  commentResultsOnPr();
  process.exit(2);
});

process.on("exit", (code) => {
  const isSigIntCode = code === 2;
  if (isSigIntCode) return;

  printResults();
  commentResultsOnPr();
});

main();
