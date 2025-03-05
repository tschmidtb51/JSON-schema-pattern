use std::env;
use std::fs;
use std::process;
use serde::Deserialize;
use regex::Regex;

#[derive(Deserialize)]
struct Schema {
    properties: Properties,
}

#[derive(Deserialize)]
struct Properties {
    test_values: TestValues,
}

#[derive(Deserialize)]
struct TestValues {
    items: Items,
}

#[derive(Deserialize)]
struct Items {
    properties: ValueProperties,
}

#[derive(Deserialize)]
struct ValueProperties {
    value: PatternValue,
}

#[derive(Deserialize)]
struct PatternValue {
    pattern: String,
}

#[derive(Deserialize)]
struct TestFile {
    test_values: Vec<TestValue>,
}

#[derive(Deserialize, Debug)]
struct TestValue {
    value: String,
    assertion: bool,
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 3 {
        eprintln!("Usage: {} <json-file-path> <schema-file-path>", args[0]);
        process::exit(1);
    }

    let file_path = &args[1];
    let schema_path = &args[2];

    let file_data = fs::read_to_string(file_path).expect("Error reading JSON file");
    let schema_data = fs::read_to_string(schema_path).expect("Error reading schema file");

    let test_file: TestFile = serde_json::from_str(&file_data).expect("Error parsing JSON file");
    let schema: Schema = serde_json::from_str(&schema_data).expect("Error parsing schema file");

    let pattern = &schema.properties.test_values.items.properties.value.pattern;
    let regex = Regex::new(pattern).expect("Error compiling regex pattern");

    // Use Debug trait (":?") to properly show the raw pattern, avoiding the backslash being escaped
    println!("Current regex to test:\n{:?}", pattern);

    let mut failed = false;

    for element in test_file.test_values {
        let result = regex.is_match(&element.value);
        failed |= result != element.assertion;
        println!("{}\t{:?}", result, element);
    }

    println!("\nOverall result:\t{}", if failed { "Error occurred" } else { "Test successful" });

    if failed {
        process::exit(1);
    }
}