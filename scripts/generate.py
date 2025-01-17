#!/usr/bin/env python

import argparse
import os
import string
import random
import polars as pl


def parse_size(size_str):
    units = {"B": 1, "K": 1024, "KB": 1024, "M": 1024**2, "MB": 1024**2, "G": 1024**3, "GB": 1024**3, "T": 1024**4, "TB": 1024**4}
    size_str = size_str.upper().strip()
    for unit in sorted(units.keys(), key=len, reverse=True):  # Match longest unit first
        if size_str.endswith(unit):
            number_part = size_str[: -len(unit)].strip()
            try:
                return int(float(number_part) * units[unit])
            except ValueError:
                raise ValueError(f"Invalid size format: {size_str}")
    return int(size_str)


def generate_large_file(file_size, line_size, output_file):
    chars = string.ascii_letters + string.digits
    lines_count = file_size // (line_size + 1)  # +1 for newline character

    print(f"Generating file: {output_file}")
    with open(output_file, 'w') as f:
        for _ in range(lines_count):
            line = ''.join(random.choices(chars, k=line_size))
            f.write(line + '\n')
    print(f"Done. Generated file: {output_file}")


def sort_large_file(input_file, output_file):
    print(f"Reading and sorting file: {input_file}")

    # Load the data as a DataFrame
    df = pl.read_csv(input_file, has_header=False, separator='\n', new_columns=['line'])

    # Sort the data lexicographically
    sorted_df = df.sort('line')

    # Save the sorted data back to a file
    sorted_df.write_csv(output_file, include_header=False, separator='\n')

    os.remove(input_file)

    print(f"Done. Sorted file saved to: {output_file}")
    print(f"File size: {os.path.getsize(output_file) / 1024**2:.2f} MB")


def main():
    parser = argparse.ArgumentParser(description="Generate and sort a large text file with Polars.")
    parser.add_argument("total_size", help="Total file size (e.g., 10GB, 1MB)")
    parser.add_argument("line_size", help="Size of each line (e.g., 10MB, 512KB)")

    args = parser.parse_args()
    total_size = parse_size(args.total_size)
    line_size = parse_size(args.line_size)

    if total_size < 1 or line_size < 1 or line_size > total_size:
        raise ValueError("Invalid sizes. Ensure total_size >= line_size and both > 0.")

    output_file = f"test_{args.total_size}_{args.line_size}.txt.tmp"
    sorted_output_file = f"test_{args.total_size}_{args.line_size}.txt"

    if not os.path.exists(output_file) or os.path.getsize(output_file) != total_size:
        generate_large_file(total_size, line_size, output_file)

    # Sort the generated file
    sort_large_file(output_file, sorted_output_file)


if __name__ == "__main__":
    main()
