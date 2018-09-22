Tea.context(function () {
    this.switchHttpOn = function () {
        var message = this.proxy.http ? "确定要关闭HTTP吗？" : "确定要开启HTTP吗？";
        if (!window.confirm(message)) {
            return false;
        }

        this.proxy.http = !this.proxy.http;
        if (this.proxy.http) {
            this.$get(".httpOn").params({
                "filename": this.filename
            });
        } else {
            this.$get(".httpOff").params({
                "filename": this.filename
            });
        }
    };

    // 代理ID
    this.proxyIdEditing = false;

    this.editId = function () {
        this.proxyIdEditing = !this.proxyIdEditing;
    };

    this.editIdSave = function () {
        this.$post(".updateId")
            .params({
                "filename": this.filename,
                "id": this.proxy.id
            });
    };

    // 代理说明
    this.proxyDescriptionEditing = false;
    this.editDescription = function () {
        this.proxyDescriptionEditing = !this.proxyDescriptionEditing;
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
});