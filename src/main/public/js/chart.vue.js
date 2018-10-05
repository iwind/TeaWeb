Tea.context(function () {
    var that = this;

    this.CHART = {
        "id": ""
    };

    this.CHART.id = function (options) {
        if (typeof(options["id"]) == "string") {
            return "chart-" + options["id"];
        }
        return "chart-" + Math.random().toString().replace(".", "");
    };

    this.CHART.progressBar = function (options) {
        return '<div class="chart-box progress">' +
            '   <div class="ui progress blue tiny">' +
            '       <div class="bar" style="width:' + options.value.toString() + '%"></div>' +
            '       <div class="label">' + options.name + " <em>(" + options.detail + ')</em></div>' +
            '   </div>' +
            '</div>';
    };

    this.CHART.line = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chart = echarts.init(document.getElementById(chartId));
            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: -4,
                    x: "center",
                    textStyle: {
                        fontSize: 12,
                        fontWeight: "bold",
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    }
                },
                legend: {
                    data: options.lines.$map(function (_, line) {
                        return line.name;
                    }),
                    bottom: -10,
                    y: "bottom",
                    textStyle: {
                        fontSize: 10
                    }
                },
                xAxis: {
                    data: options.labels
                },
                axisLabel: {
                    formatter: function (v) {
                        return v;
                    },
                    textStyle: {
                        fontSize: 10
                    }
                },
                yAxis: {},
                series: options.lines.$map(function (_, line) {
                    return {
                        name: line.name,
                        type: 'line',
                        data: line.values,
                        lineStyle: {
                            width: 1.2
                        },
                    };
                }),
                grid: {
                    left: 30,
                    right: 0,
                    bottom: 50,
                    top: 20
                },
                axisPointer: {
                    show: false
                },
                tooltip: {
                    formatter: 'X:{b0} Y:{c0}',
                    show: false
                },
                animation: false
            };

            chart.setOption(option);
        });
        return '<div class="chart-box line" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.pie = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chart = echarts.init(document.getElementById(chartId));

            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: 1,
                    bottom: 10,
                    x: "center",
                    textStyle: {
                        fontSize: 12,
                        fontWeight: "bold",
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    }
                },
                legend: {
                    orient: 'vertical',
                    x: 'right',
                    y: 'center',
                    data: options.labels,
                    itemWidth: 6,
                    itemHeight: 6,
                    textStyle: {
                        fontSize: 10
                    }
                },
                xAxis: {
                    data: []
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'pie',
                    data: options.values.$map(function (k, v) {
                        return {
                            name: options.labels.$get(k),
                            value: v
                        };
                    }),
                    radius: ['68%', '75%'],
                    center: ['50%', '56%'],
                    label: {
                        normal: {
                            show: false,
                            position: 'center'
                        },
                        emphasis: {
                            show: false,
                            textStyle: {
                                fontSize: '30',
                                fontWeight: 'bold'
                            }
                        }
                    }
                }],

                grid: {
                    left: -2,
                    right: 0,
                    bottom: 0,
                    top: 0
                },
                axisPointer: {
                    show: false
                },

                tooltip: {
                    formatter: 'X:{b0} Y:{c0}',
                    show: false
                },
                animation: false
            };

            chart.setOption(option);
        });
        return '<div class="chart-box pie" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.gauge = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chart = echarts.init(document.getElementById(chartId));

            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: 1,
                    bottom: 0,
                    x: "center",
                    textStyle: {
                        fontSize: 12,
                        fontWeight: "bold",
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    }
                },
                legend: {
                    data: [""]
                },
                xAxis: {
                    data: []
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'gauge',
                    min: 0,
                    max: 100,

                    data: [
                        {
                            "name": options.detail,
                            "value": options.value
                        }
                    ],
                    radius: "80%",
                    center: ["50%", "60%"],

                    splitNumber: 5,
                    splitLine: {
                        length: 6
                    },

                    axisLine: {
                        lineStyle: {
                            width: 8
                        }
                    },
                    axisTick: {
                        show: true
                    },
                    axisLabel: {
                        formatter: function (v) {
                            return v + options.unit
                        },
                        textStyle: {
                            fontSize: 8
                        }
                    },
                    detail: {
                        formatter: function (v) {
                            return v + options.unit
                        },
                        textStyle: {
                            fontSize: 12
                        }
                    },

                    pointer: {
                        width: 2
                    }
                }],

                grid: {
                    left: -2,
                    right: 0,
                    bottom: 0,
                    top: 0
                },
                axisPointer: {
                    show: false
                },
                tooltip: {
                    formatter: 'X:{b0} Y:{c0}',
                    show: false
                },
                animation: true
            };

            chart.setOption(option);
        });
        return '<div class="chart-box gauge" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.table = function (options) {
        var chartId = this.id(options);
        var chartBox = Tea.element("#" + chartId);
        var s = '<table class="ui table selectable"><thead><tr><th colspan="' + ((options.rows.length > 0) ? options.rows[0].columns.length : 1) + '">' + options.name  + '</th></tr></thead>';
        for (var i = 0; i < options.rows.length; i++) {
            s += "<tr>";
            for (var j = 0; j < options.rows[i].columns.length; j++) {
                var column = options.rows[i].columns[j];
                if (column.width > 0) {
                    s += "<td width=\"" + column.width + "%\">" + column.text + "</td>";
                } else {
                    s += "<td>" + column.text + "</td>";
                }
            }
            s += "</tr>";
        }
        s += '</table>';

        if (chartBox.length > 0) {
            chartBox.html(s);
        } else {
            s = '<div class="chart-box table" id="' + chartId + '">' + s + '</div>';
        }

        return s;
    };

    this.CHART.updateWidgetGroups = function (newGroups) {
        if (that.widgetGroups == null || that.widgetGroups.length == 0) {
            that.widgetGroups = newGroups;
            return;
        }

        newGroups.$each(function (_, group) {
            group.widgets.$each(function (_, widget) {
                widget.charts.$each(function (_, chart) {
                    switch (chart.type) {
                        case "line":
                            that.CHART.line(chart);
                            break;
                        case "progressBar":
                            that.CHART.progressBar(chart);
                            break;
                        case "pie":
                            that.CHART.pie(chart);
                            break;
                        case "gauge":
                            that.CHART.gauge(chart);
                            break;
                        case "table":
                            that.CHART.table(chart);
                            break;
                    }
                });
            });
        });
    };
});