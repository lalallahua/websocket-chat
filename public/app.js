new Vue({
    el: '#app', // 将Vue实例绑定到HTML中的id为'app'的元素

    data: {
        ws: null, // WebSocket连接对象
        newMsg: '', // 新消息的内容
        chatContent: '', // 聊天内容的列表，存储在服务器端
        email: null, // 用户的电子邮件
        username: null, // 用户名
        joined: false // 用户是否已加入聊天室，初始为false
    },

    created: function() {
        var self = this; // 保存Vue实例的上下文以便在回调函数中使用
        this.ws = new WebSocket('ws://' + window.location.host + '/ws'); // 创建WebSocket连接
        this.ws.addEventListener('message', function(e) { // 监听WebSocket的消息事件
            var msg = JSON.parse(e.data); // 将收到的消息解析为JSON对象
            self.chatContent += '<div class="chip">'
                + msg.username
                + ': '
                + '</div>'
                + msg.message + '<br/>'; // 将消息内容添加到chatContent中
            var element = document.getElementById('chatMessages');
            element.scrollTop = element.scrollHeight; // 自动滚动到底部
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') { // 如果消息内容不为空
                this.ws.send(
                    JSON.stringify({
                        email: this.email, // 用户的电子邮件
                        username: this.username, // 用户名
                        message: $('<p>').html(this.newMsg).text() // 消息内容
                    })
                );
                this.newMsg = ''; // 发送后清空消息输入框
            }
        },

        join: function () {
            if (!this.email) { // 如果电子邮件为空
                Materialize.toast('An email must be entered', 2000); // 显示提示消息
                return;
            }
            if (!this.username) { // 如果用户名为空
                Materialize.toast('A Username must be entered', 2000); // 显示提示消息
                return;
            }
            // 将电子邮件和用户名中的HTML标签转义，防止XSS攻击
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true; // 标记用户已加入聊天室
        }
    }
});
