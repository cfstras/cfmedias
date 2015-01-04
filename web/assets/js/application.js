window.App = Em.Application.create();

App.Router.map(function() {
});

App.ApplicationView = Em.View.extend({
    classNames: ['height'],
});

App.IndexRoute = Em.Route.extend({
    model: function() {
        var res = this.controllerFor('api').get('list', {
            query: "artist LIKE 'dead%'"
        });
        return Em.ArrayController.create(res);
    },
});

App.ApiController = Em.Controller.extend({
    authToken: null,

    get: function(url, options) {
        return $.getJSON(apiurl+url, options);
    },
})
