Tea.context(function () {
    this.switchOn = function () {
        this.location.on = !this.location.on;

        if (this.location.on) {
            this.$post(".on")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".off")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.reverse = function () {
        this.location.reverse = !this.location.reverse;

        if (this.location.reverse) {
            this.$post(".updateReverse")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "reverse": 1
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".updateReverse")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "reverse": 0
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.switchCaseInsensitive = function () {
        this.location.caseInsensitive = !this.location.caseInsensitive;

        if (this.location.caseInsensitive) {
            this.$post(".updateCaseInsensitive")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "caseInsensitive": 1
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".updateCaseInsensitive")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "caseInsensitive": 0
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    /**
     * 重写规则
     */
    this.rewriteAdding = false;
    this.addingPattern = "";
    this.addingReplace = "";
    this.targetType = "url";
    this.proxyId = "";

    this.location.rewrite = this.location.rewrite.$map(function (k, rewrite) {
        if (/^proxy:\/\//.test(rewrite.replace)) {
            var index = rewrite.replace.indexOf("/", "proxy://".length);
            rewrite.proxy = rewrite.replace.substring(0, index);
            rewrite.replace = rewrite.replace.substring(index);
            rewrite.type = "proxy";
            rewrite.proxyId = rewrite.proxy.substr("proxy://".length);
        } else {
            rewrite.proxy = "";
            rewrite.type = "url";
            rewrite.proxyId = "";
        }
        return rewrite;
    });

    this.addRewrite = function () {
        this.rewriteAdding = !this.rewriteAdding;
    };

    this.cancelRewrite = function () {
        this.rewriteAdding = false;
    };

    this.saveRewrite = function () {
        this.$post("/proxy/rewrite/add")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "pattern": this.addingPattern,
                "replace": this.addingReplace,
                "targetType": this.targetType,
                "proxyId": this.proxyId
            });
    };

    this.deleteRewrite = function (index) {
        if (!window.confirm("确定要删除此重写规则吗？")) {
            return;
        }
        this.$post("/proxy/rewrite/delete")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "rewriteIndex": index
            });
    };

    /**
     * 修改重写规则
     */
    this.editRewrite = function (rewrite, index) {
        var that = this;
        this.location.rewrite.$each(function (k, v) {
            if (k == index) {
                if (typeof(rewrite.isEditing) == "undefined") {
                    rewrite.isEditing = true;
                } else {
                    rewrite.isEditing = !rewrite.isEditing;
                }
                Tea.Vue.$set(that.location.rewrite, index, rewrite);
            } else {
                v.isEditing = false;
                Tea.Vue.$set(that.location.rewrite, k, v);
            }
        });
    };

    this.switchRewriteIndex = function (rewrite, index) {
        rewrite.on = !rewrite.on;
        if (rewrite.on) {
            this.$post("/proxy/rewrite/on")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "rewriteIndex": index
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post("/proxy/rewrite/off")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "rewriteIndex": index
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.cancelEditRewrite = function (rewrite, index) {
        rewrite.isEditing = false;
        Tea.Vue.$set(this.location.rewrite, index, rewrite);
    };

    this.updateRewrite = function (rewrite, index) {
        this.$post("/proxy/rewrite/update")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "rewriteIndex": index,
                "pattern": rewrite.pattern,
                "replace": rewrite.replace,
                "targetType": rewrite.type,
                "proxyId": rewrite.proxyId
            });
    };
});