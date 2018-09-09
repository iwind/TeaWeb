Tea.context(function () {
    this.names = [{
        "key": Tea.key()
    }];
    this.listenArray = [ {
        "key": Tea.key()
    } ];
    this.backendsArray = [{
        "key": Tea.key()
    }];

    this.addName = function () {
        this.names.push({
            "key": Tea.key()
        });
    };

    this.removeName = function (index) {
        this.names.$remove(index);
    };

    this.addListen = function () {
        this.listenArray.push({
            "key": Tea.key()
        });
    };

    this.removeListen = function (index) {
       this.listenArray.$remove(index);
    };

    this.addBackend = function () {
        this.backendsArray.push({
            "key": Tea.key()
        });
    };

    this.removeBackend = function (index) {
        this.backendsArray.$remove(index);
    };

});