Tea.context(function () {
    this.logs = [];
    this.fromId = -1;
    this.total = 0;
    this.countSuccess = 0;
    this.countFail = 0;
    this.qps = 0;

    Tea.delay(function () {
        this.loadLogs();
    });

    this.loadLogs = function () {
        Tea.action(".get")
            .params({
                "fromId": this.fromId,
                "size": 100
            })
            .success(function (response) {
                this.total = Math.ceil(response.data.total * 100 / 1000) / 100;
                this.countSuccess = Math.ceil(response.data.countSuccess * 100 / 1000) / 100;
                this.countFail = Math.ceil(response.data.countFail * 100 / 1000) / 100;
                this.qps = response.data.qps;

                this.logs = response.data.logs.concat(this.logs);
                this.logs.$each(function (_, v) {
                    if (typeof(v["isOpen"]) === "undefined") {
                        v.isOpen = false;
                    }
                });

                if (this.logs.length > 0) {
                    this.fromId = this.logs.$first().id;

                    if (this.logs.length > 100) {
                        this.logs = this.logs.slice(0, 100);
                    }
                }
            })
            .done(function () {
                // 每1秒刷新一次
                Tea.delay(function () {
                    this.loadLogs();
                }, 1000)
            })
            .get();
    };

    this.showLog = function (index) {
        var log = this.logs[index];
        log.isOpen = !log.isOpen;

        // 由于Vue的限制直接设置 log.isOpen 并不起作用
        Tea.Vue.$set(Tea.Vue.logs, index, log);
        //Tea.Vue.$forceUpdate();
    };

    this.formatCost = function (seconds) {
        var s = (seconds * 1000).toString();
        var pieces = s.split(".");
        if (pieces.length < 2) {
            return s;
        }

        return pieces[0] + "." + pieces[1].substr(0, 3);
    }
});