var reporter = require('cucumber-html-reporter');

var options = {
    brandTitle: "SES",
    name: "mPower App Integration Test Suite",
    theme: 'bootstrap',
    // theme: 'hierarchy',
    // theme: 'foundation',
    jsonDir: '/results',
    output: '/reports/report.html',
    reportSuiteAsScenarios: true,
    ignoreBadJsonFile: true,
    launchReport: false,
    metadata: {
        "Test Environment": "Docker",
    }
};

reporter.generate(options);