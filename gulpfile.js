var gulp = require('gulp');
var exec = require('gulp-exec');
var child = require('child_process');
var chalk = require('chalk');

var cmdTest = 'go test';
var cmdVet = 'go tool vet */*.go';
var files = '*/*.go';

var logger = {
    error: function (err, trace) {
        if (err === '') {
            console.log(chalk.red(trace));
        } else {
            console.log(chalk.red(err));
            console.log(chalk.gray(trace));
        }
    },
    success: function (msg) {
        console.log(chalk.green(msg));
    },
    line: function () {
        var i = 0, msg = '';
        while (i < 80) {
            i++;
            msg += '-';
        }
        console.log(chalk.gray(msg));
    }
};

var handleConsoleOutput = function (err, stdout, stderr) {
    if (err === null) {
        logger.success(stdout);
    } else {
        logger.error(stdout, stderr);
    }
};

gulp.task('test', function () {
    child.exec(cmdTest, function (err, stdout, stderr) {
        handleConsoleOutput(err, stdout, stderr);
    });
    child.exec(cmdVet, function (err, stdout, stderr) {
        logger.line();
        if (stderr !== '') {
            console.log(chalk.bold.gray('VET'));
        }
        handleConsoleOutput(err, stdout, stderr);
    });
});
gulp.task('watch:test', ['test'], function () {
    gulp.watch(files, function () {
        gulp.src(files)
            .pipe(exec(cmdTest, function (err, stdout, stderr) {
                handleConsoleOutput(err, stdout, stderr);
            }));
        gulp.src(files)
            .pipe(exec(cmdVet, function (err, stdout, stderr) {
                logger.line();
                if (stderr !== '') {
                    console.log(chalk.bold.gray('VET'));
                }
                handleConsoleOutput(err, stdout, stderr);
            }));
    });
});
gulp.task('dev', ['watch:test']);
gulp.task('default', ['dev']);