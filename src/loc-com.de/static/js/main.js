/**
 * @license Copyright 2013 Robin Gloster (LocCom) <robin@loc-com.de>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

    // Backbone.sync = function (method, model, success, error) {
    //     success();
    // };

    var FeedItem = Backbone.Model.extend({
    });

    var FeedItems = Backbone.Collection.extend({
        model: FeedItem
    });

    var Feed = Backbone.Model.extend({
        'url': '/feed/'
    });

    var Feeds = Backbone.Collection.extend({
        model: Feed
    });

    var FeedItemView = Backbone.View.extend({
        initialize: function () {
            _.bindAll(this, 'render'); // every function that uses 'this' as the current object should be in here
            this.feedItemTmpl = Handlebars.compile($("#feed-item-tmpl").html());
        },
        render: function () {
            this.el = this.feedItemTmpl(this.model.toJSON());
            return this; // for chainable calls, like .render().el
        }
    });

    var FeedView = Backbone.View.extend({
        el: $('#content'),
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
            this.feedTmpl = Handlebars.compile($("#feed-tmpl").html());
        },
        render: function () {
            this.el = this.feedTmpl(this.model.toJSON());
            return this; // for chainable calls, like .render().el
        }
    });

    var FeedsNavView = Backbone.View.extend({
        el: $('#leftnav ul'),
        events: {
            'click #addfeed': 'addFeedDialog'
        },
        initialize: function () {
            _.bindAll(this, 'render', 'appendFeed', 'addFeed', 'addFeedDialog'); // every function that uses 'this' as the current object should be in here
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
        },
        addFeedDialog: function () {
            $("#add-feed-dialog").dialog("show");
        },
        addFeed: function (url) {
            this.collection.create({
                url: url
            }, {
                wait: true
            });
        }
    });

    var feedsView;

    $.getJSON('/feeds/1', function (msg) {
        var feeds = new Feeds();
        feeds.reset(msg);
        feedsView = new FeedsNavView({
            collection: feeds
        });
    });

    $('#add-feed-dialog').dialog({
        title: "Feed hinzufügen",
        submit: function () {
            feedsView.addFeed($('#add-feed-url').val());
        }
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
