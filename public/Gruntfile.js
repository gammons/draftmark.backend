module.exports = function(grunt) {
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    concat: {
      options: {
        separator: ";",
      },
      vendor: {
        src: [
          'app/components/foundation/foundation.js',
          'app/components/lodash/lodash.js',
          'app/components/jquery/dist/jquery.js',
          'app/components/director/build/director.min.js',
          'app/components/marked/marked.min.js',
          'app/components/react/react.js',

         ],
         dest: 'dist/vendor.js',
      },
      app: {
        src: [
          'tmp/index.js',
          'tmp/note.js',
          'tmp/note_card.js',
          'tmp/app.js'
        ],
        dest: 'dist/app.js'
      }
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
                'tmp/index.js': 'src/index.js.jsx',
                'tmp/app.js': 'src/app.js.jsx',
                'tmp/note.js': 'src/note.js.jsx',
                'tmp/note_card.js': 'src/note_card.js.jsx'
            }
        }
    }
  });
  grunt.loadNpmTasks("grunt-contrib-concat");
  grunt.loadNpmTasks("grunt-contrib-sass");
  grunt.loadNpmTasks("grunt-babel");

  grunt.registerTask('build', ['babel', 'concat']);
}

