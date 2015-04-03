module.exports = function(grunt) {
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    concat: {
      options: {
        separator: ";",
      },
      dist: {
        src: [
          'app/components/foundation/foundation.js',
          'app/components/lodash/lodash.js',
          'app/components/jquery/dist/jquery.js',
          'app/components/react/react.js',
         ],
         dest: 'dist/bundle.js',
      },
    },
    sass: {
      options: {
        loadPath: ['app/components/foundation/scss']
      },
      dist: {
        options: {
          sourcemap: 'none',
          style: 'nested'
        },
        files: [{
          expand: true,
          cwd: 'src/scss',
          src: ['*.scss'],
          dest: 'dist/assets/css',
          ext: '.css'
        }],
      }
    },
    'babel': {
        options: {
            sourceMap: false
        },
        dist: {
            files: {
                'dist/index.js': 'src/index.js.jsx',
                'dist/app.js': 'src/app.js.jsx'
            }
        }
    }
  });
  grunt.loadNpmTasks("grunt-contrib-concat");
  grunt.loadNpmTasks("grunt-contrib-sass");
  grunt.loadNpmTasks("grunt-babel");

  grunt.registerTask('build', ['concat', 'babel']);
}

