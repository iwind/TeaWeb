Tea.context(function () {
    var that = this;
    this.widgetHeight = parseInt((window.innerWidth - 220) * 0.75 / 3);

    window.addEventListener("resize", function () {
        that.widgetHeight = parseInt((window.innerWidth - 220) * 0.75 / 3);
    });
});