Tea.context(function () {

    /**
     * 后端地址
     */
    this.backendAdding = false;
    this.newBackendAddress = "";
    this.backendEditing = false;

    this.addBackend = function () {
        this.backendAdding = !this.backendAdding;
    };

    this.addBackendSave = function () {
        Tea.action("/proxy/backend/add")
            .params({
                "filename": this.filename,
                "address": this.newBackendAddress
            })
            .post();
    };

    this.editBackendCancel = function (index, backend) {
        backend.isEditing = !backend.isEditing;
        Tea.Vue.$set(this.proxy.backends, index, backend);
    };

    this.editBackend = function (index, backend) {
        backend.isEditing = true;
        this.backendEditing = !this.backendEditing;

        Tea.Vue.$set(this.proxy.backends, index, backend);
    };

    this.editBackendSave = function (index, backend) {
        console.log(backend);
        Tea.action("/proxy/backend/update")
            .params({
                "filename": this.filename,
                "index": index,
                "address": backend.address
            })
            .post();
    };

    this.deleteBackend = function (index) {
        if (!window.confirm("确定要删除此服务器吗？")) {
            return;
        }
        Tea.action("/proxy/backend/delete")
            .params({
                "filename": this.filename,
                "index": index
            })
            .post();
    };

});