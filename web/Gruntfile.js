  'use strict';

  module.exports = function(grunt) {

    // Project configuration.
    grunt.initConfig({
      // Metadata.
      pkg: grunt.file.readJSON('package.json'),
      banner: '/*! <%= pkg.title || pkg.name %> - v<%= pkg.version %> - ' +
        '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
        '<%= pkg.homepage ? "* " + pkg.homepage + "\\n" : "" %>' +
        '* Copyright (c) <%= grunt.template.today("yyyy") %> <%= pkg.author.name %>;' +
        ' Licensed <%= _.pluck(pkg.licenses, "type").join(", ") %> */\n',
      /*emberTemplates: {
        build: {
          files: {
          }
        },
      },*/
      uglify: {
        build: {
          options: {
            sourceMap: true,
          },
          files: {
            'assets/vendor/libs.min.js': [
              'bower_components/bootstrap/dist/js/bootstrap.js',
              'bower_components/ember/ember.js',
              'bower_components/ember-data/ember-data.js',
              'bower_components/moment/moment.js'
            ],
          }
        },
      },
      copy: {
        build: {
          files: [
            {expand: true, cwd: 'bower_components/bootstrap/dist/css/', src: ['bootstrap.min.css'],
            dest: 'assets/vendor/css/'},
            {expand: true, cwd: 'bower_components/bootstrap/dist/fonts/', src: ['**'],
            dest: 'assets/vendor/fonts/'},
          ]
        }
      }
    });

    //grunt.loadNpmTasks('grunt-contrib-sass');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-ember-templates');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-copy');

    // Default task.
    grunt.registerTask('default', ['uglify', 'copy']);
  };
