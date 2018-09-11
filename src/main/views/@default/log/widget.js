Tea.context(function () {
    var that = this;

    this.qps = 0;
    this.inputBandwidth = "-";
    this.outputBandwidth = "-";

   this.refresh = function () {
       Tea.action(".widget")
           .post()
           .success(function (response) {
                this.qps = response.data.qps;

                var input = response.data.inputBandwidth;
                if (input <= 0) {
                    input = "-";
                } else if (input < 1024) {
                    input = input.toString() + "B"
                } else if (input < 1024 * 1024) {
                    input = (Math.ceil(input * 100 / 1024) / 100).toString() + "K"
                } else {
                    input = (Math.ceil(input * 100 / 1024/ 1024) / 100).toString() + "M"
                }
                this.inputBandwidth = input;

                var output = response.data.outputBandwidth;
                if (output <= 0) {
                    output = "-";
                } else if (output < 1024) {
                   output = output.toString() + "B"
                } else if (output < 1024 * 1024) {
                   output = (Math.ceil(output * 100 / 1024) / 100).toString() + "K"
                } else {
                   output = (Math.ceil(output * 100 / 1024/ 1024) / 100).toString() + "M"
                }
                this.outputBandwidth = output;
           })
           .done(function () {
               setTimeout(function () {
                    that.refresh();
               }, 1000);
           });
   };

   this.refresh();
});