Tea.context(function () {
   this.servers = this.servers.$map(function (k, server) {
       server.config.backends = server.config.backends.$map(function (_, backend) {
           return backend.address;
       });
       return server;
   });
});