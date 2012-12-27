(function ($, Backbone, _, document, window, undefined) {
    "use strict";

    var MyRouter = Backbone.Router.extend({
        routes : {
            "feed/:id" : "feed"
        },

        feed : function (id) {
            $.getJSON('/feed/' + id, function (msg) {
                var feedItems = new FeedItems();
                feedItems.reset(msg.Feed.Items);
                var feedView = new FeedView({
                    collection: feedItems
                });
            });
        }
    });

    var yC = new MyRouter();

    Backbone.history.start({
        pushState: true
    });

    // yC.navigate("feed/50d706a055acd67a106a4a14", {trigger: true});

    Backbone.sync = function (method, model, success, error) {
        success();
    };

    var FeedItem = Backbone.Model.extend({
    });

    var FeedItems = Backbone.Collection.extend({
        model: FeedItem
    });

    var Feed = Backbone.Model.extend({
    });

    var Feeds = Backbone.Collection.extend({
        model: Feed
    });

    var FeedItemView = Backbone.View.extend({
        initialize: function () {
            _.bindAll(this, 'render'); // every function that uses 'this' as the current object should be in here
        },
        render: function () {
            var feedItemTmpl = Handlebars.compile($("#feed-item-tmpl").html());
            this.el = feedItemTmpl(this.model.toJSON());
            return this; // for chainable calls, like .render().el
        }
    });

    var FeedView = Backbone.View.extend({
        el: $('#content'),
        // events: {
        //     'click button#add': 'addFeedItem'
        // },
        initialize: function () {
            _.bindAll(this, 'render', 'appendFeedItem'); //, 'addFeedItem'); // every function that uses 'this' as the current object should be in here

            this.collection.bind('add', this.appendFeedItem); // collection event binder

            this.counter = 0;
            this.render();
        },
        render: function () {
            var self = this;
            $(this.el).html('');
            _(this.collection.models).each(function (item) { // in case collection is not empty
                self.appendFeedItem(item);
            }, this);
        },
        // addFeedItem: function () {
        //     this.counter++;
        //     var item = new FeedItem();
        //     item.set({
        //         part2: item.get('part2') + this.counter // modify item defaults
        //     });
        //     this.collection.add(item);
        // },
        appendFeedItem: function (item) {
            var itemView = new FeedItemView({
                model: item
            });
            $(this.el).append(itemView.render().el);
        }
    });

    var FeedNavView = Backbone.View.extend({
        initialize: function () {
            _.bindAll(this, 'render'); // every function that uses 'this' as the current object should be in here
        },
        render: function () {
            var feedTmpl = Handlebars.compile($("#feed-tmpl").html());
            this.el = feedTmpl(this.model.toJSON());
            // $(this.el).html(feedTmpl(this.model.toJSON()));
            return this; // for chainable calls, like .render().el
        }
    });

    var FeedsNavView = Backbone.View.extend({
        el: $('#leftnav ul'),
        initialize: function () {
            _.bindAll(this, 'render', 'appendFeed'); // every function that uses 'this' as the current object should be in here
            this.collection.bind('add', this.appendFeed);
            this.render();
        },
        render: function () {
            var self = this;
            _(this.collection.models).each(function (item) { // in case collection is not empty
                self.appendFeed(item);
            }, this);
        },
        appendFeed: function (feed) {
            var feedNavView = new FeedNavView({
                model: feed
            });
            $(this.el).append(feedNavView.render().el);
        }
    });


    $.getJSON('/feeds/1', function (msg) {
        var feeds = new Feeds();
        feeds.reset(msg);
        var feedsView = new FeedsNavView({
            collection: feeds
        });
    });

    $(document).on("click", "a:not([data-bypass])", function (evt) {
        // Get the anchor href and protcol
        var href = $(this).attr("href");
        var protocol = this.protocol + "//";

        // Ensure the protocol is not part of URL, meaning its relative.
        if (href && href.slice(0, protocol.length) !== protocol &&
                href.indexOf("javascript:") !== 0) {
            evt.preventDefault();

            Backbone.history.navigate(href, true);
        }
    });

})(jQuery, Backbone, _, document, window);
