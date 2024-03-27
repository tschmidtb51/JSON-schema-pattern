#!/usr/bin/env node

/**
 * @file Script to validate JSON files against given schema
 *
 * Usage: node <script-path>.js json/file/path.json schema/file/path.json
 *
 */

const { exit } = require('process')
const fs = require('fs')

const [, , filePath, schemaPath] = process.argv

const json = JSON.parse(fs.readFileSync(filePath, { encoding: 'utf-8' }))
const schema = JSON.parse(fs.readFileSync(schemaPath, { encoding: 'utf-8' }))

const pattern = schema.properties.test_values.items.properties.value.pattern
const r = new RegExp(pattern)

console.log('Current regex to test:', '\n', pattern)

let failed = false

json.test_values.forEach(element => {
    const result = (r.exec(element.value) != null)
    failed = failed | (result != element.assertion) 
    console.log(result, '\t', element)
});

console.log("\nOverall result", '\t', failed ? "Error occurred" : "Test successful")

if (failed) {
    exit(1)
}