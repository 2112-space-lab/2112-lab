var cukemerge = require('cucumber-json-merge');

files = cukemerge.listJsonFiles('/suites', true);
merged = cukemerge.mergeFiles(files);
console.log(merged)
cukemerge.writeMergedFile('/results/godog-cucumber.json', merged, true)
