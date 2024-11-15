#!/usr/bin/env deno run --allow-read

/**
 * @file Script to validate JSON files against given schema
 *
 * Usage: deno run --allow-read <script-path>.ts json/file/path.json schema/file/path.json
 */

const [filePath, schemaPath] = Deno.args;

if (!filePath || !schemaPath) {
    console.error("Usage: deno run --allow-read <script-path>.ts json/file/path.json schema/file/path.json");
    Deno.exit(1);
}

const decoder = new TextDecoder("utf-8");

const jsonFile = await Deno.readFile(filePath);
const schemaFile = await Deno.readFile(schemaPath);

const json = JSON.parse(decoder.decode(jsonFile));
const schema = JSON.parse(decoder.decode(schemaFile));

const pattern = schema.properties.test_values.items.properties.value.pattern;
const r = new RegExp(pattern);

console.log('Current regex to test:', '\n', pattern);

let failed = false;

json.test_values.forEach((element: { value: string; assertion: boolean }) => {
    const result = r.test(element.value);
    failed = failed || (result !== element.assertion);
    console.log(result, '\t', element);
});

console.log("\nOverall result", '\t', failed ? "Error occurred" : "Test successful");

if (failed) {
    Deno.exit(1);
}