var gulp  = require('gulp');
var exec  = require('gulp-exec');
var child = require('child_process');

var cmd   = 'go test';
var files = '*.go';

gulp.task('test', function() {
    child.exec(cmd, function(err, stdout) {
        console.log(stdout);
    });
});
gulp.task('watch:test', function() {
    gulp.watch(files, function () {
        gulp.src(files)
            .pipe(exec(cmd, function(err, stdout) {
                console.log(stdout);
            }));
    });
});
gulp.task('dev', ['watch:test']);
gulp.task('default', ['dev']);