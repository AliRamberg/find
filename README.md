# Find
This tool is a simple alternative to the [look](https://git.kernel.org/pub/scm/utils/util-linux/util-linux.git/tree/misc-utils/look.c) command from util-linux package. It returns the first line that matches the pattern or is lexicographically greater than the pattern.


## Search for a pattern
Searching for a pattern in a file is as simple as running the following command:

```bash
finder <pattern> <file>
```

## Build
To build the project, use the following command:

```bash
go build -ldflags="-s -w" -o ~/.local/bin/finder ./pkg
```

## Generate sample files
In order to generate an ordered file of random numbers and ascii characters, you can use the `generate.py` script.

```bash
./scripts/generate.py <file_size> <line_size>
```

- `file_size` is the size of the file in bytes in human readable format (e.g. 1KB, 1MB, 1GB)
- `line_size` is the size of each line in bytes in human readable format (e.g. 1KB, 1MB, 1GB)


## Examples
In this example we're going to search for a line that already exists in the file.
```bash
> sed -n -e 3100000p scripts/test_15GB_5KB.txt > scripts/input.txt
> head -c10 scripts/input.txt 
z75DKBQ0nS

> cat scripts/input.txt | go finder -f scripts/test_15GB_5KB.txt --term - 
Result: z75DKBQ0nS...pZpUDoHC8K

# Using a third-party tool to find the line number
> rg -no "^z75DKBQ0nS" scripts/test_15GB_5KB.txt
3100000:z75DKBQ0nS # As expected, the line number is 3100000
```

Now, let's add just one character to the line and search for it, and we should get the next line.
```bash
> echo -n a >> scripts/input.txt
> cat scripts/input.txt | go finder -f scripts/test_15GB_5KB.txt --term - 
Result: z75Mgqd5w6...2vF0cCO4Ul

> rg -no "^z75Mgqd5w6" scripts/test_15GB_5KB.txt
3100001:z75Mgqd5w6 # We got the next line
```