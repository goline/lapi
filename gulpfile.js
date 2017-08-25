var gulp  = require('gulp');
var exec  = require('gulp-exec');
var child = require('child_process');
var chalk = require('chalk');

var cmd   = 'go test';
var files = '*.go';

var logger = {
    error: function(err, trace) {
        console.log(chalk.red(err));
        console.log(chalk.gray(trace));
    },
    success: function (msg) {
        console.log(chalk.green(msg));
    },
    line: function () {
        var i = 0, msg = '';
        while (i < 80) { i++; msg += '-'; }
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

gulp.task('test', function() {
    child.exec(cmd, function(err, stdout, stderr) {
        logger.line();
        handleConsoleOutput(err, stdout, stderr);
    });
});
gulp.task('watch:test', ['test'], function() {

    gulp.watch(files, function () {
        gulp.src(files)
            .pipe(exec(cmd, function(err, stdout, stderr) {
                logger.line();
                handleConsoleOutput(err, stdout, stderr);
            }));
    });
});
gulp.task('dev', ['watch:test']);
gulp.task('default', ['dev']);