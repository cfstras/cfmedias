window.App = Em.Application.create();

App.Router.map(function() {});

App.ListView = Ember.ListView.extend({
    height: 500,
    rowHeight: 35,
    //elementWidth: 80,
    //width: 500,
    itemViewClass: Ember.ListItemView.extend({
        templateName: "index_row",
    //    tagName: "tr",
    }),

});

App.IndexRoute = Ember.Route.extend({
    model: function() {
        var items = [];
        for (var i = 0; i < 100; i++) {
            items.push({
                col_1: "Item",
                col_2: i,
            });
        }
        return items;
    }
});
