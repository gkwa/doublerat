import re
import pathlib
import sys


def transform_and_write(line):
    match = re.match(r"^Files (.*) and (.*) differ$", line)
    if match:
        file1, file2 = match.group(1), match.group(2)
        transformed_line = f"diff -uw {file1} {file2}\n"
        with open("test1.sh", "a") as output_file:
            output_file.write(transformed_line)


if __name__ == "__main__":
    pathlib.Path("test1.sh").unlink(missing_ok=True)
    for line in sys.stdin:
        line = line.strip()
        transform_and_write(line)
