;(function($, document, window, undefined) {
    "use strict";


    $.widget( "loccom.dialog", {
        // default options
        options: {
            submitButton: true,
            cancelButton: true,
            defaultCancel: true,
            title: '',

            // callbacks
            submit: null,
            cancel: null
        },

        _create: function () {
            var dialog = this;

            this.element.addClass("ui-loccom-dialog");

            this.overlay = this.element.wrap($('<div/>', {
                "class": "ui-loccom-overlay"
            })).parent();
            this.overlay.on('click', function () {
                dialog._trigger("cancel");
            });
            this.element.on('click', function (event) {
                event.stopPropagation();
            });

            if (this.options.title !== '') {
                this.title = $('<h2/>', {
                    "text": this.options.title
                });
                this.element.prepend(this.title);
            }
            if (this.options.submitButton || this.options.cancelButton) {
                this.buttons = $('<div/>', {
                    "class": "ui-loccom-dialog-buttons"
                });
                this.element.append(this.buttons);
            }
            if (this.options.submitButton) {
                this.submitButton = $('<input/>', {
                    "class": "ui-loccom-dialog-submit",
                    "type": "button",
                    "value": this.options.submitButton === true ? "Abschicken" :
                            this.options.submitButton
                });
                this.buttons.append(this.submitButton);
            }
            if (this.options.cancelButton) {
                this.cancelButton = $('<input/>', {
                    "class": "ui-loccom-dialog-cancel",
                    "type": "button",
                    "value": this.options.cancelButton === true ? "Abbrechen" :
                            this.options.cancelButton
                });
                this.buttons.append(this.cancelButton);
            }
            this.element.on('dialogcancel', function() {
                if (dialog.options.defaultCancel) {
                    dialog.hide(function () {
                        dialog.empty();
                    });
                }
            });
            this.element.on('dialogsubmit', function() {
                dialog.hide(function () {
                    dialog.empty();
                });
            });

            this.element.find(".ui-loccom-dialog-cancel").on('click', function () {
                dialog._trigger("cancel");
            });
            this.element.find(".ui-loccom-dialog-submit").on('click', function () {
                dialog._trigger("submit");
            });
        },

        show: function (callback) {
            this.element.fadeIn();
            this.overlay.fadeIn(callback);
        },

        hide: function (callback) {
            this.overlay.fadeOut(callback);
        },

        empty: function () {
            this.element.find('input:text, input:password,'+
                'input:file, select, textarea').val('');
            this.element.find('input:radio, input:checkbox')
                .removeAttr('checked').removeAttr('selected');
            this.element.find('.cancel-empty').empty();
        },

        _destroy: function() {
            this.overlay.remove();
        },

        // _setOption is called for each individual option that is changing
        _setOption: function(key, value) {
            // prevent invalid color values
            if (key === "title") {
                this.title.text(value);
            }
            if (key === "cancelButton") {
                this.cancelButton.val(value);
            }
            if (key === "submitButton") {
                this.submitButton.val(value);
            }
            // in 1.9 would use _super
            $.Widget.prototype._setOption.call(this, key, value);
        }
    });

})(jQuery, document, window);

