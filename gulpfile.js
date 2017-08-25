var gulp  = require('gulp');
var exec  = require('gulp-exec');
var child = require('child_process');
var color = require('gulp-color');
var chalk = require('chalk');

var cmd   = 'go test';
var files = '*.go';

var logger = {
    error: function(err) {
        console.log(chalk.red(err));
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

var handleConsoleOutput = function (err, stdout) {
    if (err === null) {
        logger.success(stdout);
    } else {
        logger.error(stdout);
    }
};

gulp.task('test', function() {
    child.exec(cmd, function(err, stdout) {
        logger.line();
        handleConsoleOutput(err, stdout);
    });
});
gulp.task('watch:test', ['test'], function() {

    gulp.watch(files, function () {
        gulp.src(files)
            .pipe(exec(cmd, function(err, stdout) {
                logger.line();
                handleConsoleOutput(err, stdout);
            }));
    });
});
gulp.task('dev', ['watch:test']);
gulp.task('default', ['dev']);