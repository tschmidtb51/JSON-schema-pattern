#!/usr/bin/env python3

import sys
import json
import re

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 <script-path>.py json/file/path.json schema/file/path.json")
        sys.exit(1)

    file_path = sys.argv[1]
    schema_path = sys.argv[2]

    with open(file_path, 'r', encoding='utf-8') as file:
        json_data = json.load(file)

    with open(schema_path, 'r', encoding='utf-8') as schema_file:
        schema_data = json.load(schema_file)

    pattern = schema_data['properties']['test_values']['items']['properties']['value']['pattern']
    regex = re.compile(pattern)

    print("Current regex to test:\n", pattern)

    failed = False

    for element in json_data['test_values']:
        result = bool(regex.match(element['value']))
        failed |= (result != element['assertion'])
        print(result, "\t", element)

    print("\nOverall result:\t", "Error occurred" if failed else "Test successful")

    if failed:
        sys.exit(1)

if __name__ == "__main__":
    main()