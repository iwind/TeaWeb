Tea.context(function () {
    this.proxyDescriptionEditing = false;
    this.editDescription = function () {
        this.proxyDescriptionEditing = true;
    };

    this.editDescriptionSave = function () {
        Tea.action(".updateDescription")
            .params({
                "filename": this.filename,
                "description": this.proxy.description
            })
            .post();
    };

    /**
     * 域名管理
     */
    this.newName = "";
    this.nameAdding = false;
    this.addName = function () {
        this.nameAdding = true;
    };

    this.addNameSave = function () {
        Tea.action(".addName")
            .params({
                "filename": this.filename,
                "name": this.newName
            })
            .post();
    };

    this.editNameIndex = -1;
    this.editName = function (index, name) {
        this.editNameIndex = index;
    };

    this.editNameSave = function (index, name) {
        Tea.action(".updateName").params({
                "filename": this.filename,
                "index": index,
                "name": name
            })
            .post();
    };

    this.editNameCancel = function () {
        this.editNameIndex = -1;
    };

    this.deleteName = function (index) {
        if (!window.confirm("确定要删除此域名吗？")) {
            return;
        }

        Tea.action(".deleteName").params({
            "filename": this.filename,
            "index": index
        })
            .post();
    };

    /**
     * 监听地址管理
     */
    this.newListen = "";
    this.listenAdding = false;
    this.addListen = function () {
        this.listenAdding = true;
    };

    this.addListenSave = function () {
        Tea.action(".addListen")
            .params({
                "filename": this.filename,
                "listen": this.newListen
            })
            .post();
    };

    this.editListenIndex = -1;
    this.editListen = function (index, listen) {
        this.editListenIndex = index;
    };

    this.editListenSave = function (index, listen) {
        Tea.action(".updateListen").params({
            "filename": this.filename,
            "index": index,
            "listen": listen
        })
            .post();
    };

    this.deleteListen = function (index) {
        if (!window.confirm("确定要删除此域名吗？")) {
            return;
        }

        Tea.action(".deleteListen").params({
            "filename": this.filename,
            "index": index
        })
            .post();
    };

    /**
     * SSL管理
     */
    this.showSSLOptions = (this.proxy.ssl != null) ? this.proxy.ssl.on : false;
    this.sslCertFile = null;
    this.sslCertEditing = false;

    this.sslKeyFile = null;
    this.sslKeyEditing = false;

    this.switchSSLOn = function () {
        this.showSSLOptions = !this.showSSLOptions;
        if (this.proxy.ssl == null) {
            this.proxy.ssl = { "on": this.showSSLOptions };
        } else {
            this.proxy.ssl.on = !this.proxy.ssl.on;
        }

        if (this.proxy.ssl.on) {
            Tea.action("/proxy/ssl/on")
                .params({
                    "filename": this.filename
                })
                .post();
        } else {
            Tea.action("/proxy/ssl/off")
                .params({
                    "filename": this.filename
                })
                .post();
        }
    };

    this.changeSSLCertFile = function (event) {
        if (event.target.files.length > 0) {
            this.sslCertFile = event.target.files[0];
        }
    };

    this.uploadSSLCertFile = function () {
        if (this.sslCertFile == null) {
            alert("请先选择证书文件");
            return;
        }

        Tea.action("/proxy/ssl/uploadCert")
            .params({
                "filename": this.filename,
                "certFile": this.sslCertFile
            })
            .post();
    };

    this.changeSSLKeyFile = function (event) {
        if (event.target.files.length > 0) {
            this.sslKeyFile = event.target.files[0];
        }
    };

    this.uploadSSLKeyFile = function () {
        if (this.sslKeyFile == null) {
            alert("请先选择密钥文件");
            return;
        }

        Tea.action("/proxy/ssl/uploadKey")
            .params({
                "filename": this.filename,
                "keyFile": this.sslKeyFile
            })
            .post();
    };

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