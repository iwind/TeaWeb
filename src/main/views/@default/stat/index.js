Tea.context(function () {
    this.dataType = "pv";
    this.dataRange = "daily";
    this.chartTitle = "";

    this.$delay(function () {
        this.loadChart();

        var that = this;
        window.addEventListener("resize", function () {
            var chart = echarts.init(that.$find(".main-box .chart")[0]);
            chart.resize();
        });
    });

    this.loadChart = function () {
        this.$get("/stat/data?type=" + this.dataType + "&range=" + this.dataRange)
            .success(function (response) {
                this.chartTitle = response.data.title;

                var chart = echarts.init(this.$find(".main-box .chart")[0]);

                // 指定图表的配置项和数据
                var option = {
                    textStyle: {
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    },
                    title: {
                        text: "",
                        top: 0,
                        x: "center"
                    },
                    tooltip: {},
                    legend: {
                        data: []
                    },
                    xAxis: {
                        data: response.data.labels
                    },
                    yAxis: {},
                    series: [{
                        name: '',
                        type: 'line',
                        data: response.data.data,
                        areaStyle: {}
                    }],
                    grid: {
                        left: 50,
                        right: 50,
                        bottom: 50,
                        top: 10
                    },
                    axisPointer: {
                        show: true
                    },
                    tooltip: {
                        formatter: 'X:{b0} Y:{c0}'
                    }
                };

                chart.setOption(option);
            });
    };

    this.changeType = function (dataType) {
        this.dataType = dataType;
        this.loadChart();
    };

    this.changeRange = function (dataRange) {
        this.dataRange = dataRange;
        this.loadChart();
    };

    this.formatPercent = function (num) {
        return Math.ceil(num * 10000) / 100;
    };

    this.formatMS = function (seconds) {
        return Math.ceil(seconds * 10000) / 10;
    };
});