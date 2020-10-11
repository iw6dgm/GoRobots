# GoRobots: Crobots Batch Tournament Manager in Golang

## Install

You can download GoRobots source code and compile it with your current version of the Golang compiler via
```
go build gorobots.go
```
or alternatively download, unpack and install (copy/move into your executable path) the binary from any of the available builds.

GoRobots will automatically detect the number of available cores and use them to run concurrent matches via multiple threads.

This repo doesn't include the robot source codes nor the `crobots` executable: you may want to check the [Crobots official site](https://crobots.deepthought.it) for that.

## Configuration

GoRobots does not need anything special but just a simple YAML configuration file, very similar to the one used by the Crobots Python scripts.
A configuration file must include a list of properties as such as:

* `label` Currently unused but added for compatibility with the Python version. A short name for the tournament. It might be used in future versions to automatically choose the output filename
* `matchF2F` Repetition factor used for the `crobots -m` command line option in face-to-face matches
* `match3VS3` Repetition factor used for the `crobots -m` command line option in 3-vs-3 matches
* `match4VS4` Repetition factor used for the `crobots -m` command line option in 4-vs-4 matches
* `sourcePath` Robots source codes folder. If not needed, you can use `'.'` as it is required and can't be empty. Use path separator either `/` or `\` based on your operating system
* `listRobots` Robots list (filenames) without `.r` or `.ro` extensions, in the format of `['name1', 'name2', ...]`. Filenames can include a path like `tournament/micro/robot1`. Use path separator either `/` or `\` based on your operating system

Robots will be searched by concatenating the `sourcePath` and the robot filename with the path separator. E.g. if on Linux/macOS

`sourcePath` is `'./test/v4'` and `listRobots` contains `'micro/rabbit2'` GoRobots will search for `'./test/v4/micro/rabbit2.ro'` first and, if not found, try to compile `'./test/v4/micro/rabbit2.r'`

On Windows one may use `sourcePath: '.\test\v4'` and `'micro\rabbit2'` as robot filename hence GoRobots will concatenate them using `\` path separator generating `'.\test\v4\micro\rabbit2'`.

Example (Linux/macOS):
```yml
label: '2013'
matchF2F: 1000
match3VS3: 250
match4VS4: 168
sourcePath: '.'
listRobots: ['2013/lamela', '2013/eternity', '2013/ride', '2013/okapi', '2013/pjanic']
```

## Command line options

  * `-bench string`
    	robot (full path, no extension) to create a benchmark tournament for
  * `-config string`
    	YAML configuration file (default "config.yml")
  * `-cpu int`
    	number of threads (CPUs/cores) (default the number of automatically detected cores)
  * `-exe string`
    	Crobots executable (default "crobots")
  * `-limit int`
    	limit random number of matches (random mode only)
  * `-out string`
    	output results to file
  * `-parse string`
    	parse log file only (no tournament)
  * `-random`
    	random mode: generate random matches for 4vs4 only
  * `-sql string`
    	output results as SQL updates to file
  * `-test`
    	test mode, check configuration and exit
  * `-type string`
    	tournament type: f2f, 3vs3 or 4vs4
  * `-verbose`
    	verbose mode: print tournament progression percentage

## Examples

Below a list of common command line examples, using Linux/macOs path separator unless specified otherwise. If not using the binary executable `gorobots` (or `gorobots.exe`) one can run the code by simply calling `go run gorobots.go` therefore in the following examples these two options are perfectly interchangeable.

### Benchmark

Run a benchmark to estimate your PC SGNIPs performance:

`gorobots -type 4vs4 -config ./bench.yml`

and the numner of SGNIPs can be estimated via
`(15 * 200010 * 500) / T` where T is the execution time in seconds (as printed out by GoRobots).

### Crobots executable not in path

If the `crobots` executable is not available in your path, you can specify it via the `-exe` option, e.g.:

`gorobots -type 4vs4 -config ./conf/tournament.yml -exe ~/git/crobots/bin/crobots`

and `~/git/crobots/bin/crobots` executable will be used to run matches.

On Windows machines the following version can be used (assuming `gorobots.exe` is in your path):

`gorobots -type 4vs4 -config .\conf\tournament.yml -exe C:\git\crobots\bin\crobots.exe`

### Test configuration

Test configuration, verifying that the `crobots` executable is available and robots binaries and/or source codes are readable, e.g.:

`gorobots -type 4vs4 -config ./conf/tournament.yml -exe ~/git/crobots/bin/crobots -test`

will test that
* `~/git/crobots/bin/crobots` is runnable
* all robots specified in the configuration have a readable `.ro` binary or a `.r` source code: first time GoRobots doesn't find a `.ro` it will try to compile the `.r` source code and then use the binary for further tournaments (unless manually removed/deleted).

By using the `-test` option GoRobots won't be running the tournament but will exit after the verification has ended.

### Save output to file

Results can be saved into text and SQL files, e.g.:

`gorobots -type 4vs4 -config ./conf/tournament.yml -exe ~/git/crobots/bin/crobots -sql ./output.sql -csv ./output.csv`

will save a text file `output.csv` in TAB separated format and a text file `output.sql` with SQL `UPDATE` statements. For the SQL statements to be usable you may need to install somewhere a SQL-compatible database of your choice (SQLite, MySQL, ...), create results tables via SQL code similar to

```sql
CREATE TABLE `results_f2f` (
 `robot` TEXT NOT NULL,
 `games` INTEGER NOT NULL DEFAULT 0,
 `wins` INTEGER NOT NULL DEFAULT 0,
 `ties` INTEGER NOT NULL DEFAULT 0,
 `points` INTEGER NOT NULL DEFAULT 0,
 PRIMARY KEY (`robot`)
);

CREATE TABLE `results_3vs3` (
 `robot` TEXT NOT NULL,
 `games` INTEGER NOT NULL DEFAULT 0,
 `wins` INTEGER NOT NULL DEFAULT 0,
 `ties` INTEGER NOT NULL DEFAULT 0,
 `points` INTEGER NOT NULL DEFAULT 0,
 PRIMARY KEY (`robot`)
);

CREATE TABLE `results_4vs4` (
 `robot` TEXT NOT NULL,
 `games` INTEGER NOT NULL DEFAULT 0,
 `wins` INTEGER NOT NULL DEFAULT 0,
 `ties` INTEGER NOT NULL DEFAULT 0,
 `points` INTEGER NOT NULL DEFAULT 0,
 PRIMARY KEY (`robot`)
);
```

and insert the robots names first (without extension) before running the updates.

Note: on Windows it is _not_ recommended to redirect the output (e.g. with `> output.csv`) to file but to use the `-csv` option instead as redirection might cause GoRobots to unexpectedely crash.

### Robot benchmarking

The `-bench <robot>` command line option enables to test (benchmark) the performance of a single robot `<robot>` (you can add its full path) against a list of opponents. The robot specified as part of the `-bench` option will face all combinations of the robots provided with the YAML configuration.
E.g.
`listRobots: ['alice', 'bob', 'charlie']`

`gorobots -type 3vs3 -config config.yml -bench zombie`

will generate

```
zombie alice bob
zombie bob charlie
zombie alice charlie
```
matches. Same for the face-to-face and 4-vs-4 modes. GoRobots does not allow to include the tested robot as part of the YAML list of opponents. For instance, in the previous example if providing
`listRobots: ['alice', 'bob', 'charlie', 'zombie']` will result in a error message.

### Random matches

If the list of robots is huge (like in the King of the hill AKA KOTH) it might come in handy running a tournament where matches are randomly generated using the `-random` option and `-limit <n>` to tell GoRobots to stop after `<n>` matches. These options can be used alongside the `-bench <robot>` to test (benchmark) a robot against a randomly generated list of matches.
Note: `-random` option requires to use `limit <n>` and only supports `-type 4vs4`. E.g.:

`gorobots -type 4vs4 -config config.yml -bench test -random -limit 150000`

### Limitations and future developments

Compared to the Python tournament manager, GoRobots lacks any support for databases therefore tournaments cannot be paused nor resumed. This feature did not turn out to be particularly useful in the past - actually quite taxing in terms of performance - and will not probably be implemented in the future.
Future developments might include support for f2f, 3vs3 and 4vs4 tournament modes via the same command line execution (so, no need to call GoRobots three times), for instance using the `label` configuration property to generate output filenames accordingly. This means the behaviour of the `-type`, `-csv` and `-sql` options, one day, may unexpectedly change.