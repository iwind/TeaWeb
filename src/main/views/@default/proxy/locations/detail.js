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
        } else {
            rewrite.proxy = "";
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
});