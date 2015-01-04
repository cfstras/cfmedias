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
            sourceMapIncludeSources: true,
          },
          files: {
            'assets/vendor/js/libs.min.js': [
              'bower_components/jquery/dist/jquery.js',
              'bower_components/handlebars/handlebars.js',
              'bower_components/ember/ember.js',
              'bower_components/ember-data/ember-data.js',
              'bower_components/moment/moment.js',
              'bower_components/bootstrap/dist/js/bootstrap.js',
            ],
          }
        },
      },
      less: {
        build: {
          files: {
            'assets/vendor/css/libs.min.css': [
              'bower_components/bootstrap/less/bootstrap.less',
              'bower_components/fontawesome/less/font-awesome.less'
            ]
          }
        }
      },
      copy: {
        build: {
          files: [
            {expand: true, cwd: 'bower_components/bootstrap/dist/fonts/', src: ['**'],
            dest: 'assets/vendor/fonts/'},
            {expand: true, cwd: 'bower_components/fontawesome/fonts/', src: ['**'],
            dest: 'assets/vendor/fonts/'},
          ]
        }
      }
    });

    grunt.loadNpmTasks('grunt-contrib-less');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-ember-templates');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-copy');

    // Default task.
    grunt.registerTask('default', ['uglify', 'less', 'copy']);
  };
