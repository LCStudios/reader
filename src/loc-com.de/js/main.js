(function ($, Backbone, document, window, undefined) {
    "use strict";

    // var MyRouter = Backbone.Router.extend({
    //     routes : {
    //         "say/:something" : "say"
    //     },

    //     say : function (something) {
    //         alert(something);
    //     }
    // });

    // var yC = new MyRouter();

    // Backbone.history.start({
    //     pushState: true
    // });

    // yC.navigate("say/something", {trigger: true});

    Backbone.sync = function (method, model, success, error) {
        success();
    };

    var FeedItem = Backbone.Model.extend({
    });

    var FeedItems = Backbone.Collection.extend({
        model: FeedItem
    });

    var FeedItemView = Backbone.View.extend({
        initialize: function () {
            _.bindAll(this, 'render'); // every function that uses 'this' as the current object should be in here
        },

        render: function () {
            var feedItemTmpl = Handlebars.compile($("#feed-item-tmpl").html());
            $(this.el).html(feedItemTmpl(this.model.toJSON()));
            return this; // for chainable calls, like .render().el
        }
    });

    var FeedView = Backbone.View.extend({
        el: $('#content'),
        events: {
            'click button#add': 'addFeedItem'
        },
        initialize: function () {
            _.bindAll(this, 'render', 'addFeedItem', 'appendFeedItem'); // every function that uses 'this' as the current object should be in here

            this.collection.bind('add', this.appendFeedItem); // collection event binder

            this.counter = 0;
            this.render();
        },
        render: function () {
            var self = this;
            _(this.collection.models).each(function (item) { // in case collection is not empty
                self.appendFeedItem(item);
            }, this);
        },
        addFeedItem: function () {
            this.counter++;
            var item = new FeedItem();
            item.set({
                part2: item.get('part2') + this.counter // modify item defaults
            });
            this.collection.add(item);
        },
        appendFeedItem: function (item) {
            var itemView = new FeedItemView({
                model: item
            });
            $(this.el).append(itemView.render().el);
        }
    });

    $.getJSON('/feed', function (msg) {
        var feedItems = new FeedItems();
        feedItems.reset(msg.Feed.Items);
        var feedView = new FeedView({
            collection: feedItems
        });
    });


})(jQuery, Backbone, document, window);
