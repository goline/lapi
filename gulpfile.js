var gulp  = require('gulp');
var exec  = require('gulp-exec');
var child = require('child_process');
var color = require('gulp-color');

var cmd   = 'go test';
var files = '*.go';

gulp.task('test', function() {
    child.exec(cmd, function(err, stdout, stderr) {
        console.log(stdout, stderr);
    });
});
gulp.task('watch:test', ['test'], function() {
    gulp.watch(files, function () {
        gulp.src(files)
            .pipe(exec(cmd, function(err, stdout, stderr) {
                console.log(color("---------------------------------------------------------------------------------", 'MAGENTA'));
                console.log(stdout, stderr);
            }));
    });
});
gulp.task('dev', ['watch:test']);
gulp.task('default', ['dev']);